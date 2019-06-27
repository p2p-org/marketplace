#!/usr/bin/env bash

echo "Clearing previous files..."
rm -rf ~/.mp*

echo "Building..."
make install

echo "Initialization..."
mpd init node0 --chain-id mpchain

echo "Adding keys..."

echo "Adding genesis account..."
mpcli keys add validator1 --recover <<< "12345678
base figure planet hazard sail easily honey advance tuition grab across unveil random kiss fence connect disagree evil recall latin cause brisk soft lunch
"

mpd add-genesis-account $(mpcli keys show validator1 -a) 1000token,100000000stake

echo "Configuring..."
mpcli config chain-id mpchain
mpcli config output json
mpcli config indent true
mpcli config trust-node true

mpd gentx --name validator1 <<< "12345678"
mpd collect-gentxs
mpd validate-genesis

echo "Starting node..."

mpd start
