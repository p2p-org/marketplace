#!/usr/bin/env bash

sleep_time=4

echo "run test 01:"
echo "Create some NFTs. Put this NFT on market as batch with correct token ID, price and seller beneficiary address."


user1_id=$(mpcli keys show user1 -a)
declare -a uu_arr
declare -a nft_arr
for i in {0..4}
do
  uu=$(uuidgen)
  echo $uu
  mpcli tx nft mint name $uu $user1_id --from user1 -y <<< '12345678' >/dev/null
  uu_arr+=($uu)
  sleep $sleep_time

done

echo

for i in {0..4}
do
  nft_id=$(mpcli query marketplace nft ${uu_arr[i]} | grep -oP '(?<=\"id\": \")(.*)(?=\".*)' -m 1)
  nft_arr+=($nft_id)
  if [[ ${uu_arr[i]} != $nft_id ]]
  then
        echo "Error: token $uu_arr not created"
        exit 1
  else
        echo "token created: $nft_id"
  fi
done

seller_id=$(mpcli keys show sellerBeneficiary -a)
JSON_STRING="{\"${uu_arr[0]}\": \"100token\", \"${uu_arr[1]}\": \"20token\", \"${uu_arr[2]}\": \"12token\", \"${uu_arr[3]}\": \"1token\", \"${uu_arr[4]}\": \"100token\"}"
echo $JSON_STRING

mpcli tx marketplace batch_put_on_market $seller_id "$JSON_STRING" --from user1

