include Makefile.ledger

all: lint install

install: go.sum
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/mpd
		go install -mod=readonly $(BUILD_FLAGS) ./cmd/mpcli

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

test:
	@echo "Start testing:"
	go test ./...

build: go.sum
	go build $(BUILD_FLAGS) -o build/mpd ./cmd/mpd
	go build $(BUILD_FLAGS) -o build/mpcli ./cmd/mpcli

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

build-docker-mpdnode:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: build-linux localnet-stop
	@if ! [ -f build/node0/mpd/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/mpd:Z tendermint/mpdnode testnet --v 4 -o . --starting-ip-address 192.168.10.2 ;	fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down
.PHONY: test

###############################################################################
###                                Protobuf                                 ###
###############################################################################

proto-all: proto-tools proto-gen proto-lint proto-check-breaking

proto-gen:
	@./scripts/protocgen.sh

# This generates the SDK's custom wrapper for google.protobuf.Any. It should only be run manually when needed
proto-gen-any:
	@./scripts/protocgen-any.sh

proto-lint:
	@buf check lint --error-format=json

proto-check-breaking:
	@buf check breaking --against-input '.git#branch=master'

proto-lint-docker:
	@$(DOCKER_BUF) check lint --error-format=json
.PHONY: proto-lint

proto-check-breaking-docker:
	@$(DOCKER_BUF) check breaking --against-input $(HTTPS_GIT)#branch=master
.PHONY: proto-check-breaking-ci

TM_URL           = https://raw.githubusercontent.com/tendermint/tendermint/v0.33.1
GOGO_PROTO_URL   = https://raw.githubusercontent.com/regen-network/protobuf/cosmos
COSMOS_PROTO_URL = https://raw.githubusercontent.com/regen-network/cosmos-proto/master

TM_KV_TYPES         = third_party/proto/tendermint/libs/kv
TM_MERKLE_TYPES     = third_party/proto/tendermint/crypto/merkle
TM_ABCI_TYPES       = third_party/proto/tendermint/abci/types
GOGO_PROTO_TYPES    = third_party/proto/gogoproto
COSMOS_PROTO_TYPES  = third_party/proto/cosmos-proto
SDK_PROTO_TYPES     = third_party/proto/cosmos-sdk/types
AUTH_PROTO_TYPES    = third_party/proto/cosmos-sdk/x/auth/types
VESTING_PROTO_TYPES = third_party/proto/cosmos-sdk/x/auth/vesting/types
SUPPLY_PROTO_TYPES  = third_party/proto/cosmos-sdk/x/supply/types

proto-update-deps:
	@mkdir -p $(GOGO_PROTO_TYPES)
	@curl -sSL $(GOGO_PROTO_URL)/gogoproto/gogo.proto > $(GOGO_PROTO_TYPES)/gogo.proto

	@mkdir -p $(COSMOS_PROTO_TYPES)
	@curl -sSL $(COSMOS_PROTO_URL)/cosmos.proto > $(COSMOS_PROTO_TYPES)/cosmos.proto

	@mkdir -p $(TM_ABCI_TYPES)
	@curl -sSL $(TM_URL)/abci/types/types.proto > $(TM_ABCI_TYPES)/types.proto
	@sed -i '' '8 s|crypto/merkle/merkle.proto|third_party/proto/tendermint/crypto/merkle/merkle.proto|g' $(TM_ABCI_TYPES)/types.proto
	@sed -i '' '9 s|libs/kv/types.proto|third_party/proto/tendermint/libs/kv/types.proto|g' $(TM_ABCI_TYPES)/types.proto

	@mkdir -p $(TM_KV_TYPES)
	@curl -sSL $(TM_URL)/libs/kv/types.proto > $(TM_KV_TYPES)/types.proto

	@mkdir -p $(TM_MERKLE_TYPES)
	@curl -sSL $(TM_URL)/crypto/merkle/merkle.proto > $(TM_MERKLE_TYPES)/merkle.proto


.PHONY: proto-all proto-gen proto-lint proto-check-breaking proto-update-deps


proto-tools: proto-tools-stamp
proto-tools-stamp:
	@echo "Installing protoc compiler..."
	@(cd /tmp; \
	curl -OL "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/${PROTOC_ZIP}"; \
	unzip -o ${PROTOC_ZIP} -d $(PREFIX) bin/protoc; \
	unzip -o ${PROTOC_ZIP} -d $(PREFIX) 'include/*'; \
	rm -f ${PROTOC_ZIP})

	@echo "Installing protoc-gen-buf-check-breaking..."
	@curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/protoc-gen-buf-check-breaking-${UNAME_S}-${UNAME_M}" \
    -o "${BIN}/protoc-gen-buf-check-breaking" && \
	chmod +x "${BIN}/protoc-gen-buf-check-breaking"

	@echo "Installing buf..."
	@curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-${UNAME_S}-${UNAME_M}" \
    -o "${BIN}/buf" && \
	chmod +x "${BIN}/buf"

	touch $@

protoc-gen-gocosmos:
	@echo "Installing protoc-gen-gocosmos..."
	@go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos

tools-clean:
	rm -f $(STATIK) $(GOLANGCI_LINT) $(RUNSIM)
	rm -f tools-stamp proto-tools-stamp

.PHONY: tools-clean statik runsim \
	protoc-gen-gocosmos