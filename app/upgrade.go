package app

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// UpgradeHandler h for software upgrade proposal
type UpgradeHandler struct {
	*MinitiaApp
}

// NewUpgradeHandler return new instance of UpgradeHandler
func NewUpgradeHandler(app *MinitiaApp) UpgradeHandler {
	return UpgradeHandler{app}
}

func (h UpgradeHandler) CreateUpgradeHandler() upgradetypes.UpgradeHandler {
	return func(ctx context.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		return h.ModuleManager.RunMigrations(ctx, h.configurator, vm)
	}
}
