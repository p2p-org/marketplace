module github.com/dgamingfoundation/marketplace

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20191013030331-92ea174ea6e6
	github.com/cosmos/modules/incubator/nft v0.0.0-20191015123508-50d0c8092493
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.3
	github.com/magiconair/properties v1.8.1
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/common v0.4.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.32.6
	github.com/tendermint/tm-db v0.2.0
)

replace github.com/cosmos/cosmos-sdk => github.com/cosmos/cosmos-sdk v0.34.4-0.20191013030331-92ea174ea6e6
