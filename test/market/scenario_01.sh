#!/usr/bin/env bash

sleep_time=5

echo "run test 01:"
echo "Create an NFT."
echo "Expected: NFT is created"

uu=$(uuidgen)
user1_id=$(mpcli keys show user1 -a)
mpcli tx nft mint name $uu $user1_id --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

nft_id=$(mpcli query marketplace nft $uu | grep -oP '(?<=\"id\": \")(.*)(?=\".*)' -m 1)

if [[ $uu != $nft_id ]]
then
      echo "test FAILURE"
      exit 1
else
      echo "test SUCCESS $nft_id"
      exit 0
fi