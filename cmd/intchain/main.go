package main

import (
	"fmt"
	"github.com/intfoundation/intchain/bridge"
	"github.com/intfoundation/intchain/chain"
	"github.com/intfoundation/intchain/cmd/geth"
	"github.com/intfoundation/intchain/cmd/utils"
	"github.com/intfoundation/intchain/console"
	"github.com/intfoundation/intchain/log"
	"github.com/intfoundation/intchain/metrics"
	"github.com/intfoundation/intchain/version"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

const (
	clientIdentifier = "intchain" // Client identifier to advertise over the network; it also is the main chain's id
)

func main() {

	cliApp := newCliApp(version.Version, "the intchain command line interface")
	cliApp.Action = intchainCmd
	cliApp.Commands = []cli.Command{

		{
			Action:      versionCmd,
			Name:        "version",
			Usage:       "",
			Description: "Print the version",
		},

		{
			Action:      chain.InitIntGenesis,
			Name:        "init_int_genesis",
			Usage:       "init_int_genesis balance:{\"1000000000000000000000000000\",\"100000000000000000000000\"}",
			Description: "Initialize the balance of accounts",
		},

		{
			Action:      chain.InitCmd,
			Name:        "init",
			Usage:       "init genesis.json",
			Description: "Initialize the files",
		},

		{
			Action:      chain.InitChildChainCmd,
			Name:        "init_child_chain",
			Usage:       "./intchain --datadir=~/.intchain --childChain=child_0,child_1,child_2 init_child_chain",
			Description: "Initialize child chain genesis from chain info db",
		},

		{
			//Action: GeneratePrivateValidatorCmd,
			Action: utils.MigrateFlags(GeneratePrivateValidatorCmd),
			Name:   "gen_priv_validator",
			Usage:  "gen_priv_validator address", //generate priv_validator.json for address
			Flags: []cli.Flag{
				utils.DataDirFlag,
			},
			Description: "Generate priv_validator.json for address",
		},

		//gethmain.ConsoleCommand,
		gethmain.AttachCommand,
		//gethmain.JavascriptCommand,
		gethmain.ImportChainCommand,
		gethmain.ExportChainCommand,
		gethmain.CountBlockStateCommand,

		//walletCommand,
		accountCommand,
	}
	cliApp.HideVersion = true // we have a command to print the version

	cliApp.Before = func(ctx *cli.Context) error {

		// Log Folder
		logFolderFlag := ctx.GlobalString(LogDirFlag.Name)

		// Setup the Global Logger
		commonLogDir := path.Join(ctx.GlobalString("datadir"), logFolderFlag, "common")
		log.NewLogger("", commonLogDir, ctx.GlobalInt(verbosityFlag.Name), ctx.GlobalBool(debugFlag.Name), ctx.GlobalString(vmoduleFlag.Name), ctx.GlobalString(backtraceAtFlag.Name))

		runtime.GOMAXPROCS(runtime.NumCPU())

		if err := bridge.Debug_Setup(ctx, logFolderFlag); err != nil {
			return err
		}

		// Start system runtime metrics collection
		go metrics.CollectProcessMetrics(3 * time.Second)

		return nil
	}

	cliApp.After = func(ctx *cli.Context) error {
		bridge.Debug_Exit()
		console.Stdin.Close() // Resets terminal mode.
		return nil
	}

	if err := cliApp.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newCliApp(version, usage string) *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Author = ""
	//app.Authors = nil
	app.Email = ""
	app.Version = version
	app.Usage = usage
	app.Flags = []cli.Flag{
		utils.IdentityFlag,
		//utils.UnlockedAccountFlag,
		utils.PasswordFileFlag,
		utils.BootnodesFlag,
		utils.BootnodesV4Flag,
		utils.BootnodesV5Flag,
		utils.DataDirFlag,
		utils.KeyStoreDirFlag,
		utils.NoUSBFlag,
		utils.TxPoolNoLocalsFlag,
		utils.TxPoolJournalFlag,
		utils.TxPoolRejournalFlag,
		utils.TxPoolPriceLimitFlag,
		utils.TxPoolPriceBumpFlag,
		utils.TxPoolAccountSlotsFlag,
		utils.TxPoolGlobalSlotsFlag,
		utils.TxPoolAccountQueueFlag,
		utils.TxPoolGlobalQueueFlag,
		utils.TxPoolLifetimeFlag,
		utils.SyncModeFlag,
		utils.GCModeFlag,
		utils.CacheFlag,
		utils.CacheDatabaseFlag,
		utils.CacheTrieFlag,
		utils.CacheGCFlag,
		utils.ListenPortFlag,
		utils.MaxPeersFlag,
		utils.MaxPendingPeersFlag,
		utils.MinerThreadsFlag,
		utils.MinerGasTargetFlag,
		utils.MinerGasLimitFlag,
		utils.MinerGasPriceFlag,
		utils.MinerCoinbaseFlag,
		utils.NATFlag,
		utils.NoDiscoverFlag,
		utils.DiscoveryV5Flag,
		utils.NetrestrictFlag,
		utils.NodeKeyFileFlag,
		utils.NodeKeyHexFlag,
		utils.TestnetFlag,
		utils.VMEnableDebugFlag,
		utils.NetworkIdFlag,
		utils.PruneFlag,
		//utils.PruneBlockFlag,

		utils.EthStatsURLFlag,
		utils.MetricsEnabledFlag,
		utils.NoCompactionFlag,
		utils.GpoBlocksFlag,
		utils.GpoPercentileFlag,
		utils.ExtraDataFlag,
		//gethmain.ConfigFileFlag,
		// RPC HTTP Flag
		utils.RPCEnabledFlag,
		utils.RPCListenAddrFlag,
		utils.RPCPortFlag,
		utils.RPCApiFlag,
		utils.RPCCORSDomainFlag,
		utils.RPCVirtualHostsFlag,
		// RPC WS Flag
		utils.WSEnabledFlag,
		utils.WSListenAddrFlag,
		utils.WSPortFlag,
		utils.WSApiFlag,
		utils.WSAllowedOriginsFlag,

		utils.IPCDisabledFlag,
		utils.IPCPathFlag,

		utils.SolcPathFlag,

		utils.PerfTestFlag,

		LogDirFlag,
		ChildChainFlag,

		/*
			//Tendermint flags
			MonikerFlag,
			NodeLaddrFlag,
			SeedsFlag,
			FastSyncFlag,
			SkipUpnpFlag,
			RpcLaddrFlag,
			AddrFlag,
		*/
	}
	app.Flags = append(app.Flags, DebugFlags...)

	return app
}

func versionCmd(ctx *cli.Context) error {
	fmt.Println(version.Version)
	return nil
}
