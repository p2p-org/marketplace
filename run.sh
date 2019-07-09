#!/usr/bin/env bash

echo "Clearing previous files..."
rm -rf ~/.mp*

echo "Building..."
make install

echo "Initialization..."
mpd init node0 --chain-id mpchain

echo "Adding keys..."

echo "Adding genesis account..."
mpcli keys add user1 --recover <<< "12345678
define hurry shoot window find now soul fly live cruel elevator harvest cradle great charge such box post brass midnight glimpse forest jaguar ankle
"

mpcli keys add user2 --recover <<< "12345678
base figure planet hazard sail easily honey advance tuition grab across unveil random kiss fence connect disagree evil recall latin cause brisk soft lunch
"

mpd add-genesis-account $(mpcli keys show user1 -a) 1000token,100000000stake
mpd add-genesis-account $(mpcli keys show user2 -a) 1000token,100000000stake

echo "Configuring..."
mpcli config chain-id mpchain
mpcli config output json
mpcli config indent true
mpcli config trust-node true

mpd gentx --name user1 <<< "12345678"
mpd gentx --name user2 <<< "123456789"
mpd collect-gentxs
mpd validate-genesis

echo "Starting node..."

mpd start
