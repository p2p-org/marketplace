#!/bin/zsh

MP_BRANCH=feat/ibc-nftt
MP_DIR=$(mktemp -d)
CONF_DIR=$(mktemp -d)

echo "MP_DIR: ${MP_DIR}"
echo "CONF_DIR: ${CONF_DIR}"

sleep 1

set -x

echo "Killing existing mpd instances..."

killall mpd

set -e

echo "Building Gaia..."

cd $MP_DIR
git clone git@github.com:corestario/marketplace.git
cd marketplace
git checkout $MP_BRANCH
make install

echo "Generating configurations..."

cd $CONF_DIR && mkdir ibc-testnets && cd ibc-testnets
echo -e "\n" | mpd testnet -o ibc0 --v 1 --chain-id ibc0 --node-dir-prefix n
echo -e "\n" | mpd testnet -o ibc1 --v 1 --chain-id ibc1 --node-dir-prefix n

if [ "$(uname)" = "Linux" ]; then
  sed -i 's/"leveldb"/"goleveldb"/g' ibc0/n0/mpd/config/config.toml
  sed -i 's/"leveldb"/"goleveldb"/g' ibc1/n0/mpd/config/config.toml
  sed -i 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:26556"#g' ibc1/n0/mpd/config/config.toml
  sed -i 's#"tcp://0.0.0.0:26657"#"tcp://0.0.0.0:26557"#g' ibc1/n0/mpd/config/config.toml
  sed -i 's#"localhost:6060"#"localhost:6061"#g' ibc1/n0/mpd/config/config.toml
  sed -i 's#"tcp://127.0.0.1:26658"#"tcp://127.0.0.1:26558"#g' ibc1/n0/mpd/config/config.toml
else
  sed -i '' 's/"leveldb"/"goleveldb"/g' ibc0/n0/mpd/config/config.toml
  sed -i '' 's/"leveldb"/"goleveldb"/g' ibc1/n0/mpd/config/config.toml
  sed -i '' 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:26556"#g' ibc1/n0/mpd/config/config.toml
  sed -i '' 's#"tcp://0.0.0.0:26657"#"tcp://0.0.0.0:26557"#g' ibc1/n0/mpd/config/config.toml
  sed -i '' 's#"localhost:6060"#"localhost:6061"#g' ibc1/n0/mpd/config/config.toml
  sed -i '' 's#"tcp://127.0.0.1:26658"#"tcp://127.0.0.1:26558"#g' ibc1/n0/mpd/config/config.toml
fi;

mpcli config --home ibc0/n0/mpcli/ chain-id ibc0
mpcli config --home ibc1/n0/mpcli/ chain-id ibc1
mpcli config --home ibc0/n0/mpcli/ output json
mpcli config --home ibc1/n0/mpcli/ output json
mpcli config --home ibc0/n0/mpcli/ node http://localhost:26657
mpcli config --home ibc1/n0/mpcli/ node http://localhost:26557

mpcli config --home ibc0/n0/mpcli/ trust-node true
mpcli config --home ibc1/n0/mpcli/ trust-node true

echo "Importing keys..."

SEED0=$(jq -r '.secret' ibc0/n0/mpcli/key_seed.json)
SEED1=$(jq -r '.secret' ibc1/n0/mpcli/key_seed.json)
echo -e "12345678\n" | mpcli --home ibc1/n0/mpcli keys delete n0

echo "Seed 0: ${SEED0}"
echo "Seed 1: ${SEED1}"

mpcli keys test --home ibc0/n0/mpcli n1 "$(jq -r '.secret' ibc1/n0/mpcli/key_seed.json)" 12345678
mpcli keys test --home ibc1/n0/mpcli n0 "$(jq -r '.secret' ibc0/n0/mpcli/key_seed.json)" 12345678
mpcli keys test --home ibc1/n0/mpcli n1 "$(jq -r '.secret' ibc1/n0/mpcli/key_seed.json)" 12345678

echo "Keys should match:"

mpcli --home ibc0/n0/mpcli keys list | jq '.[].address'
mpcli --home ibc1/n0/mpcli keys list | jq '.[].address'

echo "Starting Gaiad instances..."

nohup mpd --home ibc0/n0/mpd --log_level="*:debug" start > ibc0.log &
nohup mpd --home ibc1/n0/mpd --log_level="*:debug" start > ibc1.log &

sleep 20

echo "Creating clients..."

echo -e "12345678\n" | mpcli --home ibc0/n0/mpcli \
  tx ibc client create ibconeclient \
  $(mpcli --home ibc1/n0/mpcli q ibc client node-state) \
  --from n0 -y -o text

echo -e "12345678\n" | mpcli --home ibc1/n0/mpcli \
  tx ibc client create ibczeroclient \
  $(mpcli --home ibc0/n0/mpcli q ibc client node-state) \
  --from n1 -y -o text

sleep 3

echo "Querying clients..."

mpcli --home ibc0/n0/mpcli q ibc client consensus-state ibconeclient --indent
mpcli --home ibc1/n0/mpcli q ibc client consensus-state ibczeroclient --indent

echo "Establishing a connection..."

mpcli \
  --home ibc0/n0/mpcli \
  tx ibc connection handshake \
  connectionzero ibconeclient $(mpcli --home ibc1/n0/mpcli q ibc client path) \
  connectionone ibczeroclient $(mpcli --home ibc0/n0/mpcli q ibc client path) \
  --chain-id2 ibc1 \
  --from1 n0 --from2 n1 \
  --node1 tcp://localhost:26657 \
  --node2 tcp://localhost:26557

sleep 2

echo "Querying connection..."

mpcli --home ibc0/n0/mpcli q ibc connection end connectionzero --indent --trust-node
mpcli --home ibc1/n0/mpcli q ibc connection end connectionone --indent --trust-node

echo "Establishing a channel..."

mpcli \
  --home ibc0/n0/mpcli \
  tx ibc channel handshake \
  ibconeclient transfernft channelzero connectionzero \
  ibczeroclient transfernft channelone connectionone \
  --node1 tcp://localhost:26657 \
  --node2 tcp://localhost:26557 \
  --chain-id2 ibc1 \
  --from1 n0 --from2 n1

sleep 2

echo "Querying channel..."

mpcli --home ibc0/n0/mpcli q ibc channel end transfernft channelzero --indent --trust-node
mpcli --home ibc1/n0/mpcli q ibc channel end transfernft channelone --indent --trust-node

echo "Sending token packets from ibc0..."

DEST=$(mpcli --home ibc0/n0/mpcli keys show n1 -a)

#mpcli \
#  --home ibc0/n0/mpcli \
#  tx ibc transfer transfer \
#  bank channelzero \
#  $DEST 1stake \
#  --from n0 \
#  --source

TOKEN_ID=$(uuidgen)
mpcli --home ibc0/n0/mpcli tx nft mint name $TOKEN_ID $(mpcli --home ibc0/n0/mpcli keys show n0 -a) --from n0

sleep 5

mpcli \
  --home ibc0/n0/mpcli \
  tx marketplace transferNFT \
  transfernft channelzero \
  $DEST $TOKEN_ID \
  --from n0 \
  --source

echo "Enter height:"

read -r HEIGHT

TIMEOUT=$(echo "$HEIGHT + 1000" | bc -l)

echo "Account before:"
mpcli --home ibc1/n0/mpcli query marketplace nfts

echo "Recieving token packets on ibc1..."

sleep 3

mpcli \
  tx ibc transfer recv-packet \
  transfernft channelzero ibczeroclient \
  --home ibc1/n0/mpcli \
  --packet-sequence 1 \
  --timeout $TIMEOUT \
  --from n1 \
  --node2 tcp://localhost:26657 \
  --chain-id2 ibc0 \
  --source

echo "Account after:"

mpcli --home ibc1/n0/mpcli query marketplace nfts