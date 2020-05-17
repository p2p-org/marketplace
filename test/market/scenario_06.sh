#!/usr/bin/env bash

sleep_time=5

echo "run test 06:"
echo "user1 creates an NFT with valid params, user2 creates an NFT with valid params. user1 puts on market the token created by user2"
echo "Expected: error."

nft_1_id=$(uuidgen)
nft_2_id=$(uuidgen)
user1_id=$(mpcli keys show user1 -a)
user2_id=$(mpcli keys show user2 -a)

mpcli tx nft mint name $nft_1_id $user1_id --from user1 -y <<< '12345678' >/dev/null
mpcli tx nft mint name $nft_2_id $user2_id --from user2 -y <<< '12345678' >/dev/null

sleep $sleep_time

nft_id=$(mpcli query marketplace nft $nft_2_id | ggrep -oP '(?<=\"id\": \")(.*)(?=\".*)' -m 1)

if [[ -z "$nft_id" ]]
then
      echo "Error: token not created"
      exit 1
else
      echo "token created: $nft_id"
fi

seller_id=$(mpcli keys show sellerBeneficiary -a)

mpcli tx marketplace put_on_market $nft_2_id 150token $seller_id --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

nft_sel_ben_id=$(mpcli query marketplace nft $nft_id | ggrep -oP '(?<=\"seller_beneficiary\": \")(.*)(?=\".*)' -m 1)
status=$(mpcli query marketplace nft $nft_id | ggrep -oP '(?<=\"status\": \")(.*)(?=\".*)' -m 1 | tr -d ,)
price=$(mpcli query marketplace nft $nft_id | ggrep -oP '(?<=\"price\": ).*' -m 1)

if [[ $seller_id == $nft_sel_ben_id ]] || [[ $status == "on_market" ]] || [[ $price != "[]," ]]
then
      echo "test FAILURE"
      exit 1
else
      echo "test SUCCESS, no put_on_market passed. $seller_id"
fi