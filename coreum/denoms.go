package coreum

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// QueryDenoms returns list of available denoms. the paginationKey should nil for the first page.
// the nextPaginationKey will be nil if there are no more pages.
func (r *Reader) QueryDenoms(
	ctx context.Context, bankClient banktypes.QueryClient, paginationKey []byte,
) (data types.Coins, nextPaginationKey []byte, err error) {
	res, err := bankClient.TotalSupply(ctx, &banktypes.QueryTotalSupplyRequest{
		Pagination: &query.PageRequest{Key: paginationKey},
	})
	if err != nil {
		return nil, nil, err
	}
	return res.Supply, res.Pagination.NextKey, nil
}

// QueryDenomsMetadata returns list of available denoms metadata. the paginationKey should nil for the first page.
// the nextPaginationKey will be nil if there are no more pages.
func (r *Reader) QueryDenomsMetadata(
	ctx context.Context, bankClient banktypes.QueryClient, paginationKey []byte,
) (data []banktypes.Metadata, nextPaginationKey []byte, err error) {
	res, err := bankClient.DenomsMetadata(ctx, &banktypes.QueryDenomsMetadataRequest{
		Pagination: &query.PageRequest{Key: paginationKey},
	})
	if err != nil {
		return nil, nil, err
	}
	return res.Metadatas, res.Pagination.NextKey, nil
}
