package keeper

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/iovns/x/domain/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// QueryDomainRequest is the request made to
type QueryDomainRequest struct {
	Name string
}

type QueryDomainResponse struct {
	Domain types.Domain `json:"domain"`
}

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryDomain:
			return queryGet(ctx, path[1:], req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query request", types.ModuleName)
		}
	}
}

func queryGet(ctx sdk.Context, args []string, _ abci.RequestQuery, k Keeper) ([]byte, error) {
	resp, ok := k.GetDomain(ctx, args[0])
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, args[0])
	}
	return codec.MustMarshalJSONIndent(k.cdc, QueryDomainResponse{
		Domain: resp,
	}), nil
}

type QueryAccountsInDomain struct {
	Domain         string `json:"domain"`
	ResultsPerPage int    `json:"results_per_page"`
	Offset         int    `json:"offset"`
}

// Validate will validate the query model and set defaults
func (q *QueryAccountsInDomain) Validate() error {
	if q.Domain == "" {
		return sdkerrors.Wrapf(types.ErrInvalidDomainName, "empty")
	}
	// if results per page is unset then use default
	if q.ResultsPerPage <= 0 {
		q.ResultsPerPage = 100
	}
	// if offset is zero then use default
	if q.Offset <= 0 {
		q.Offset = 1
	}
	return nil
}

type QueryAccountsInDomainResponse struct {
	Accounts []types.Account `json:"accounts"`
}

// queryGetAccountsInDomain returns all accounts in a domain
func queryGetAccountsInDomain(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var query *QueryAccountsInDomain
	err := json.Unmarshal(req.Data, &query)
	if err != nil {
		return nil, err
	}
	// verify request
	if err = query.Validate(); err != nil {
		return nil, err
	}
	accKeys := k.GetAccountsInDomain(ctx, query.Domain)
	nKeys := len(accKeys) // total number of keys
	// no results
	if nKeys == 0 {
		return codec.MustMarshalJSONIndent(k.cdc, QueryAccountsInDomainResponse{}), nil
	}
	// get the index of the first object we want
	firstObjectIndex := query.Offset*query.ResultsPerPage - query.ResultsPerPage
	// check if there is at least one object at that index
	if nKeys < firstObjectIndex+1 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid offset")
	}

	// get the index for the last object
	lastObjectIndex := firstObjectIndex + query.ResultsPerPage - 1
	// check if last object index would outbound our acc keys slice
	if lastObjectIndex > nKeys {
		lastObjectIndex = nKeys - 1 // if it does then set last index as the last element of our slice
	}
	accounts := make([]types.Account, 0, lastObjectIndex-firstObjectIndex+1)
	// fill accounts
	for currIndex := firstObjectIndex; currIndex <= lastObjectIndex; currIndex++ {
		// get account
		account, _ := k.GetAccount(ctx, query.Domain, accountKeyToString(accKeys[currIndex]))
		// append
		accounts = append(accounts, account)
	}
	// return response
	return codec.MustMarshalJSONIndent(k.cdc, QueryAccountsInDomainResponse{Accounts: accounts}), nil
}
