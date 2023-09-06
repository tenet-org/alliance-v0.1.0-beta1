package alliance

import (
	"fmt"

	"github.com/terra-money/alliance/x/alliance/keeper"
	"github.com/terra-money/alliance/x/alliance/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	// note: disable EndBlocker due to panic in RebalanceHook
	return []abci.ValidatorUpdate{}

	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyEndBlocker)
	k.CompleteRedelegations(ctx)
	if err := k.CompleteUndelegations(ctx); err != nil {
		panic(fmt.Errorf("failed to complete undelegations from x/alliance module: %s", err))
	}

	assets := k.GetAllAssets(ctx)
	k.InitializeAllianceAssets(ctx, assets)
	if _, err := k.DeductAssetsHook(ctx, assets); err != nil {
		panic(fmt.Errorf("failed to deduct take rate from alliance in x/alliance module: %s", err))
	}
	if err := k.RewardWeightChangeHook(ctx, assets); err != nil {
		panic(fmt.Errorf("failed to update assets reward weights in x/alliance module: %s", err))
	}
	if err := k.RebalanceHook(ctx, assets); err != nil {
		panic(fmt.Errorf("failed to rebalance assets in x/alliance module: %s", err))
	}
	return []abci.ValidatorUpdate{}
}
