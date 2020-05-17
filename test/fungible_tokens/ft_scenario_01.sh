#!/usr/bin/env bash

sleep_time=5

echo "run FT test 01:"
echo "Create an FT."
echo "Expected: the currency is created, it is listed among the fungible tokens list,
its creator has the exact amount of tokens that she requested, user dgaming received token creation fee."

mpcli tx marketplace createFT fungible 577 --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

amount=$(mpcli query marketplace fungible_tokens | ggrep -A1 '"denom": "fungible",' | ggrep -oP '(?<=\"emission_amount\": \")(.*)(?=\".*)' -m 1)
user_amount=$(mpcli query account $(mpcli keys show user1 -a) | ggrep -A1 '"denom": "fungible",' | ggrep -oP '(?<=\"amount\": \")(.*)(?=\".*)' -m 1)

if [[ $amount == 577 ]] && [[ $amount == $user_amount ]]
then
      echo "test SUCCESS amount=$amount"
      exit 0
else
      echo "test FAILURE"
      exit 1
fi