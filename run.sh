#!/usr/bin/env bash

account_num=${ACC_CNT:-5}
money_count="100000token"
file_output=${ACC_OUT_FILE:-"out.txt"}
file_input=${ACC_IN_FILE:-""}
stake_count=100000000
PSW="12345678"

while test $# -gt 0; do
  case "$1" in
    -h|--help)
      echo "run marketplace node"
      echo " "
      echo "run.sh [options]"
      echo " "
      echo "options:"
      echo "-h, --help                show brief help"
      echo "--build                   manual rebuild marketplace (not for docker)"
      echo "--reset                   force reset blockchain state (clear all data)"
      echo "--demo                    add demo accounts (used only with --reset)"
      echo "--input_mnemonic          specify mnemonics file"
      echo "-n, --num_account=n       specify number of demo accounts | 200 default"
      echo "-m, --money=m             specify token amount for demo account | 100000token default"
      echo "-s, --stake=s             specify stake amount for demo account | 100000000 default"
      echo "-o, --output_file=o       specify output file | out.txt default"
      echo "--norun                   omit service start"
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
    --build)
      BUILD=true
      shift
      ;;
    --demo)
      DEMO=true
      shift
      ;;
    --reset)
      RESET=true
      shift
      ;;
    --norun)
      NORUN=true
      shift
      ;;
    --input_mnemonic)
      shift
      if test $# -gt 0; then
        file_input=$1
      else
        echo "no mnemonic_file specified"
        exit 1
      fi
      shift
      ;;
    *)
      break
      ;;
  esac
done

if [[ $BUILD ]]; then
  echo "Building..."
  make install
fi

if [[ $RESET ]]; then
  echo "Clearing previous files..."
  rm -rf ~/.mp*
fi

if [ ! -f ~/.mpd/config/config.toml ]; then
  mkdir -p ~/.mpd/config
  cp config.toml ~/.mpd/config

  echo "Initialization..."
  mpd init node0 --chain-id mpchain
  mpcli config keyring-backend test

  echo "Adding genesis accounts..."
  mpcli keys add user1 -i < data_u1.txt
  mpcli keys add user2 -i < data_u2.txt
  mpcli keys add user3 -i < data_u3.txt
  mpcli keys add sellerBeneficiary -i < data_sb.txt
  mpcli keys add buyerBeneficiary -i < data_bb.txt
  mpcli keys add relay -i < data_dg.txt

  mpd add-genesis-account $(mpcli keys show user1 -a) 999999token,100000000stake
  mpd add-genesis-account $(mpcli keys show user2 -a) 999999token,100000000stake
  mpd add-genesis-account $(mpcli keys show user3 -a) 999999token,100000000stake
  mpd add-genesis-account $(mpcli keys show sellerBeneficiary -a) 100000000stake
  mpd add-genesis-account $(mpcli keys show buyerBeneficiary -a) 100000000stake
  mpd add-genesis-account $(mpcli keys show relay -a) 100000000stake

  if [[ $DEMO ]]; then
    if [[ -e $file_input ]]; then
      echo "read prepared mnemonics"
      i=1
      while read -r line; do
        echo $line > tmp.txt
        echo "" >> tmp.txt
#        echo $PSW >> tmp.txt
#        echo $PSW >> tmp.txt
#        echo "" >> tmp.txt

        echo "gen user demo${i}"
        mpcli keys add demo$i -i < tmp.txt
        mpd add-genesis-account $(mpcli keys show demo$i -a) $money_count --keyring-backend test
        i=$((i+1))
      done < $file_input
      rm tmp.txt
    else
      rm $file_output
      echo "generate mnemonics"
      for ((i=1;i<=$account_num;i++));
      do
        mnemonic=$(mpcli keys add demo$i |& tail -1)
        mpd add-genesis-account $(mpcli keys show demo$i -a) $money_count --keyring-backend test
        echo "demo$i      $pwd        $money_count       $mnemonic" >> $file_output
      done
    fi
  fi

  echo "Configuring..."
  mpcli config chain-id mpchain
  mpcli config output json
  mpcli config indent true
  mpcli config trust-node true

  mpd gentx --name user1 --keyring-backend test
  mpd collect-gentxs
  mpd validate-genesis
fi

if [[ -z $NORUN ]]; then
  echo "Starting node..."
  #mpcli rest-server --chain-id mpchain --trust-node --laddr tcp://0.0.0.0:1317 > /dev/null &
  mpd start #&> mplog.log
fi
