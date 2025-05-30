package keeper

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/omni-network/omni/halo/comet"
	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/ethclient"
	"github.com/omni-network/omni/lib/k1util"
	"github.com/omni-network/omni/octane/evmengine/types"

	"github.com/ethereum/go-ethereum/beacon/engine"
	etypes "github.com/ethereum/go-ethereum/core/types"

	ormv1alpha1 "cosmossdk.io/api/cosmos/orm/v1alpha1"
	"cosmossdk.io/core/store"
	"cosmossdk.io/orm/model/ormdb"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	grpc1 "github.com/cosmos/gogoproto/grpc"
)

type Keeper struct {
	cdc                    codec.BinaryCodec
	storeService           store.KVStoreService
	headTable              ExecutionHeadTable
	withdrawalTable        WithdrawalTable
	engineCl               ethclient.EngineClient
	txConfig               client.TxConfig
	voteProvider           types.VoteExtensionProvider
	eventProcs             []types.EvmEventProcessor
	cmtAPI                 comet.API
	addrProvider           types.AddressProvider
	feeRecProvider         types.FeeRecipientProvider
	buildDelay             time.Duration
	buildOptimistic        bool
	maxWithdrawalsPerBlock uint64

	// mutablePayload contains the previous optimistically triggered payload.
	// It is optimistic because the validator set can change,
	// so we might not actually be the next proposer.
	mutablePayload struct {
		sync.Mutex
		ID        engine.PayloadID
		Height    uint64
		UpdatedAt time.Time
	}
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	engineCl ethclient.EngineClient,
	txConfig client.TxConfig,
	addrProvider types.AddressProvider,
	feeRecProvider types.FeeRecipientProvider,
	maxWithdrawalsPerBlock uint64,
	eventProcs ...types.EvmEventProcessor,
) (*Keeper, error) {
	schema := &ormv1alpha1.ModuleSchemaDescriptor{SchemaFile: []*ormv1alpha1.ModuleSchemaDescriptor_FileEntry{
		{Id: 1, ProtoFileName: File_octane_evmengine_keeper_evmengine_proto.Path()},
	}}

	modDB, err := ormdb.NewModuleDB(schema, ormdb.ModuleDBOptions{KVStoreService: storeService})
	if err != nil {
		return nil, errors.Wrap(err, "create module db")
	}

	dbStore, err := NewEvmengineStore(modDB)
	if err != nil {
		return nil, errors.Wrap(err, "create evmengine store")
	}

	if err := verifyProcs(eventProcs); err != nil {
		return nil, err
	}

	return &Keeper{
		cdc:                    cdc,
		storeService:           storeService,
		headTable:              dbStore.ExecutionHeadTable(),
		withdrawalTable:        dbStore.WithdrawalTable(),
		engineCl:               engineCl,
		txConfig:               txConfig,
		addrProvider:           addrProvider,
		feeRecProvider:         feeRecProvider,
		eventProcs:             eventProcs,
		maxWithdrawalsPerBlock: maxWithdrawalsPerBlock,
	}, nil
}

// verifyProcs ensures that all event processors have distinct names.
func verifyProcs(eventProcs []types.EvmEventProcessor) error {
	names := make(map[string]bool)
	for _, proc := range eventProcs {
		name := proc.Name()
		if names[name] {
			return errors.New("duplicate event processors", "name", name)
		}
		names[name] = true
	}

	return nil
}

func (k *Keeper) SetVoteProvider(p types.VoteExtensionProvider) {
	k.voteProvider = p
}

// SetCometAPI sets the comet API client.
func (k *Keeper) SetCometAPI(c comet.API) {
	k.cmtAPI = c
}

// SetBuildDelay sets the build delay parameter.
func (k *Keeper) SetBuildDelay(d time.Duration) {
	k.buildDelay = d
}

// SetBuildOptimistic sets the optimistic build parameter.
func (k *Keeper) SetBuildOptimistic(b bool) {
	k.buildOptimistic = b
}

// RegisterProposalService registers the proposal service on the provided router.
// This implements abci.ProcessProposal verification of new proposals.
func (k *Keeper) RegisterProposalService(server grpc1.Server) {
	types.RegisterMsgServiceServer(server, NewProposalServer(k))
}

// parseAndVerifyProposedPayload parses and verifies and returns the proposed payload
// if comparing it against the latest execution block succeeds.
func (k *Keeper) parseAndVerifyProposedPayload(ctx context.Context, msg *types.MsgExecutionPayload) (engine.ExecutableData, error) {
	// Block of magellan network upgrade height is built by uluwatu logic,
	// support both that and new magellan proto payloads.
	// Note this should strictly only be allowed for that first block, but is tricky to enforce.
	if msg.ExecutionPayloadDeneb != nil && len(msg.ExecutionPayload) > 0 {
		return engine.ExecutableData{}, errors.New("only one payload type allowed")
	} else if msg.ExecutionPayloadDeneb == nil && len(msg.ExecutionPayload) == 0 {
		return engine.ExecutableData{}, errors.New("no payload provided")
	}

	// Parse the payload.
	var payload engine.ExecutableData
	if msg.ExecutionPayloadDeneb != nil {
		var err error
		payload, err = types.PayloadFromProto(msg.ExecutionPayloadDeneb)
		if err != nil {
			return engine.ExecutableData{}, errors.Wrap(err, "unmarshal proto payload")
		}
	} else if err := json.Unmarshal(msg.ExecutionPayload, &payload); err != nil {
		return engine.ExecutableData{}, errors.Wrap(err, "unmarshal payload")
	}

	eligibleWithdrawals, err := k.EligibleWithdrawals(ctx, false)
	if err != nil {
		return engine.ExecutableData{}, errors.Wrap(err, "eligible withdrawals")
	}

	if !withdrawalsEqual(payload.Withdrawals, eligibleWithdrawals) {
		return engine.ExecutableData{}, errors.New("mismatch with eligible withdrawals")
	}

	// Ensure no witness
	if payload.ExecutionWitness != nil {
		return engine.ExecutableData{}, errors.New("witness not allowed in payload")
	}

	// Ensure fee recipient using provider
	if err := k.feeRecProvider.VerifyFeeRecipient(payload.FeeRecipient); err != nil {
		return engine.ExecutableData{}, errors.Wrap(err, "verify proposed fee recipient")
	}

	// Fetch the latest execution head from the local keeper DB.
	head, err := k.GetExecutionHead(ctx)
	if err != nil {
		return engine.ExecutableData{}, errors.Wrap(err, "latest execution block")
	}
	headHash, err := head.Hash()
	if err != nil {
		return engine.ExecutableData{}, err
	}

	// Ensure the parent hash and block height matches
	if payload.Number != head.GetBlockHeight()+1 {
		return engine.ExecutableData{}, errors.New("invalid proposed payload number", "proposed", payload.Number, "head", head.GetBlockHeight())
	} else if payload.ParentHash != headHash {
		return engine.ExecutableData{}, errors.New("invalid proposed payload parent hash", "proposed", payload.ParentHash, "head", headHash)
	}

	// Ensure the payload timestamp is after latest execution block and before or equaled to the current consensus block.
	minTimestamp := head.GetBlockTime() + 1
	maxTimestamp := uint64(sdk.UnwrapSDKContext(ctx).BlockTime().Unix())
	if maxTimestamp < minTimestamp { // Execution block minimum takes precedence
		maxTimestamp = minTimestamp
	}
	if payload.Timestamp < minTimestamp || payload.Timestamp > maxTimestamp {
		return engine.ExecutableData{}, errors.New("invalid payload timestamp",
			"proposed", payload.Timestamp, "min", minTimestamp, "max", maxTimestamp,
		)
	}

	// Ensure the Randao Digest is equaled to parent hash as this is our workaround at this point.
	if payload.Random != headHash {
		return engine.ExecutableData{}, errors.New("invalid payload random", "proposed", payload.Random, "latest", headHash)
	}

	return payload, nil
}

// isNextProposer returns true if the local node is the proposer
// for the next block.
//
// Note that the validator set can change, so this is an optimistic check.
func (k *Keeper) isNextProposer(ctx context.Context, currentHeight int64) (bool, error) {
	// cometAPI is lazily set and may be nil on startup (e.g. rollbacks).
	if k.cmtAPI == nil {
		return false, nil
	}

	valset, err := k.cmtAPI.Validators(ctx, currentHeight)
	if err != nil {
		return false, err
	}

	nextProposer := valset.CopyIncrementProposerPriority(1).Proposer
	nextAddr, err := k1util.PubKeyToAddress(nextProposer.PubKey)
	if err != nil {
		return false, err
	}

	isNextProposer := nextAddr == k.addrProvider.LocalAddress()

	return isNextProposer, nil
}

func (k *Keeper) setOptimisticPayload(id engine.PayloadID, height uint64) {
	k.mutablePayload.Lock()
	defer k.mutablePayload.Unlock()

	k.mutablePayload.ID = id
	k.mutablePayload.Height = height
	k.mutablePayload.UpdatedAt = time.Now()
}

func (k *Keeper) getOptimisticPayload() (engine.PayloadID, uint64, time.Time) {
	k.mutablePayload.Lock()
	defer k.mutablePayload.Unlock()

	return k.mutablePayload.ID, k.mutablePayload.Height, k.mutablePayload.UpdatedAt
}

func withdrawalsEqual(w1, w2 []*etypes.Withdrawal) bool {
	if len(w1) != len(w2) {
		return false
	}

	for i, w := range w1 {
		if w.Index != w2[i].Index ||
			w.Validator != w2[i].Validator ||
			w.Address != w2[i].Address ||
			w.Amount != w2[i].Amount {
			return false
		}
	}

	return true
}
