#!/usr/bin/env bash

sleep_time=5

echo "run test 02:"
echo "Create an NFT. Transfer NFT from user1 to user2."
echo "Expected: user2 becomes the owner of the NFT."

uu=$(uuidgen)
user1_id=$(mpcli keys show user1 -a)
mpcli tx nft mint name $uu $user1_id --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

nft_id=$(mpcli query marketplace nft $uu | grep -oP '(?<=\"id\": \")(.*)(?=\".*)' -m 1)

if [[ $uu != $nft_id ]]
then
      echo "Error: token not created"
      exit 1
else
      echo "token created: $nft_id"
fi

echo "transfer token"
mpcli tx nft transfer $user1_id $(mpcli keys show user2 -a) name $nft_id --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

new_owner_id=$(mpcli query marketplace nft $nft_id | grep -oP '(?<=\"owner\": \")(.*)(?=\".*)' -m 1)
u2_id=$(mpcli keys show user2 -a)

if [[ $new_owner_id != $u2_id ]]
then
      echo "test FAILURE"
      exit 1
else
      echo "test SUCCESS, $new_owner_id"
      exit 0
fi