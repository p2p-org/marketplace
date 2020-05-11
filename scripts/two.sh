#!/bin/bash

# Ensure gopath is set and go is installed
if [[ ! -d $GOPATH ]] || [[ ! -d $GOBIN ]] || [[ ! -x "$(which go)" ]]; then
  echo "Your \$GOPATH is not set or go is not installed,"
  echo "ensure you have a working installation of go before trying again..."
  echo "https://golang.org/doc/install"
  exit 1
fi

MP_DATA="$(pwd)/data"

# ARGS:
# $1 -> local || remote, defaults to remote

# Ensure user understands what will be deleted
if [[ -d $MP_DATA ]] && [[ ! "$2" == "skip" ]]; then
  read -p "$0 will delete \$(pwd)/data folder. Do you wish to continue? (y/n): " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
      exit 1
  fi
fi

rm -rf $MP_DATA &> /dev/null
killall mpd &> /dev/null

set -e

chainid0=ibc0
chainid1=ibc1

echo "Generating mp configurations..."
mkdir -p $MP_DATA && cd $MP_DATA
echo -e "\n" | mpd testnet -o $chainid0 --v 1 --chain-id $chainid0 --node-dir-prefix n --keyring-backend test #&> /dev/null
echo -e "\n" | mpd testnet -o $chainid1 --v 1 --chain-id $chainid1 --node-dir-prefix n --keyring-backend test #&> /dev/null

cfgpth="n0/gaiad/config/config.toml"
if [ "$(uname)" = "Linux" ]; then
  # TODO: Just index *some* specified tags, not all
  sed -i 's/index_all_keys = false/index_all_keys = true/g' $chainid0/$cfgpth
  sed -i 's/index_all_keys = false/index_all_keys = true/g' $chainid1/$cfgpth

  # Set proper defaults and change ports
  sed -i 's/"leveldb"/"goleveldb"/g' $chainid0/$cfgpth
  sed -i 's/"leveldb"/"goleveldb"/g' $chainid1/$cfgpth
  sed -i 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:26556"#g' $chainid1/$cfgpth
  sed -i 's#"tcp://0.0.0.0:26657"#"tcp://0.0.0.0:26557"#g' $chainid1/$cfgpth
  sed -i 's#"localhost:6060"#"localhost:6061"#g' $chainid1/$cfgpth
  sed -i 's#"tcp://127.0.0.1:26658"#"tcp://127.0.0.1:26558"#g' $chainid1/$cfgpth

  # Make blocks run faster than normal
  sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $chainid0/$cfgpth
  sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $chainid1/$cfgpth
  sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $chainid0/$cfgpth
  sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $chainid1/$cfgpth
else
  # TODO: Just index *some* specified tags, not all
  sed -i '' 's/index_all_keys = false/index_all_keys = true/g' $chainid0/$cfgpth
  sed -i '' 's/index_all_keys = false/index_all_keys = true/g' $chainid1/$cfgpth

  # Set proper defaults and change ports
  sed -i '' 's/"leveldb"/"goleveldb"/g' $chainid0/$cfgpth
  sed -i '' 's/"leveldb"/"goleveldb"/g' $chainid1/$cfgpth
  sed -i '' 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:26556"#g' $chainid1/$cfgpth
  sed -i '' 's#"tcp://0.0.0.0:26657"#"tcp://0.0.0.0:26557"#g' $chainid1/$cfgpth
  sed -i '' 's#"localhost:6060"#"localhost:6061"#g' $chainid1/$cfgpth
  sed -i '' 's#"tcp://127.0.0.1:26658"#"tcp://127.0.0.1:26558"#g' $chainid1/$cfgpth

  # Make blocks run faster than normal
  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $chainid0/$cfgpth
  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $chainid1/$cfgpth
  sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $chainid0/$cfgpth
  sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $chainid1/$cfgpth
fi

gclpth="n0/gaiacli/"
mpcli config --home $chainid0/$gclpth chain-id $chainid0 &> /dev/null
mpcli config --home $chainid1/$gclpth chain-id $chainid1 &> /dev/null
mpcli config --home $chainid0/$gclpth output json &> /dev/null
mpcli config --home $chainid1/$gclpth output json &> /dev/null
mpcli config --home $chainid0/$gclpth node http://localhost:26657 &> /dev/null
mpcli config --home $chainid1/$gclpth node http://localhost:26557 &> /dev/null

echo "Starting mpd instances..."
mpd --home $MP_DATA/$chainid0/n0/gaiad start --pruning=nothing --log_level=debug > $chainid0.log 2>&1 &
mpd --home $MP_DATA/$chainid1/n0/gaiad start --pruning=nothing --log_level=debug > $chainid1.log 2>&1 &

#tail -f $chainid1.log