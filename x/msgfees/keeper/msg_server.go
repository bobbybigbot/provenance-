package keeper

import (
	"context"

	"github.com/provenance-io/provenance/x/msgfees/types"
)

type msgServer struct {
	Keeper
}

func (m msgServer) CreateAdditionalFeeForMsgType(ctx context.Context, request *types.MsgAddFeeForMsgTypeRequest) (*types.CreateAdditionalFeeForMsgTypeResponse, error) {
	panic("implement me")
}

func (m msgServer) CalculateMsgBasedFees(ctx context.Context, request *types.CalculateFeePerMsgRequest) (*types.CalculateMsgBasedFeesResponse, error) {
	panic("implement me")
}

// NewMsgServerImpl returns an implementation of the msgfees MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}