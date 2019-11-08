#!/usr/bin/env bash

echo "Clearing previous files..."
rm -rf ~/.mp*
mkdir -p ~/.mpd/config
cp config.toml ~/.mpd/config

echo "Building..."
make install

echo "Initialization..."

cur_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )

account_num=200
money_count="100000token"
file_output=$cur_path"/out.txt"
stake_count=100000000

while test $# -gt 0; do
  case "$1" in
    -h|--help)
      echo "run marketplace node"
      echo " "
      echo "run.sh [options]"
      echo " "
      echo "options:"
      echo "-h, --help                show brief help"
      echo "--demo                    set demo mode (add demo accounts)"
      echo "-n, --num_account=n       specify number of demo accounts | 200 default"
      echo "-m, --money=m             specify token amount for demo account | 100000token default"
      echo "-s, --stake=s             specify stake amount for demo account | 100000000 default"
      echo "-o, --output_file=o       specify output file | out.txt default"
      echo "--embeded                 set embeded mode (for docker)"
      exit 0
      ;;
    -n|--num_account)
      shift
      if test $# -gt 0; then
        account_num=$1
      else
        echo "no number of accounts specified"
        exit 1
      fi
      shift
      ;;
    -m|--money)
      shift
      if test $# -gt 0; then
        money_count=$1
      else
        echo "no money amount specified"
        exit 1
      fi
      shift
      ;;
    -s|--stake)
      shift
      if test $# -gt 0; then
        stake_count=$1
      else
        echo "no money amount specified"
        exit 1
      fi
      shift
      ;;
    -o|--output_file)
      shift
      if test $# -gt 0; then
        file_output=$1
      else
        echo "no output_file specified"
        exit 1
      fi
      shift
      ;;
    --embeded)
      EMBEDED=true
      shift
      ;;
    --demo)
      DEMO=true
      shift
      ;;
    *)
      break
      ;;
  esac
done

rm $file_output

mpd init node0 --chain-id mpchain

echo "Adding keys..."

mpcli keys add user1 <<< "12345678"
mpcli keys add user2 <<< "12345678"
mpcli keys add sellerBeneficiary <<< "12345678"
mpcli keys add buyerBeneficiary <<< "12345678"
mpcli keys add dgaming <<< "12345678"

mpd add-genesis-account $(mpcli keys show user1 -a) 1000token,100000000stake
mpd add-genesis-account $(mpcli keys show user2 -a) 1000token,100000000stake
mpd add-genesis-account $(mpcli keys show sellerBeneficiary -a) 1000token,100000000stake
mpd add-genesis-account $(mpcli keys show buyerBeneficiary -a) 1000token,100000000stake
mpd add-genesis-account $(mpcli keys show dgaming -a) 1000token,100000000stake


if [[ $DEMO ]]; then
  for ((i=1;i<=$account_num;i++));
  do
    pwd=$(gpg --gen-random --armor 1 14)
    mnemonic=$(mpcli keys add demo$i <<< $pwd 2>&1 >/dev/null | tail -1)

    mpd add-genesis-account $(mpcli keys show demo$i -a) $money_count,${stake_count}stake
    echo "demo$i      $pwd        $money_count   ${stake_count}stake       $mnemonic" >> $file_output
  done
fi


echo "Adding genesis account..."

echo "Configuring..."
mpcli config chain-id mpchain
mpcli config output json
mpcli config indent true
mpcli config trust-node true

mpd gentx --name user1 <<< "12345678"
mpd collect-gentxs
mpd validate-genesis

if [[ $EMBEDED ]]; then
  sed -i 's/proxy_app = "tcp:\/\/127.0.0.1:26658"/proxy_app = "tcp:\/\/0.0.0.0:26658"/' /root/.mpd/config/config.toml
  sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/' /root/.mpd/config/config.toml
else
  echo "Starting node..."

  mpd start
fi