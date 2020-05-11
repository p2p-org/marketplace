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
