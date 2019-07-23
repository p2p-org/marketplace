# DGaming Marketplace

## Init and run a node

Init node configuration with:

```bash 
mpd init [node_name] --chain-id [chain_id] 
```

To specify maximum beneficiary fee init with flag 'max-commission'. Example:

```bash
mpd init node0 --chain-id mpchain --max-commission 0.07
```

To run a node with a script:

```bash
./run.sh
```

This will start a node with two users, `user1` and `user2` (both are validators).

## Client commands
To get information about account:
```
$ mpcli query account $(mpcli keys show [name] -a)
```
Example:

```
$ mpcli query account $(mpcli keys show user1 -a)
```

To mint an NFT for that user:

```bash
mpcli tx marketplace mint $(uuidgen) name description image token_uri --from user1
```

The token is **not** put on the market when minted.

To transfer a token from user1 to user2:

```bash
mpcli tx marketplace transfer 686769b1-9395-4821-8a9e-36008ad4ca7c cosmos16y2vaas25ea8n353tfve45rwvt4sx0gl627pzn --from user1
```

To put a token on market (to make it purchasable by anybody who offers the exact price you specified):

```bash
mpcli tx marketplace put_on_market 686769b1-9395-4821-8a9e-36008ad4ca7c 150token cosmos16y2vaas25ea8n353tfve45rwvt4sx0gl627pzn --from user1
```

Note that you *must* provide the beneficiary address.

To buy a token:

```bash
mpcli tx marketplace buy 686769b1-9395-4821-8a9e-36008ad4ca7c cosmos16y2vaas25ea8n353tfve45rwvt4sx0gl627pzn --from user2
```

To buy a token with specified commission add 'beneficiary-commission' (--beneficiary-commission or -c) flag

```bash
mpcli tx marketplace buy 686769b1-9395-4821-8a9e-36008ad4ca7c cosmos16y2vaas25ea8n353tfve45rwvt4sx0gl627pzn -c 0.013 --from user2
```
or
```bash
mpcli tx marketplace buy 686769b1-9395-4821-8a9e-36008ad4ca7c cosmos16y2vaas25ea8n353tfve45rwvt4sx0gl627pzn --beneficiary-commission 0.013 --from user2
```

To create some number of fungible tokens:
```bash
mpcli tx marketplace createFT fungible 1000 --from user1
```

To transfer some amount of fungible tokens:
```bash
mpcli tx marketplace transferFT $(mpcli keys show user1 -a) fungible 500  --from user1
```

## Full scenario

After running `./run.sh`, 4 users are created: `user1` (minter and seller), `user2` (buyer), `sellerBeneficiary` and `buyerBeneficiary` (each has 1000token coins in the beginning).

Mint a new token:

```
$ mpcli tx marketplace mint $(uuidgen) name description image token_uri --from user1
```
*Output:*
```
{
  "chain_id": "mpchain",
  "account_number": "0",
  "sequence": "1",
  "fee": {
    "amount": [],
    "gas": "200000"
  },
  "msgs": [
    {
      "type": "marketplace/MintNFT",
      "value": {
        "owner": "cosmos1qv79nvxnkq7pf2tgrgjz53w9as6hlp7zszcpvr",
        "name": "name",
        "description": "description",
        "image": "image",
        "token_uri": "token_uri"
      }
    }
  ],
  "memo": ""
}

confirm transaction before signing and broadcasting [y/N]: y
Password to sign with 'user1':
{
  "height": "0",
  "txhash": "12AAB743F568E72E22E05C040AFA9CB5450C70FF709AFBFF1B51D6A8BDED2359"
}
```

Change token params:
Use 'update_params' with flags --image --price --description --token_uri --name
```
mpcli tx marketplace update_params 4eb281c9-1eea-4aab-b508-a3c27828b572 --from user1 -i newimage -p 500token -d newdescription -u newuri -n newname
```
List nfts:
```
$ mpcli query marketplace nfts
```
*Output:*
```
{
  "nfts": [
    {
      "nft": {
        "id": "4eb281c9-1eea-4aab-b508-a3c27828b572",
        "owner": "cosmos1qv79nvxnkq7pf2tgrgjz53w9as6hlp7zszcpvr",
        "name": "name",
        "description": "description",
        "image": "image",
        "token_uri": "token_uri"
      },
      "price": [],
      "on_sale": false,
      "seller_beneficiary": ""
    }
  ]
}
```

Put the new token on the market (and specify `sellerBeneficiary`):

```
$ mpcli tx marketplace put_on_market 4eb281c9-1eea-4aab-b508-a3c27828b572 650token cosmos1497eedaprzjvydwvgj5tu9e97agw30d7ksj99r --from user1
```
*Output:*
```
{
  "chain_id": "mpchain",
  "account_number": "0",
  "sequence": "2",
  "fee": {
    "amount": [],
    "gas": "200000"
  },
  "msgs": [
    {
      "type": "marketplace/PutOnMarketNFT",
      "value": {
        "owner": "cosmos1qv79nvxnkq7pf2tgrgjz53w9as6hlp7zszcpvr",
        "beneficiary": "cosmos1497eedaprzjvydwvgj5tu9e97agw30d7ksj99r",
        "token_id": "4eb281c9-1eea-4aab-b508-a3c27828b572",
        "price": [
          {
            "denom": "token",
            "amount": "650"
          }
        ]
      }
    }
  ],
  "memo": ""
}

confirm transaction before signing and broadcasting [y/N]: y
Password to sign with 'user1':
{
  "height": "0",
  "txhash": "B5CA8C210EECBD58E2B35EE5B7BD35BA0F64B032CC7D151997FE96B162CB932A"
}
```

Buy the token (and specify `buyerBeneficiary`):

```
$ mpcli tx marketplace buy 4eb281c9-1eea-4aab-b508-a3c27828b572 cosmos1qgq89a2xquyasydkyu6x7x96fq822z3em2t8xf --from user2
```
*Output:*
```
{
  "chain_id": "mpchain",
  "account_number": "1",
  "sequence": "0",
  "fee": {
    "amount": [],
    "gas": "200000"
  },
  "msgs": [
    {
      "type": "marketplace/BuyNFT",
      "value": {
        "buyer": "cosmos19608kpjnmmhzc2r9qp45eqd89m4c0z0wv7fy3j",
        "beneficiary": "cosmos1qgq89a2xquyasydkyu6x7x96fq822z3em2t8xf",
        "token_id": "4eb281c9-1eea-4aab-b508-a3c27828b572"
      }
    }
  ],
  "memo": ""
}

confirm transaction before signing and broadcasting [y/N]: y
Password to sign with 'user2':
{
  "height": "0",
  "txhash": "4A320FF2637F4274FDAE00F0D9CADFAF6772F223200D45FCA34D630E2B30A138"
}
```

After this we have token balances updated:

```
$ mpcli query account $(mpcli keys show buyerBeneficiary -a)
```
*Output:*
```
{
  "type": "auth/Account",
  "value": {
    "address": "cosmos1qgq89a2xquyasydkyu6x7x96fq822z3em2t8xf",
    "coins": [
      {
        "denom": "stake",
        "amount": "100000000"
      },
      {
        "denom": "token",
        "amount": "1004"
      }
    ],
    "public_key": null,
    "account_number": "3",
    "sequence": "0"
  }
}
```

```
$ mpcli query account $(mpcli keys show sellerBeneficiary -a)
```
*Output:*
```
{
  "type": "auth/Account",
  "value": {
    "address": "cosmos1497eedaprzjvydwvgj5tu9e97agw30d7ksj99r",
    "coins": [
      {
        "denom": "stake",
        "amount": "100000000"
      },
      {
        "denom": "token",
        "amount": "1004"
      }
    ],
    "public_key": null,
    "account_number": "2",
    "sequence": "0"
  }
}
```

```
$ mpcli query account $(mpcli keys show user1 -a)
```
*Output:*
```
{
  "type": "auth/Account",
  "value": {
    "address": "cosmos1qv79nvxnkq7pf2tgrgjz53w9as6hlp7zszcpvr",
    "coins": [
      {
        "denom": "token",
        "amount": "1634"
      }
    ],
    "public_key": {
      "type": "tendermint/PubKeySecp256k1",
      "value": "Ap0y3b8HOn1unrSvTOwSJ82ykJnqHE4RkL0Tj56d3mEX"
    },
    "account_number": "0",
    "sequence": "3"
  }
}
```

```
$ mpcli query account $(mpcli keys show user2 -a)
```
*Output:*
```
{
  "type": "auth/Account",
  "value": {
    "address": "cosmos19608kpjnmmhzc2r9qp45eqd89m4c0z0wv7fy3j",
    "coins": [
      {
        "denom": "stake",
        "amount": "100000000"
      },
      {
        "denom": "token",
        "amount": "349"
      }
    ],
    "public_key": {
      "type": "tendermint/PubKeySecp256k1",
      "value": "ApjMM44kZ8YAolktUY4Qj5nbwjGRsCcfhVtim/FM8rLs"
    },
    "account_number": "1",
    "sequence": "1"
  }
}
```