package app

import (
	"encoding/json"
	"fmt"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	tmos "github.com/tendermint/tendermint/libs/os"
	"io"
	"log"
	"os"

	tlog "github.com/tendermint/tendermint/libs/log"

	"github.com/corestario/marketplace/common"
	"github.com/corestario/marketplace/x/marketplace"
	"github.com/corestario/marketplace/x/marketplace/config"
	transferNFT "github.com/corestario/marketplace/x/nftIBC"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	transfer "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	"github.com/cosmos/modules/incubator/nft"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const appName = "marketplace"

var (
	// default home directories for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.mpcli")

	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.mpd")

	// ModuleBasicManager is in charge of setting up basic module elemnets
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		params.AppModuleBasic{},
		mint.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		nft.AppModuleBasic{},
		ibc.AppModuleBasic{},
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler, distr.ProposalHandler, upgradeclient.ProposalHandler,
		),
		upgrade.AppModuleBasic{},

		marketplace.AppModule{},
		transfer.AppModuleBasic{},
		transferNFT.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:              nil,
		distr.ModuleName:                   nil,
		mint.ModuleName:                    {auth.Minter},
		staking.BondedPoolName:             {auth.Burner, auth.Staking},
		staking.NotBondedPoolName:          {auth.Burner, auth.Staking},
		gov.ModuleName:                     {auth.Burner},
		transfer.GetModuleAccountName():    {auth.Minter, auth.Burner},
		transferNFT.GetModuleAccountName(): {auth.Minter, auth.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		distr.ModuleName: true,
	}
)

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

type marketplaceApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// subspaces
	subspaces map[string]params.Subspace

	// Keepers
	accountKeeper     auth.AccountKeeper
	bankKeeper        bank.Keeper
	mintKeeper        mint.Keeper
	stakingKeeper     staking.Keeper
	crisisKeeper      crisis.Keeper
	govKeeper         gov.Keeper
	slashingKeeper    slashing.Keeper
	distrKeeper       distr.Keeper
	paramsKeeper      params.Keeper
	nftKeeper         *nft.Keeper
	ibcKeeper         *ibc.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	evidenceKeeper    evidence.Keeper
	transferKeeper    transfer.Keeper
	transferNFTKeeper transferNFT.Keeper
	capabilityKeeper  *capability.Keeper
	upgradeKeeper     upgrade.Keeper

	mpKeeper *marketplace.Keeper

	// make scoped keepers public for test purposes
	scopedIBCKeeper         capability.ScopedKeeper
	scopedTransferKeeper    capability.ScopedKeeper
	scopedTransferNFTKeeper capability.ScopedKeeper

	// Module Manager
	mm *module.Manager
}

// NewMarketplaceApp is a const	ructor function for marketplaceApp
func NewMarketplaceApp(logger tlog.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string,
	baseAppOptions ...func(*bam.BaseApp)) *marketplaceApp {

	// First define the top level codec that will be shared by the different modules
	//cdc := std.MakeCodec(ModuleBasics)
	//cdc := std.MakeCodec(ModuleBasics)
	//appCodec := std.NewAppCodec(cdc)
	appCodec, cdc := MakeCodecs()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)

	keys := sdk.NewKVStoreKeys(
		auth.StoreKey, bank.StoreKey, staking.StoreKey,
		mint.StoreKey, distr.StoreKey, slashing.StoreKey,
		gov.StoreKey, params.StoreKey, ibc.StoreKey, upgrade.StoreKey,
		evidence.StoreKey, transfer.StoreKey, capability.StoreKey,
		nft.StoreKey, marketplace.StoreKey, marketplace.RegisterCurrencyKey, marketplace.AuctionKey,
		marketplace.DeletedNFTKey, transferNFT.StoreKey,
	)

	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capability.MemStoreKey)

	// Here you initialize your application with the store keys it requires
	var app = &marketplaceApp{
		BaseApp:   bApp,
		cdc:       cdc,
		keys:      keys,
		tkeys:     tkeys,
		memKeys:   memKeys,
		subspaces: make(map[string]params.Subspace),
	}

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(appCodec, keys[params.StoreKey], tkeys[params.TStoreKey])
	// Set specific supspaces
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[mint.ModuleName] = app.paramsKeeper.Subspace(mint.DefaultParamspace)
	app.subspaces[distr.ModuleName] = app.paramsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	app.subspaces[crisis.ModuleName] = app.paramsKeeper.Subspace(crisis.DefaultParamspace)

	bApp.SetParamStore(app.paramsKeeper.Subspace(bam.Paramspace).WithKeyTable(std.ConsensusParamsKeyTable()))

	app.capabilityKeeper = capability.NewKeeper(appCodec, keys[capability.StoreKey], memKeys[capability.MemStoreKey])
	scopedIBCKeeper := app.capabilityKeeper.ScopeToModule(ibc.ModuleName)
	scopedTransferKeeper := app.capabilityKeeper.ScopeToModule(transfer.ModuleName)
	scopedTransferNFTKeeper := app.capabilityKeeper.ScopeToModule(transferNFT.ModuleName)

	app.accountKeeper = auth.NewAccountKeeper(
		appCodec, keys[auth.StoreKey], app.subspaces[auth.ModuleName], auth.ProtoBaseAccount, maccPerms,
	)
	app.bankKeeper = bank.NewBaseKeeper(
		appCodec, keys[bank.StoreKey], app.accountKeeper, app.subspaces[bank.ModuleName], app.BlacklistedAccAddrs(),
	)
	stakingKeeper := staking.NewKeeper(
		appCodec, keys[staking.StoreKey], app.accountKeeper, app.bankKeeper, app.subspaces[staking.ModuleName],
	)
	app.mintKeeper = mint.NewKeeper(
		appCodec, keys[mint.StoreKey], app.subspaces[mint.ModuleName], &stakingKeeper,
		app.accountKeeper, app.bankKeeper, auth.FeeCollectorName,
	)
	app.distrKeeper = distr.NewKeeper(
		appCodec, keys[distr.StoreKey], app.subspaces[distr.ModuleName], app.accountKeeper, app.bankKeeper,
		&stakingKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	app.slashingKeeper = slashing.NewKeeper(
		appCodec, keys[slashing.StoreKey], &stakingKeeper, app.subspaces[slashing.ModuleName],
	)
	app.crisisKeeper = crisis.NewKeeper(
		app.subspaces[crisis.ModuleName], invCheckPeriod, app.bankKeeper, auth.FeeCollectorName,
	)
	app.upgradeKeeper = upgrade.NewKeeper(skipUpgradeHeights, keys[upgrade.StoreKey], appCodec, DefaultNodeHome)

	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper)).
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.upgradeKeeper))
	app.govKeeper = gov.NewKeeper(
		appCodec, keys[gov.StoreKey], app.subspaces[gov.ModuleName], app.accountKeeper, app.bankKeeper,
		&stakingKeeper, govRouter,
	)

	// The staking keeper
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	// The NFTKeeper is the Keeper from the module NFTs.
	newKeeper := nft.NewKeeper(
		app.cdc,
		app.keys[nft.StoreKey],
	)
	app.nftKeeper = &newKeeper

	nftModule := nft.NewAppModule(newKeeper, app.accountKeeper)

	srvCfg := ReadSrvConfig()
	fmt.Printf("Server Config: \n %+v \n", srvCfg)

	app.ibcKeeper = ibc.NewKeeper(
		app.cdc, appCodec, keys[ibc.StoreKey], app.stakingKeeper, scopedIBCKeeper,
	)

	app.mpKeeper = marketplace.NewKeeper(
		app.bankKeeper,
		app.stakingKeeper,
		app.distrKeeper,
		app.keys[marketplace.StoreKey],
		app.keys[marketplace.RegisterCurrencyKey],
		app.keys[marketplace.AuctionKey],
		app.keys[marketplace.DeletedNFTKey],
		app.cdc,
		srvCfg,
		common.NewPrometheusMsgMetrics("marketplace"),
		app.nftKeeper,
		&app.mintKeeper,
		&app.accountKeeper,
		app.ibcKeeper,
	)

	// Create Transfer Keepers
	app.transferKeeper = transfer.NewKeeper(
		appCodec, keys[transfer.StoreKey],
		app.ibcKeeper.ChannelKeeper, &app.ibcKeeper.PortKeeper,
		app.accountKeeper, app.bankKeeper,
		scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.transferKeeper)

	app.transferNFTKeeper = transferNFT.NewKeeper(
		appCodec, keys[transferNFT.StoreKey],
		app.ibcKeeper.ChannelKeeper, &app.ibcKeeper.PortKeeper,
		app.accountKeeper, app.bankKeeper,
		scopedTransferNFTKeeper, app.mpKeeper, app.nftKeeper,
	)
	transferNFTModule := transferNFT.NewAppModule(app.transferNFTKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := port.NewRouter()
	ibcRouter.AddRoute(transfer.ModuleName, transferModule)
	ibcRouter.AddRoute(transferNFT.ModuleName, transferNFTModule)
	app.ibcKeeper.SetRouter(ibcRouter)

	overriddenNFTModule := marketplace.NewNFTModuleMarketplace(nftModule, app.nftKeeper, app.mpKeeper)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(appCodec, app.accountKeeper),
		bank.NewAppModule(appCodec, app.bankKeeper, app.accountKeeper),
		capability.NewAppModule(appCodec, *app.capabilityKeeper),
		mint.NewAppModule(appCodec, app.mintKeeper, app.accountKeeper),
		slashing.NewAppModule(appCodec, app.slashingKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		distr.NewAppModule(appCodec, app.distrKeeper, app.accountKeeper, app.bankKeeper, app.stakingKeeper),
		staking.NewAppModule(appCodec, app.stakingKeeper, app.accountKeeper, app.bankKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		gov.NewAppModule(appCodec, app.govKeeper, app.accountKeeper, app.bankKeeper),
		upgrade.NewAppModule(app.upgradeKeeper),
		evidence.NewAppModule(app.evidenceKeeper),
		ibc.NewAppModule(app.ibcKeeper),
		params.NewAppModule(app.paramsKeeper),
		mint.NewAppModule(appCodec, app.mintKeeper, app.accountKeeper),
		transferModule,
		transferNFTModule,
		marketplace.NewAppModule(app.mpKeeper, app.bankKeeper, app.nftKeeper),
		overriddenNFTModule,
	)

	app.mm.SetOrderBeginBlockers(
		upgrade.ModuleName, mint.ModuleName, distr.ModuleName, slashing.ModuleName,
		evidence.ModuleName, staking.ModuleName, ibc.ModuleName,
	)
	app.mm.SetOrderEndBlockers(crisis.ModuleName, gov.ModuleName, staking.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	app.mm.SetOrderInitGenesis(
		capability.ModuleName,
		auth.ModuleName,
		distr.ModuleName,
		staking.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		gov.ModuleName,
		mint.ModuleName,
		crisis.ModuleName,
		nft.ModuleName,
		ibc.ModuleName,
		genutil.ModuleName,
		transfer.ModuleName,
		transferNFT.ModuleName,
		marketplace.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)

	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(
		ante.NewAnteHandler(
			app.accountKeeper, app.bankKeeper, *app.ibcKeeper,
			ante.DefaultSigVerificationGasConsumer,
		),
	)

	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	if err := app.LoadLatestVersion(); err != nil {
		tmos.Exit(err.Error())
	}

	// Initialize and seal the capability keeper so all persistent capabilities
	// are loaded in-memory and prevent any further modules from creating scoped
	// sub-keepers.
	// This must be done during creation of baseapp rather than in InitChain so
	// that in-memory capabilities get regenerated on app restart
	ctx := app.BaseApp.NewContext(true, abci.Header{Height: -1})
	app.capabilityKeeper.InitializeAndSeal(ctx)

	app.scopedIBCKeeper = scopedIBCKeeper
	app.scopedTransferKeeper = scopedTransferKeeper
	app.scopedTransferNFTKeeper = scopedTransferNFTKeeper

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
	cdc := std.MakeCodec(ModuleBasics)
	return ModuleBasics.DefaultGenesis(cdc)
}

func (app *marketplaceApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	// VERY important. registers default denoms (e.g. "token")
	app.mpKeeper.RegisterBasicDenoms(ctx)

	err := app.cdc.UnmarshalJSON(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}

	return app.mm.InitGenesis(ctx, app.cdc, genesisState)
}

func (app *marketplaceApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *marketplaceApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	app.mpKeeper.CheckFinishedAuctions(ctx)
	return app.mm.EndBlock(ctx, req)
}
func (app *marketplaceApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *marketplaceApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[auth.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlacklistedAccAddrs returns all the app's module account addresses black listed for receiving tokens.
func (app *marketplaceApp) BlacklistedAccAddrs() map[string]bool {
	blacklistedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blacklistedAddrs[auth.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	return blacklistedAddrs
}

//_________________________________________________________

// ExportAppStateAndValidators export the state of gaia for a genesis file
func (app *marketplaceApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string,
) (appState json.RawMessage, validators []tmtypes.GenesisValidator, cp *abci.ConsensusParams, err error) {
	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	if forZeroHeight {
		app.prepForZeroHeightGenesis(ctx, jailWhiteList)
	}

	genState := app.mm.ExportGenesis(ctx, app.cdc)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, nil, err
	}
	validators = staking.WriteValidators(ctx, app.stakingKeeper)
	return appState, validators, app.BaseApp.GetConsensusParams(ctx), nil
}

func (app *marketplaceApp) prepForZeroHeightGenesis(ctx sdk.Context, jailWhiteList []string) {
	applyWhiteList := false

	//Check if there is a whitelist
	if len(jailWhiteList) > 0 {
		applyWhiteList = true
	}

	whiteListMap := make(map[string]bool)

	for _, addr := range jailWhiteList {
		_, err := sdk.ValAddressFromBech32(addr)
		if err != nil {
			log.Fatal(err)
		}
		whiteListMap[addr] = true
	}

	/* Just to be safe, assert the invariants on current state. */
	app.crisisKeeper.AssertInvariants(ctx)

	/* Handle fee distribution state. */

	// withdraw all validator commission
	app.stakingKeeper.IterateValidators(ctx, func(_ int64, val staking.ValidatorI) (stop bool) {
		_, err := app.distrKeeper.WithdrawValidatorCommission(ctx, val.GetOperator())
		if err != nil {
			log.Fatal(err)
		}
		return false
	})

	// withdraw all delegator rewards
	dels := app.stakingKeeper.GetAllDelegations(ctx)
	for _, delegation := range dels {
		_, err := app.distrKeeper.WithdrawDelegationRewards(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
		if err != nil {
			log.Fatal(err)
		}
	}

	// clear validator slash events
	app.distrKeeper.DeleteAllValidatorSlashEvents(ctx)

	// clear validator historical rewards
	app.distrKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	// set context height to zero
	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	// reinitialize all validators
	app.stakingKeeper.IterateValidators(ctx, func(_ int64, val staking.ValidatorI) (stop bool) {

		// donate any unwithdrawn outstanding reward fraction tokens to the community pool
		scraps := app.distrKeeper.GetValidatorOutstandingRewards(ctx, val.GetOperator()).Rewards
		feePool := app.distrKeeper.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
		app.distrKeeper.SetFeePool(ctx, feePool)

		app.distrKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
		return false
	})

	// reinitialize all delegations
	for _, del := range dels {
		app.distrKeeper.Hooks().BeforeDelegationCreated(ctx, del.DelegatorAddress, del.ValidatorAddress)
		app.distrKeeper.Hooks().AfterDelegationModified(ctx, del.DelegatorAddress, del.ValidatorAddress)
	}

	// reset context height
	ctx = ctx.WithBlockHeight(height)

	/* Handle staking state. */

	// iterate through redelegations, reset creation height
	app.stakingKeeper.IterateRedelegations(ctx, func(_ int64, red staking.Redelegation) (stop bool) {
		for i := range red.Entries {
			red.Entries[i].CreationHeight = 0
		}
		app.stakingKeeper.SetRedelegation(ctx, red)
		return false
	})

	// iterate through unbonding delegations, reset creation height
	app.stakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd staking.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].CreationHeight = 0
		}
		app.stakingKeeper.SetUnbondingDelegation(ctx, ubd)
		return false
	})

	// Iterate through validators by power descending, reset bond heights, and
	// update bond intra-tx counters.
	store := ctx.KVStore(app.keys[staking.StoreKey])
	iter := sdk.KVStoreReversePrefixIterator(store, staking.ValidatorsKey)
	counter := int16(0)

	for ; iter.Valid(); iter.Next() {
		addr := sdk.ValAddress(iter.Key()[1:])
		validator, found := app.stakingKeeper.GetValidator(ctx, addr)
		if !found {
			panic("expected validator, not found")
		}

		validator.UnbondingHeight = 0
		if applyWhiteList && !whiteListMap[addr.String()] {
			validator.Jailed = true
		}

		app.stakingKeeper.SetValidator(ctx, validator)
		counter++
	}

	iter.Close()

	_ = app.stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)

	/* Handle slashing state. */

	// reset start height on signing infos
	app.slashingKeeper.IterateValidatorSigningInfos(
		ctx,
		func(addr sdk.ConsAddress, info slashing.ValidatorSigningInfo) (stop bool) {
			info.StartHeight = 0
			app.slashingKeeper.SetValidatorSigningInfo(ctx, addr, info)
			return false
		},
	)
}

func MakeCodecs() (*std.Codec, *codec.Codec) {
	cdc := std.MakeCodec(ModuleBasics)
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	appCodec := std.NewAppCodec(cdc, interfaceRegistry)

	sdk.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterInterfaceModules(interfaceRegistry)

	return appCodec, cdc
}

func ReadSrvConfig() *config.MPServerConfig {
	var cfg *config.MPServerConfig
	vCfg := viper.New()
	vCfg.SetConfigName("server")
	vCfg.AddConfigPath(DefaultNodeHome + "/config")
	err := vCfg.ReadInConfig()
	if err != nil {
		fmt.Println("ERROR: server config file not found, error:", err)
		return config.DefaultMPServerConfig()
	}
	fmt.Println(vCfg.GetString("maximum_beneficiary_commission"))
	err = vCfg.Unmarshal(&cfg)
	if err != nil {
		fmt.Println("ERROR: could not unmarshal server config file, error:", err)
		return config.DefaultMPServerConfig()
	}
	return cfg
}
