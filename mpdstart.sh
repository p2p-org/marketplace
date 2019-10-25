#!/bin/zsh

PSW="12345678"

if [[ ! -z $MPD_INIT ]]; then
    echo "Clearing previous files..."
    rm -rf ~/.mp*
    mkdir -p ~/.mpd/config
    cp config.toml ~/.mpd/config

    echo "Initialization..."
    # to specify maximum beneficiary fee use flag
    # max-commission [0.05]
    # example:
    # mpd init node0 --chain-id mpchain --max-commission 0.07

    mpd init node0 --chain-id mpchain

    echo "Adding genesis accounts..."
    mpcli keys add user1 -i <<< $PSW < data_u1.txt
    mpcli keys add user2 -i <<< $PSW < data_u2.txt
    mpcli keys add sellerBeneficiary -i <<< $PSW < data_sb.txt
    mpcli keys add buyerBeneficiary -i <<< $PSW < data_bb.txt
    mpcli keys add dgaming -i <<< $PSW < data_dg.txt

    mpd add-genesis-account $(mpcli keys show user1 -a) 999999token,100000000stake
    mpd add-genesis-account $(mpcli keys show user2 -a) 999999token,100000000stake
    mpd add-genesis-account $(mpcli keys show sellerBeneficiary -a) 1000token,100000000stake
    mpd add-genesis-account $(mpcli keys show buyerBeneficiary -a) 1000token,100000000stake
    mpd add-genesis-account $(mpcli keys show dgaming -a) 1000token,100000000stake

    echo "Configuring..."
    mpcli config chain-id mpchain
    mpcli config output json
    mpcli config indent true
    mpcli config trust-node true

    mpd gentx --name user1 <<< $PSW
    mpd collect-gentxs
    mpd validate-genesis
fi


echo "Starting node..."
mpcli rest-server --chain-id mpchain --trust-node --laddr tcp://0.0.0.0:1317 > /dev/null &
mpd start

