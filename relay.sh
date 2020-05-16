#!/bin/bash

while test $# -gt 0; do
  case "$1" in
    -h|--help)
      echo "run relayer"
      echo " "
      echo "relay.sh [options]"
      echo " "
      echo "options:"
      echo "--reset                   force reset blockchain state (clear all data)"
      echo "--norun                   omit service start"
      exit 0
      ;;
    --init)
      INIT=true
      shift
      ;;
    --norun)
      NORUN=true
      shift
      ;;
    *)
      break
      ;;
  esac
done


if [[ $INIT ]]; then
  echo "Init..."
  rm -rf ~/.relayer*

  rly config init

# Then add the chains and paths that you will need to work with the
# gaia chains spun up by the two-chains script
rly cfg add-dir demoIBC/

# NOTE: you may want to look at the config between these steps
#cat ~/.relayer/config/config.yaml

# Now, add the key seeds from each chain to the relayer to give it funds to work with
rly keys restore ibc0 testkey "$($(head -n 1 filename))"
rly keys restore ibc1 testkey "$(jq -r '.secret' data/ibc1/n0/gaiacli/key_seed.json)"

# Then its time to initialize the relayer's lite clients for each chain
# All data moving forward is validated by these lite clients.
rly lite init ibc0 -f
rly lite init ibc1 -f

fi

if [[ -z $NORUN ]]; then
  echo "Starting relayer..."
  while :
  do
    # loop infinitely
  done
fi

