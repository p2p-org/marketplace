#!/usr/bin/env bash

sleep_time=5
echo "run test 09:"
echo "Create an NFT. Put this NFT on auction. user2 makes a bid equal opening price. dgaming makes a greater bid."
echo "Expected: auction_lot's field lastBid updated with dgaming information."

uu=$(uuidgen)
user1_id=$(mpcli keys show user1 -a)
mpcli tx nft mint name $uu $user1_id --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

nft_id=$(mpcli query marketplace nft $uu | grep -oP '(?<=\"id\": \")(.*)(?=\".*)' -m 1)

if [[ $uu != $nft_id ]]
then
      echo "Error: token was not created"
      exit 1
else
      echo "token created: $nft_id"
fi

seller_id=$(mpcli keys show sellerBeneficiary -a)
buyer_id=$(mpcli keys show buyerBeneficiary -a)

mpcli tx marketplace put_on_auction $nft_id 200token $seller_id 5h -u 500token --from user1 -y <<< '12345678' >/dev/null

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

mpcli tx marketplace bid $nft_id $buyer_id 200token -c 0.02  --from user2 -y <<< '12345678' >/dev/null

sleep $sleep_time

nft_sel_ben_id=$(mpcli query marketplace nft $nft_id | grep -oP '(?<=\"seller_beneficiary\": \")(.*)(?=\".*)' -m 1)
status=$(mpcli query marketplace nft $nft_id | grep -oP '(?<=\"status\": \")(.*)(?=\".*)' -m 1 | tr -d ,)
auc_id=$(mpcli query marketplace auction_lot $nft_id | grep -oP '(?<=\"nft_id\": \")(.*)(?=\".*)' -m 1 | tr -d ,)
auc_price=$(mpcli query marketplace auction_lot $nft_id | grep -oP '(?<=\"amount\": \")(200)(?=\".*)' -m 1 | tr -d ,)
bidder=$(mpcli query marketplace auction_lot $nft_id | grep -oP '(?<=\"bidder\": \")(.*)(?=\".*)' -m 1 | tr -d ,)

if [[ $status == "on_auction" ]] && [[ $auc_price == 200 ]] && [[ $bidder == $(mpcli keys show user2 -a) ]]
then
      echo "test last bid made by user2"
else
      echo "Error: bid was not made"
      exit 1
fi

mpcli tx marketplace bid $nft_id $buyer_id 201token -c 0.02  --from dgaming -y <<< '12345678' >/dev/null

sleep $sleep_time

status=$(mpcli query marketplace nft $nft_id | grep -oP '(?<=\"status\": \")(.*)(?=\".*)' -m 1 | tr -d ,)
auc_id=$(mpcli query marketplace auction_lot $nft_id | grep -oP '(?<=\"nft_id\": \")(.*)(?=\".*)' -m 1 | tr -d ,)
auc_price=$(mpcli query marketplace auction_lot $nft_id | grep -oP '(?<=\"amount\": \")(201)(?=\".*)' -m 1 | tr -d ,)
bidder=$(mpcli query marketplace auction_lot $nft_id | grep -oP '(?<=\"bidder\": \")(.*)(?=\".*)' -m 1 | tr -d ,)


if [[ $status == "on_auction" ]] && [[ $auc_price == 201 ]] && [[ $bidder == $(mpcli keys show dgaming -a) ]]
then
      echo "test SUCCESS, bid was made by user2: $(mpcli query marketplace auction_lot $nft_id)"
else
      echo "test FAILURE"
      exit 1
fi