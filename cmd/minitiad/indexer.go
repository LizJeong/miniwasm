package main

import (
	"context"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/initia-labs/indexer"
	indexercfg "github.com/initia-labs/indexer/config"
	"github.com/initia-labs/indexer/service/collector"
	"github.com/initia-labs/indexer/service/cron/pair"
	"github.com/initia-labs/indexer/service/dashboard"
	minitiaapp "github.com/initia-labs/miniwasm/app"
)

const (
	FlagIndexer = "indexer"
)

func addIndexFlag(cmd *cobra.Command) {
	indexercfg.AddIndexerFlag(cmd)
}

func preSetupIndexer(svrCtx *server.Context, clientCtx client.Context, ctx context.Context, g *errgroup.Group, _app types.Application) error {
	app := _app.(*minitiaapp.MinitiaApp)

	// if indexer is disabled, it returns nil
	idxer, err := indexer.NewIndexer(svrCtx.Viper, app)
	if err != nil {
		return err
	}
	// if idxer is nil, it means indexer is disabled
	if idxer == nil {
		return nil
	}

	idxer.GetHandler().RegisterService(collector.CollectorSvc)
	idxer.GetHandler().RegisterService(dashboard.DashboardSvc)
	idxer.GetHandler().RegisterCronjob(pair.Tag, pair.Expr, pair.JobInit, pair.JobFunc)

	err = idxer.Validate()
	if err != nil {
		return err
	}

	// listen all keys
	keysToListen := []storetypes.StoreKey{}

	// TODO: if it downgrades performacne, have to set only keys for registered submodules and crons
	keys := app.GetKeys()
	for _, key := range keys {
		keysToListen = append(keysToListen, key)
	}
	app.CommitMultiStore().AddListeners(keysToListen)

	err = idxer.Start(nil)
	if err != nil {
		return err
	}

	streamingManager := storetypes.StreamingManager{
		ABCIListeners: []storetypes.ABCIListener{idxer},
		StopNodeOnErr: true,
	}
	app.SetStreamingManager(streamingManager)

	return nil
}

var startCmdOptions = server.StartCmdOptions{
	DBOpener: nil,
	PreSetup: preSetupIndexer,
	AddFlags: addIndexFlag,
}
