package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	cosmosauthtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/provenance-io/provenance/x/msgfees/types"
	"github.com/tendermint/tendermint/libs/log"
)

// StoreKey is the store key string for authz
const StoreKey = types.ModuleName

// Keeper of the Additional fee store
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              codec.BinaryCodec
	paramSpace       paramtypes.Subspace
	feeCollectorName string // name of the FeeCollector ModuleAccount
	defaultFeeDenom string
}

// NewKeeper returns a AdditionalFeeKeeper. It handles:
// CONTRACT: the parameter Subspace must have the param key table already initialized
func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	feeCollectorName string,
	defaultFeeDenom string,
) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		paramSpace:       paramSpace,
		feeCollectorName: feeCollectorName,
		defaultFeeDenom:  defaultFeeDenom,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) GetFeeCollectorName() string {
	return k.feeCollectorName
}

func (k Keeper) GetDefaultFeeDenom() string {
	return k.defaultFeeDenom
}

// SetMsgBasedFee sets the additional fee schedule for a Msg
func (k Keeper) SetMsgBasedFee(ctx sdk.Context, msgBasedFees types.MsgBasedFee) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&msgBasedFees)
	store.Set(types.GetMsgBasedFeeKey(msgBasedFees.MsgTypeUrl), bz)
}

// GetMsgBasedFee returns a MsgBasedFee for the msg type if it exists nil if it does not
func (k Keeper) GetMsgBasedFee(ctx sdk.Context, msgType string) (*types.MsgBasedFee, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetMsgBasedFeeKey(msgType)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil, nil
	}

	var msgBasedFee types.MsgBasedFee
	if err := k.cdc.Unmarshal(bz, &msgBasedFee); err != nil {
		return nil, err
	}

	return &msgBasedFee, nil
}

// RemoveMsgBasedFee removes MsgBasedFee or returns an error if it does not exist
func (k Keeper) RemoveMsgBasedFee(ctx sdk.Context, msgType string) error {
	store := ctx.KVStore(k.storeKey)
	key := types.GetMsgBasedFeeKey(msgType)
	bz := store.Get(key)
	if len(bz) == 0 {
		return types.ErrMsgFeeDoesNotExist
	}

	store.Delete(key)

	return nil
}

type Handler func(record types.MsgBasedFee) (stop bool)

// IterateMsgBasedFees  iterates all msg fees with the given handler function.
func (k Keeper) IterateMsgBasedFees(ctx sdk.Context, handle func(msgFees types.MsgBasedFee) (stop bool)) error {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.MsgBasedFeeKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		record := types.MsgBasedFee{}
		if err := k.cdc.Unmarshal(iterator.Value(), &record); err != nil {
			return err
		}
		if handle(record) {
			break
		}
	}
	return nil
}

// DeductFees deducts fees from the given account, the only reason it exists is that the
func (k Keeper) DeductFees(bankKeeper cosmosauthtypes.BankKeeper, ctx sdk.Context, acc cosmosauthtypes.AccountI, fees sdk.Coins) error {
	ctx.Logger().Info("NOTICE: In DeductFees:" + ctx.GasMeter().String())
	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	err := bankKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), k.feeCollectorName, fees)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}
	ctx.Logger().Info("NOTICE: End of DeductFees:" + ctx.GasMeter().String())
	return nil
}

// ExportGenesis returns a GenesisState for a given context.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// TODO Implement me
	return nil
}

// InitGenesis new msgfees genesis
func (k Keeper) InitGenesis(ctx sdk.Context, data *types.GenesisState) {
	// TODO Implement me
}
