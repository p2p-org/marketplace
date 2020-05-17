#!/usr/bin/env bash

sleep_time=5

echo "run test 03:"
echo "Create an NFT. Transfer NFT from user2 to user1."
echo "Expected: error."

uu=$(uuidgen)
user1_id=$(mpcli keys show user1 -a)
mpcli tx nft mint name $uu $user1_id --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

nft_id=$(mpcli query marketplace nft $uu | ggrep -oP '(?<=\"id\": \")(.*)(?=\".*)' -m 1)

if [[ $uu != $nft_id ]]
then
      echo "Error: token not created"
      exit 1
else
      echo "token created: $nft_id"
fi

old_owner_id=$(mpcli query marketplace nft $nft_id | ggrep -oP '(?<=\"owner\": \")(.*)(?=\".*)' -m 1)

mpcli tx nft transfer $(mpcli keys show user2 -a) $user1_id name $nft_id --from user1 -y <<< '12345678' >/dev/null

echo "transfer token"

sleep $sleep_time

new_owner_id=$(mpcli query marketplace nft $nft_id | ggrep -oP '(?<=\"owner\": \")(.*)(?=\".*)' -m 1)

if [[ $new_owner_id != $old_owner_id ]]
then
      echo "test FAILURE"
      exit 1
else
      echo "test SUCCESS, owner was not changed: $new_owner_id"
      exit 0
fi