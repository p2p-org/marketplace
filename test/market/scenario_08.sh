#!/usr/bin/env bash

sleep_time=5

echo "run test 08:"
echo "user2 buys a non-existent NFT."
echo "Expected: error."

uu=$(uuidgen)
fake_id=$(uuidgen)
user1_id=$(mpcli keys show user1 -a)
mpcli tx nft mint name $uu $user1_id --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

nft_id=$(mpcli query marketplace nft $uu | grep -oP '(?<=\"id\": \")(.*)(?=\".*)' -m 1)

if [[ -z "$nft_id" ]] || [[ $uu != $nft_id ]]
then
      echo "Error: token not created"
      exit 1
else
      echo "token created: $nft_id"
fi

sleep $sleep_time
seller_id=$(mpcli keys show sellerBeneficiary -a)
buyer_id=$(mpcli keys show buyerBeneficiary -a)

mpcli tx marketplace put_on_market $nft_id 650token $seller_id --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time
echo "token put on market"

mpcli tx marketplace buy $fake_id $buyer_id --from user2 -y <<< '12345678' >/dev/null

sleep $sleep_time
echo "token buy attempt"

new_owner_id=$(mpcli query marketplace nft $nft_id | grep -oP '(?<=\"owner\": \")(.*)(?=\".*)' -m 1)
user1_id=$(mpcli keys show user1 -a)
user2_id=$(mpcli keys show user2 -a)
status=$(mpcli query marketplace nft $nft_id | grep -oP '(?<=\"status\": \")(.*)(?=\".*)' -m 1 | tr -d ,)

if [[ $status != "on_market" ]]
then
      echo "test FAILURE: existing token is not on market"
      exit 1
fi

if [[ $new_owner_id != $user1_id ]]
then
      echo "test FAILURE: token was bought"
      exit 1
fi

fake_owner_id=$(mpcli query marketplace nft $fake_id | grep -oP '(?<=\"owner\": \")(.*)(?=\".*)' -m 1)

if [[ -n "$fake_owner_id" ]]
then
      echo "test FAILURE: fake token exists"
      exit 1
fi

balance_u1=$(mpcli query account $user1_id | grep -A1 '"denom": "token",' | grep -oP '(?<=\"amount\": \").*(?=\".*)' -m 1)
balance_u2=$(mpcli query account $user2_id | grep -A1 '"denom": "token",' | grep -oP '(?<=\"amount\": \").*(?=\".*)' -m 1)
balance_sb=$(mpcli query account $seller_id | grep -A1 '"denom": "token",' | grep -oP '(?<=\"amount\": \").*(?=\".*)' -m 1)
balance_bb=$(mpcli query account $buyer_id | grep -A1 '"denom": "token",' | grep -oP '(?<=\"amount\": \").*(?=\".*)' -m 1)

echo "user1:" $balance_u1
echo "user2:" $balance_u2
echo "sellerBeneficiary:" $balance_sb
echo "buyerBeneficiary:" $balance_bb

if [[ $balance_u1 != 1000 ]] || [[ $balance_u2 != 1000 ]] || [[ $balance_sb != 1000 ]] || [[ $balance_bb != 1000 ]]
then
      echo "FAILURE: wrong numbers"
else
      echo "SUCCESS: NFT was not bought"
fi