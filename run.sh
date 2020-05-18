#!/bin/bash

DATA_DIR=${DATA_DIR:-"data"}
CHAIN_ID=${CHAIN_ID:-"mpchain"}
RPC_HOST=${RPC_HOST:-"localhost"}
P2P_PORT=${P2P_PORT:-"26656"}
RPC_PORT=${RPC_PORT:-"26657"}
PROXY_PORT=${PROXY_PORT:-"26658"}
LCD_PORT=${LCD_PORT:-"1317"}

MP_DATA="$DATA_DIR/$CHAIN_ID"

gclpth="n0/gaiacli"
gdpth="n0/gaiad"
gclhome="$MP_DATA/$gclpth"
gdhome="$MP_DATA/$gdpth"
cfgpth="$gdhome/config/config.toml"

while test $# -gt 0; do
  case "$1" in
    -h|--help)
      echo "run relayer"
      echo " "
      echo "relay.sh [options]"
      echo " "
      echo "options:"
      echo "--init                   force reset blockchain state (clear all data)"
      echo "--norun                  omit service start"
      echo "--rest                   start rest server"
      exit 0
      ;;
    --init)
      INIT=true
      shift
      ;;
    --rest)
      REST=true
      shift
      ;;
    --norun)
      NORUN=true
      shift
      ;;
    *)
      break
      ;;
  esac
done

if [[ $INIT ]]; then
  echo "Clearing data..."
  rm -rf $MP_DATA &> /dev/null
fi


if [ ! -d "$MP_DATA" ]; then
set -e

echo "Generating mp configurations..."
mkdir -p $MP_DATA
echo -e "\n" | mpd testnet -o $MP_DATA --v 1 --chain-id $CHAIN_ID --node-dir-prefix n --keyring-backend test
 #&> /dev/null
echo "$(pwd)"

if [ "$(uname)" = "Linux" ]; then
  # TODO: Just index *some* specified tags, not all
  sed -i "s#index_all_keys = false#index_all_keys = true#g" $cfgpth

  # Set proper defaults and change ports
  sed -i 's#"leveldb"#"goleveldb"#g' $cfgpth
  sed -i "s#:26656#:$P2P_PORT#g" $cfgpth
  sed -i "s#:26657#:$RPC_PORT#g" $cfgpth
  #sed -i 's#"localhost:6060"#"localhost:6061"#g' $cfgpth
  sed -i "s#:26658#:$PROXY_PORT#g" $cfgpth

  # Make blocks run faster than normal
  sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $cfgpth
  sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $cfgpth
else
  # TODO: Just index *some* specified tags, not all
  sed -i '' "s#index_all_keys = false#index_all_keys = true#g" $cfgpth

  # Set proper defaults and change ports
  sed -i '' 's#"leveldb"#"goleveldb"#g' $cfgpth
  sed -i '' "s#:26656#:$P2P_PORT#g" $cfgpth
  sed -i '' "s#:26657#:$RPC_PORT#g" $cfgpth
  #sed -i '' 's#"localhost:6060"#"localhost:6061"#g' $cfgpth
  sed -i '' "s#:26658#:$PROXY_PORT#g" $cfgpth

  # Make blocks run faster than normal
  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $cfgpth
  sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $cfgpth
fi



mpcli --home $gclhome keys add user1 -i --keyring-backend test < data_u1.txt
mpcli --home $gclhome keys add user2 -i --keyring-backend test < data_u2.txt
mpcli --home $gclhome keys add user3 -i --keyring-backend test < data_u3.txt
#mpcli --home $gclhome keys add sellerBeneficiary -i --keyring-backend test < "data_sb.txt"
#mpcli --home $gclhome keys add buyerBeneficiary -i --keyring-backend test < "data_bb.txt"

mpd --home $gdhome add-genesis-account $(mpcli --home $gclhome keys show user1 -a --keyring-backend test) 999999token,100000000stake --keyring-backend test
mpd --home $gdhome add-genesis-account $(mpcli --home $gclhome keys show user2 -a --keyring-backend test) 999999token,100000000stake --keyring-backend test
mpd --home $gdhome add-genesis-account $(mpcli --home $gclhome keys show user3 -a --keyring-backend test) 999999token,100000000stake --keyring-backend test

mpcli --home $gclhome config chain-id $CHAIN_ID &> /dev/null
mpcli --home $gclhome config output json &> /dev/null
mpcli --home $gclhome config node http://localhost:$RPC_PORT &> /dev/null
#
##  mpcli config indent true
##  mpcli config trust-node true
#
#
#  mpd --home $gdhome gentx --name user1 --keyring-backend test
#  mpd --home $gdhome collect-gentxs
#  mpd --home $gdhome validate-genesis

fi

if [[ -z $NORUN ]]; then
  echo "Starting mpd instances..."
  if [[ -z $REST ]]; then
    mpd --home $gdhome start --pruning=nothing --log_level=debug #> $CHAIN_ID.log 2>&1 &
  else
    mpcli rest-server --chain-id $CHAIN_ID --trust-node --node tcp://$RPC_HOST:$RPC_PORT --laddr tcp://0.0.0.0:$LCD_PORT #> "${CHAIN_ID}-rest".log 2>&1 &
  fi
fi
#tail -f $chainid1.log