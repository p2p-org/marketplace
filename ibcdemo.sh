#!/bin/bash

rly config init

# Then add the chains and paths that you will need to work with the
# gaia chains spun up by the two-chains script
rly cfg add-dir demoIBC/

# NOTE: you may want to look at the config between these steps
#cat ~/.relayer/config/config.yaml

# Now, add the key seeds from each chain to the relayer to give it funds to work with
rly keys restore ibc0 testkey "$(jq -r '.secret' scripts/data/ibc0/n0/gaiacli/key_seed.json)"
rly keys restore ibc1 testkey "$(jq -r '.secret' scripts/data/ibc1/n0/gaiacli/key_seed.json)"

# Then its time to initialize the relayer's lite clients for each chain
# All data moving forward is validated by these lite clients.
rly lite init ibc0 -f
rly lite init ibc1 -f

# Now you can connect the two chains with one command:
rly tx link demo


TOKEN_ID=9E1FAAD1-BA51-4ED9-A0DB-00D096F807DD
DENOM=denom

mpcli tx nft mint $DENOM $TOKEN_ID $(mpcli keys show n0 --home scripts/data/ibc0/n0/gaiacli/ --keyring-backend test -a) --tokenURI someTOKENURI --from n0 --home scripts/data/ibc0/n0/gaiacli/ --keyring-backend test

sleep 5

echo "----------------------------"
echo "Minted NFT on ibc0"
mpcli q marketplace nfts --home scripts/data/ibc0/n0/gaiacli/

echo "----------------------------"
echo "Transfering NFT to ibc1...\n"
rly tx transferNFT ibc0 ibc1 $TOKEN_ID $DENOM true $(rly ch addr ibc1)

echo "----------------------------"
echo "Transferred NFT on ibc0 (owned by escrow account)"
mpcli q marketplace nfts --home scripts/data/ibc0/n0/gaiacli/

sleep 5
echo "----------------------------"
echo "Relaying packet"
rly tx relay demo

echo "----------------------------"
echo "Relayed NFT on ibc1"
mpcli q marketplace nfts --home scripts/data/ibc1/n0/gaiacli/