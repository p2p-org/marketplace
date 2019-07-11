#!/usr/bin/env bash

echo "Clearing previous files..."
rm -rf ~/.mp*

echo "Building..."
make install

echo "Initialization..."
mpd init node0 --chain-id mpchain

echo "Adding keys..."

echo "Adding genesis account..."
mpcli keys add user1 <<< "12345678"
mpcli keys add user2 <<< "12345678"
mpcli keys add sellerBeneficiary <<< "12345678"
mpcli keys add buyerBeneficiary <<< "12345678"

mpd add-genesis-account $(mpcli keys show user1 -a) 1000token,100000000stake
mpd add-genesis-account $(mpcli keys show user2 -a) 1000token,100000000stake
mpd add-genesis-account $(mpcli keys show sellerBeneficiary -a) 1000token,100000000stake
mpd add-genesis-account $(mpcli keys show buyerBeneficiary -a) 1000token,100000000stake

echo "Configuring..."
mpcli config chain-id mpchain
mpcli config output json
mpcli config indent true
mpcli config trust-node true

mpd gentx --name user1 <<< "12345678"
mpd collect-gentxs
mpd validate-genesis

echo "Starting node..."

mpd start
