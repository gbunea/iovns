package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/iov-one/iovns"
	types2 "github.com/iov-one/iovns/x/configuration/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/iov-one/iovns/x/fee/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	feeQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	feeQueryCmd.AddCommand(
		flags.GetCommands(
			getCmdQueryFees(queryRoute, cdc),
		)...,
	)

	return feeQueryCmd
}

func getCmdQueryFees(route string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "get-fees",
		Short: "gets the current fees",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			path := fmt.Sprintf("custom/%s/%s", route, types2.QueryFees)
			resp, _, err := cliCtx.Query(path)
			if err != nil {
				return err
			}
			var jsonResp types2.QueryFeesResponse
			err = iovns.DefaultQueryDecode(resp, &jsonResp)
			if err != nil {
				return err
			}
			return cliCtx.PrintOutput(jsonResp)
		},
	}
}
