package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/iov-one/iovns/x/domain/types"
	fee "github.com/iov-one/iovns/x/fee/calculator"
)

// CollectFees collects the fees of a msg and sends them
// to the distribution module to validators and stakers
func (k Keeper) CollectFees(ctx sdk.Context, msg types.MsgWithFeePayer, domain types.Domain) error {
	moduleFees := k.FeeKeeper.GetFees(ctx)
	// create fee calculator
	calculator := fee.NewCalculator(ctx, k, moduleFees, domain)
	// get fee
	fee := calculator.GetFee(msg)
	// transfer fee to distribution
	return k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.FeePayer(), auth.FeeCollectorName, sdk.NewCoins(fee))
}
