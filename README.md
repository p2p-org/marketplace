# DGaming Marketplace

To run a node:

```bash
./run.sh
```

This will start a node with two users, `user1` and `user2` (both are validators).

To mint an NFT for that user:

```bash
mpcli tx marketplace mint name description image token_uri --from user1
```

The token is **not** put on the market when minted.

To transfer a token from user1 to user2:

```bash
mpcli tx marketplace transfer 20cec63d-bc88-44da-94c8-b67044ff7ab2 cosmos16y2vaas25ea8n353tfve45rwvt4sx0gl627pzn --from user1
```

To sell a token (to make it purchasable by anybody who offers the exact price you specified):

```bash
mpcli tx marketplace sell 20cec63d-bc88-44da-94c8-b67044ff7ab2 --from user1
```


