#!/usr/bin/env bash

sleep_time=5

echo "run FT test 03:"
echo "Create an FT, then create an FT with the same denom."
echo "Expected: an error."

mpcli tx marketplace createFT fungible 577 --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time

echo "create FT"

old_ft_count=$(mpcli query marketplace fungible_tokens | grep denom | wc -l)

mpcli tx marketplace createFT fungible 100 --from user1 -y <<< '12345678' >/dev/null

sleep $sleep_time
echo "create dublicate FT"

new_ft_count=$(mpcli query marketplace fungible_tokens | grep denom | wc -l)

amount=$(mpcli query marketplace fungible_tokens | grep -A1 '"denom": "fungible",' | grep -oP '(?<=\"emission_amount\": \")(.*)(?=\".*)' -m 1)

if [[ $amount != 577 ]] && [[ $old_ft_count != $new_ft_count ]]
then
      echo "test FAILURE"
      exit 1
else
      echo "test SUCCESS amount=$amount"
      exit 0
fi
