#!/usr/bin/env bash

sleep_time=5

echo "run test 05:"
echo "Create an NFT. First put this NFT on market with incorrect token ID (non-existent token, ),
then with incorrect price (non-existent denomination),
and then with incorrect seller beneficiary address (incorrect address)."
echo "Expected: error (for each of the three cases)."

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

mpcli tx marketplace put_on_market $(uuidgen) 150token $seller_id --from user1 -y <<< '12345678'
echo ""
mpcli tx marketplace put_on_market $nft_id 150uuidgen $seller_id --from user1 -y <<< '12345678'
echo ""
mpcli tx marketplace put_on_market $nft_id 150token $(uuidgen) --from user1 -y <<< '12345678'
echo ""

sleep $sleep_time

nft_sel_ben_id=$(mpcli query marketplace nfts | grep -oP '(?<=\"seller_beneficiary\": \")(.*)(?=\".*)' -m 1)
is_on_sale=$(mpcli query marketplace nfts | grep -oP '(?<=\"on_sale\": )(.*)(?=.*)' -m 1 | tr -d ,)
price=$(mpcli query marketplace nfts | grep -oP '(?<=\"price\": ).*' -m 1)

if [[ $seller_id == $nft_sel_ben_id ]] || [[ $is_on_sale == "true" ]] || [[ $price != "[]," ]]
then
      echo "test FAILURE"
      exit 1
else
      echo "test SUCCESS, no put_on_market passed. $seller_id"
fi