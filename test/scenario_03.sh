#!/usr/bin/env bash

sleep_time=5

echo "run test 03:"
echo "Create an NFT. Transfer NFT from user2 to user1."
echo "Expected: error."

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

old_owner_id=$(mpcli query marketplace nfts | grep -oP '(?<=\"owner\": \")(.*)(?=\".*)' -m 1)

mpcli tx marketplace transfer $nft_id $(mpcli keys show user1 -a) --from user2 -y <<< '12345678' >/dev/null

echo "transfer token"

sleep $sleep_time

new_owner_id=$(mpcli query marketplace nfts | grep -oP '(?<=\"owner\": \")(.*)(?=\".*)' -m 1)

if [[ $new_owner_id != $old_owner_id ]]
then
      echo "test FAILURE"
      exit 1
else
      echo "test SUCCESS, owner was not changed: $new_owner_id"
      exit 0
fi