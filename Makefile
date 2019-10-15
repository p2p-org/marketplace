include Makefile.ledger

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION := 0.1.0
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/dgamingfoundation/marketplace | sed  's/ /\@/g')

export GO111MODULE = on


build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/dgamingfoundation/marketplace/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/dgamingfoundation/marketplace/version.Name=mp \
		  -X github.com/dgamingfoundation/marketplace/version.ServerName=mpd \
		  -X github.com/dgamingfoundation/marketplace/version.ClientName=mpcli \
		  -X github.com/dgamingfoundation/marketplace/version.Version=$(VERSION) \
		  -X github.com/dgamingfoundation/marketplace/version.Commit=$(COMMIT) \
		  -X "github.com/dgamingfoundation/marketplace/version.BuildTags=$(build_tags_comma_sep)"

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/dgamingfoundation/marketplace/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))


BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

#include contrib/devtools/Makefile

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
	go build -mod=readonly $(BUILD_FLAGS) -o build/mpd ./cmd/mpd
	go build -mod=readonly $(BUILD_FLAGS) -o build/mpcli ./cmd/mpcli

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