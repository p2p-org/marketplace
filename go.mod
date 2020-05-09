module github.com/corestario/marketplace

go 1.12

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20200430150743-930802e7a13c
	github.com/cosmos/modules/incubator/nft v0.0.0-20191015123508-50d0c8092493
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.4
	github.com/magiconair/properties v1.8.1
	github.com/prometheus/client_golang v1.5.1
	github.com/prometheus/common v0.9.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.4
	github.com/tendermint/tm-db v0.5.1
)

replace github.com/cosmos/cosmos-sdk => github.com/cosmos/cosmos-sdk v0.34.4-0.20200430150743-930802e7a13c
