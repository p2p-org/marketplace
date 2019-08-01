#!/usr/bin/env bash

sleep_time=5

echo "run test 04:"
echo "Create an NFT. Put this NFT on market with correct token ID, price and seller beneficiary address."
echo "Expected: the NFT is updated with price and seller beneficiary address, its OnSale field equals true."

uu=$(uuidgen)
mpcli tx marketplace mint $uu name description image token_uri --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

nft_id=$(mpcli query marketplace nfts | grep -oP '(?<=\"id\": \")(.*)(?=\".*)' -m 1)

if [[ $uu != $nft_id ]]
then
      echo "Error: token not created"
      exit 1
else
      echo "token created: $nft_id"
fi

seller_id=$(mpcli keys show sellerBeneficiary -a)

mpcli tx marketplace put_on_market $nft_id 150token $seller_id --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

nft_sel_ben_id=$(mpcli query marketplace nfts | grep -oP '(?<=\"seller_beneficiary\": \")(.*)(?=\".*)' -m 1)
is_on_sale=$(mpcli query marketplace nfts | grep -oP '(?<=\"on_sale\": )(.*)(?=.*)' -m 1 | tr -d ,)
price=$(mpcli query marketplace nfts | grep -oP '(?<=\"price\": ).*' -m 1)

if [[ $seller_id == $nft_sel_ben_id ]] && [[ $is_on_sale == "true" ]] && [[ $price != "[]," ]]
then
      echo "test SUCCESS, nft was put on market: $seller_id"
else
      echo "test FAILURE"
      exit 1
fi