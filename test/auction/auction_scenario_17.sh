#!/usr/bin/env bash

sleep_time=5

echo "run test 17:"
echo "Create an NFT. Put this NFT on auction. Make a bid. dgaming makes a greater bid. Wait for the end of the auction."
echo "Expected: NFT belongs to dgaming. NFT status set to default, lot is deleted from auction, all commissions are paid."

uu=$(uuidgen)
boutp=500
bidu=220
bidp=310

user1_id=$(mpcli keys show user1 -a)
mpcli tx nft mint name $uu $user1_id --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

nft_id=$(mpcli query marketplace nft $uu | grep -oP '(?<=\"id\": \")(.*)(?=\".*)' -m 1)

user1_id=$(mpcli keys show user1 -a)
user2_id=$(mpcli keys show user2 -a)
dgaming_id=$(mpcli keys show dgaming -a)
seller_id=$(mpcli keys show sellerBeneficiary -a)
buyer_id=$(mpcli keys show buyerBeneficiary -a)

if [[ $uu != $nft_id ]]
then
      echo "Error: token was not created"
      exit 1
else
      echo "token created: $nft_id"
fi

seller_id=$(mpcli keys show sellerBeneficiary -a)

mpcli tx marketplace put_on_auction $nft_id 200token $seller_id 20s -u 500token --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

nft_sel_ben_id=$(mpcli query marketplace nft $nft_id | grep -oP '(?<=\"seller_beneficiary\": \")(.*)(?=\".*)' -m 1)
status=$(mpcli query marketplace nft $nft_id | grep -oP '(?<=\"status\": \")(.*)(?=\".*)' -m 1 | tr -d ,)

if [[ $seller_id == $nft_sel_ben_id ]] && [[ $status == "on_auction" ]]
then
      echo "nft is on auction"
else
      echo "Error: nft was not put on auction"
      exit 1
fi

mpcli tx marketplace bid $nft_id $buyer_id -c 0.03 ${bidu}token --from user2 -y <<< '12345678' >/dev/null

sleep $sleep_time

bidder=$(mpcli query marketplace auction_lot $nft_id | grep -oP '(?<=\"bidder\": \")(.*)(?=\".*)' -m 1 | tr -d ,)
if [[ $seller_id == $nft_sel_ben_id ]] && [[ $status == "on_auction" ]] && [[ $bidder == $user2_id ]]
then
      echo "nft is on auction, bid was made by user2"
else
      echo "Error: nft was not put on auction"
      exit 1
fi


mpcli tx marketplace bid $nft_id $buyer_id -c 0.04 ${bidp}token --from dgaming -y <<< '12345678' >/dev/null

sleep $sleep_time
sleep $sleep_time
sleep $sleep_time
sleep $sleep_time

nft_sel_ben_id=$(mpcli query marketplace nft $nft_id | grep -oP '(?<=\"seller_beneficiary\": \")(.*)(?=\".*)' -m 1)
status=$(mpcli query marketplace nft $nft_id | grep -oP '(?<=\"status\": \")(.*)(?=\".*)' -m 1 | tr -d ,)
auc_id=$(mpcli query marketplace auction_lot $nft_id | grep -oP '(?<=\"nft_id\": \")(.*)(?=\".*)' -m 1 | tr -d ,)
owner=$(mpcli query marketplace nft $nft_id | grep -oP '(?<=\"owner\": \")(.*)(?=\".*)' -m 1)

balance_u1=$(mpcli query account $user1_id | grep -A1 '"denom": "token",' | grep -oP '(?<=\"amount\": \").*(?=\".*)' -m 1)
balance_u2=$(mpcli query account $user2_id | grep -A1 '"denom": "token",' | grep -oP '(?<=\"amount\": \").*(?=\".*)' -m 1)
balance_dg=$(mpcli query account $dgaming_id | grep -A1 '"denom": "token",' | grep -oP '(?<=\"amount\": \").*(?=\".*)' -m 1)
balance_sb=$(mpcli query account $seller_id | grep -A1 '"denom": "token",' | grep -oP '(?<=\"amount\": \").*(?=\".*)' -m 1)
balance_bb=$(mpcli query account $buyer_id | grep -A1 '"denom": "token",' | grep -oP '(?<=\"amount\": \").*(?=\".*)' -m 1)

echo "user1:" $balance_u1
echo "user2:" $balance_u2
echo "dgaming:" $balance_dg
echo "sellerBeneficiary:" $balance_sb
echo "buyerBeneficiary:" $balance_bb

echo ""

if [[ $balance_u1 != $(($bidp+(1000-($bidp/20)))) ]] || [[ $balance_dg != $((1000-$bidp)) ]] || [[ $balance_sb != $((1000+($bidp/50))) ]] || [[ $balance_bb != $((1000+($bidp/50))) ]] || [[ $balance_u2 != 1000 ]]
then
      echo "FAILURE: wrong numbers"
      exit 1
else
      echo "success, all balances are right"
fi

echo ""

if [[ -z $nft_sel_ben_id ]] && [[ $status == "default" ]] && [[ -z $auc_id ]] && [[ $owner == $dgaming_id ]]
then
      echo "test SUCCESS, nft was selled to bidder: $(mpcli query marketplace nft $nft_id)"
else
      echo "test FAILURE"
      exit 1
fi