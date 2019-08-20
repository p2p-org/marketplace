#!/usr/bin/env bash

sleep_time=5

echo "run test 04:"
echo "Create an NFT. Put this NFT on auction. Remove NFT from auction."
echo "Expected: NFT Status becomes default, sellerBeneficiary becomes empty. Auction lot is removed."


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

mpcli tx marketplace remove_from_auction $nft_id --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

nft_sel_ben_id=$(mpcli query marketplace nft $nft_id | grep -oP '(?<=\"seller_beneficiary\": \")(.*)(?=\".*)' -m 1)
status=$(mpcli query marketplace nft $nft_id | grep -oP '(?<=\"status\": \")(.*)(?=\".*)' -m 1 | tr -d ,)
auc_id=$(mpcli query marketplace auction_lot $nft_id | grep -oP '(?<=\"nft_id\": \")(.*)(?=\".*)' -m 1 | tr -d ,)

if [[ -z $nft_sel_ben_id ]] && [[ $status == "default" ]] && [[ -z $auc_id ]]
then
      echo "test SUCCESS, nft was returned to owner: $(mpcli query marketplace nft $nft_id)"
else
      echo "test FAILURE"
      exit 1
fi