#!/usr/bin/env bash

sleep_time=5

if [[ $1 == "" ]]
then
      start=650
else
      start=$1
fi

echo "run test 07:"
echo "Create an NFT and put it on market. user2 buys this NFT."
echo "Expected: user2 becomes the owner of NFT. user2 loses exactly NFT.Price. user1 receives exactly NFT.Price * 0.95 tokens.
Both sellerBeneficiary and buyerBeneficiary receive NFT.Price * 0.02 tokens.
Token's OnSale field equals false. (7.5).
Scenario No. 7 should be performed both with no integer rounding and with integer rounding."

uu=$(uuidgen)
user1_id=$(mpcli keys show user1 -a)
mpcli tx nft mint name $uu $user1_id --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

nft_id=$(mpcli query marketplace nft $uu | ggrep -oP '(?<=\"id\": \")(.*)(?=\".*)' -m 1)

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

mpcli tx marketplace put_on_market $nft_id ${start}token $seller_id --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time
echo "token put on market"

mpcli tx marketplace buy $nft_id $buyer_id -c 0.04 --from user2 -y <<< '12345678' >/dev/null

sleep $sleep_time
echo "token bought"

new_owner_id=$(mpcli query marketplace nft $nft_id | ggrep -oP '(?<=\"owner\": \")(.*)(?=\".*)' -m 1)
user1_id=$(mpcli keys show user1 -a)
user2_id=$(mpcli keys show user2 -a)
status=$(mpcli query marketplace nft $nft_id | ggrep -oP '(?<=\"status\": \")(.*)(?=\".*)' -m 1 | tr -d ,)

if [[ $status == "on_market" ]]
then
      echo "test FAILURE: token is still on market"
      exit 1
fi

if [[ $new_owner_id != $user2_id ]]
then
      echo "test FAILURE: token was not bought"
      exit 1
fi

balance_u1=$(mpcli query account $user1_id | ggrep -A1 '"denom": "token",' | ggrep -oP '(?<=\"amount\": \").*(?=\".*)' -m 1)
balance_u2=$(mpcli query account $user2_id | ggrep -A1 '"denom": "token",' | ggrep -oP '(?<=\"amount\": \").*(?=\".*)' -m 1)
balance_sb=$(mpcli query account $seller_id | ggrep -A1 '"denom": "token",' | ggrep -oP '(?<=\"amount\": \").*(?=\".*)' -m 1)
balance_bb=$(mpcli query account $buyer_id | ggrep -A1 '"denom": "token",' | ggrep -oP '(?<=\"amount\": \").*(?=\".*)' -m 1)

echo "user1:" $balance_u1
echo "user2:" $balance_u2
echo "sellerBeneficiary:" $balance_sb
echo "buyerBeneficiary:" $balance_bb

if [[ $balance_u1 != $(($start+(1000-($start/20)))) ]] || [[ $balance_u2 != $((1000-$start)) ]] || [[ $balance_sb != $((1000+($start/50))) ]] || [[ $balance_bb != $((1000+($start/50))) ]]
then
      echo "FAILURE: wrong numbers"
else
      echo "SUCCESS"
fi
