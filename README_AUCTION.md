# AUCTION

mpcli query account $(mpcli keys show user1 -a)
mpcli query account $(mpcli keys show user2 -a)
mpcli query account $(mpcli keys show sellerBeneficiary -a)
mpcli query account $(mpcli keys show buyerBeneficiary -a)

mpcli query marketplace auction_lots

mpcli query marketplace auction_lot [nft_id]

mpcli tx marketplace mint $(uuidgen) name description image token_uri --from user1

mpcli tx marketplace put_on_auction [nft_id] 100token [sellerBeneficiary] 5h -u 300token --from user1

mpcli tx marketplace bid [nft_id] [buyerBeneficiary] 100token --from user2

mpcli tx marketplace buyout [nft_id] [buyerBeneficiary] --from user2

mpcli tx remove_from_auction [nft_id] --from user1

mpcli tx finish_auction [nft_id] --from user1

