# DGaming Marketplace

To run a node:

```bash
./run.sh
```

This will start a node with two users, `user1` and `user2` (both are validators).

To mint an NFT for that user:

```bash
mpcli tx marketplace mint name description image token_uri 10token --from user1
```

Price must be specified; the token is **not** put on the market when minted.

To transfer a token from one user to another:

```bash
mpcli tx marketplace transfer user2 --from user1
```

