# REST

mpcli rest-server --chain-id mpchain --trust-node

##

curl -s http://localhost:1317/marketplace/nfts

curl -s http://localhost:1317/marketplace/nfts/cf9d19be-30f8-429e-9a75-8f997f524481

##

curl -XPUT -s http://localhost:1317/marketplace/mint --data-binary '{"base_req":{"from":"'$(mpcli keys show user1 -a)'","chain_id":"mpchain","sequence":"1","account_number":"0"},"token_name":"name","token_id":"'$(uuidgen)'","owner":"user1","name":"user1","password":"12345678","description":"desc","image":"ing","token_uri":"uri"}'

curl -XPUT -s http://localhost:1317/marketplace/put_on_market --data-binary '{"base_req":{"from":"'$(mpcli keys show user1 -a)'","chain_id":"mpchain","sequence":"2","account_number":"0"},"token_id":"cf9d19be-30f8-429e-9a75-8f997f524481","name":"user1","password":"12345678","beneficiary":"'$(mpcli keys show sellerBeneficiary -a)'","price":"650token"}'

curl -XPUT -s http://localhost:1317/marketplace/buy --data-binary '{"base_req":{"from":"'$(mpcli keys show user2 -a)'","chain_id":"mpchain","sequence":"0","account_number":"1"},"token_id":"cf9d19be-30f8-429e-9a75-8f997f524481","name":"user2","password":"12345678","beneficiary":"'$(mpcli keys show buyerBeneficiary -a)'"}'

curl -XPUT -s http://localhost:1317/marketplace/update_params --data-binary '{"base_req":{"from":"'$(mpcli keys show user1 -a)'","chain_id":"mpchain","sequence":"3","account_number":"0"},"token_name":"name","token_id":"cf9d19be-30f8-429e-9a75-8f997f524481","name":"user1","password":"12345678","image":"newimage"}'

curl -XPUT -s http://localhost:1317/marketplace/transfer --data-binary '{"base_req":{"from":"'$(mpcli keys show user1 -a)'","chain_id":"mpchain","sequence":"4","account_number":"0"},"token_id":"cf9d19be-30f8-429e-9a75-8f997f524481","name":"user1","password":"12345678","recipient":"'$(mpcli keys show user2 -a)'"}'

##

curl -XPUT -s http://localhost:1317/marketplace/create_ft --data-binary '{"base_req":{"from":"'$(mpcli keys show user1 -a)'","chain_id":"mpchain","sequence":"1","account_number":"0"},"denom":"pigs","amount":"100","name":"user1","password":"12345678"}'

curl -XPUT -s http://localhost:1317/marketplace/transfer_ft --data-binary '{"base_req":{"from":"'$(mpcli keys show user1 -a)'","chain_id":"mpchain","sequence":"2","account_number":"0"},"denom":"pigs","amount":"10","name":"user1","password":"12345678","recipient":"'$(mpcli keys show user2 -a)'"}'

##

curl -s http://localhost:1317/marketplace/fungible_tokens

curl -s http://localhost:1317/marketplace/fungible_tokens/token

