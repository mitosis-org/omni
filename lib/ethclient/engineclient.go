package ethclient

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/log"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/jinzhu/copier"
)

const (
	defaultRPCHTTPTimeout = time.Second * 30

	newPayloadV4 = "engine_newPayloadV4"

	forkchoiceUpdatedV3 = "engine_forkchoiceUpdatedV3"

	getPayloadV4 = "engine_getPayloadV4"
)

// EngineClient defines the Engine API authenticated JSON-RPC endpoints.
// It extends the normal Client interface with the Engine API.
type EngineClient interface {
	Client

	// NewPayloadV4 creates an Eth1 block, inserts it in the chain, and returns the status of the chain.
	// https://github.com/ethereum/execution-apis/blob/main/src/engine/prague.md#engine_newpayloadv4
	NewPayloadV4(ctx context.Context, params engine.ExecutableData, versionedHashes []common.Hash,
		beaconRoot *common.Hash, executionRequests []hexutil.Bytes) (engine.PayloadStatusV1, error)

	// ForkchoiceUpdatedV3 is equivalent to V2 with the addition of parent beacon block root in the payload attributes.
	// https://github.com/ethereum/execution-apis/blob/main/src/engine/cancun.md#engine_forkchoiceupdatedv3
	ForkchoiceUpdatedV3(ctx context.Context, update engine.ForkchoiceStateV1,
		payloadAttributes *engine.PayloadAttributes) (engine.ForkChoiceResponse, error)

	// GetPayloadV4 returns a cached payload by id.
	// https://github.com/ethereum/execution-apis/blob/main/src/engine/prague.md#engine_getpayloadv4
	GetPayloadV4(ctx context.Context, payloadID engine.PayloadID) (*engine.ExecutionPayloadEnvelope, error)
}

// engineClient implements EngineClient using JSON-RPC.
type engineClient struct {
	Wrapper
}

// NewAuthClient returns a new authenticated JSON-RPc engineClient.
func NewAuthClient(ctx context.Context, urlAddr string, jwtSecret []byte) (EngineClient, error) {
	transport := http.DefaultTransport
	if len(jwtSecret) > 0 {
		transport = newJWTRoundTripper(http.DefaultTransport, jwtSecret)
	}

	client := &http.Client{Timeout: defaultRPCHTTPTimeout, Transport: transport}

	rpcClient, err := rpc.DialOptions(ctx, urlAddr, rpc.WithHTTPClient(client))
	if err != nil {
		return engineClient{}, errors.Wrap(err, "rpc dial")
	}

	return engineClient{
		Wrapper: NewClient(rpcClient, "engine", urlAddr),
	}, nil
}

func (c engineClient) NewPayloadV4(ctx context.Context, params engine.ExecutableData, versionedHashes []common.Hash,
	beaconRoot *common.Hash, executionRequests []hexutil.Bytes,
) (engine.PayloadStatusV1, error) {
	const endpoint = "new_payload_v4"
	defer latency(c.chain, endpoint)()

	// isStatusOk returns true if the response status is valid.
	isStatusOk := func(status engine.PayloadStatusV1) bool {
		return map[string]bool{
			engine.VALID:    true,
			engine.INVALID:  true,
			engine.SYNCING:  true,
			engine.ACCEPTED: true,
		}[status.Status]
	}

	// NOTE: We should use this struct for compatibility with reth.
	// Otherwise, "Invalid params" error will be returned from reth.
	executionPayload := ExecutionPayloadV3{}
	if err := copier.Copy(&executionPayload, &params); err != nil {
		return engine.PayloadStatusV1{}, err
	}

	var resp engine.PayloadStatusV1
	var rpcErr rpc.Error
	err := c.cl.Client().CallContext(ctx, &resp, newPayloadV4, executionPayload, versionedHashes, beaconRoot, executionRequests)
	if isStatusOk(resp) {
		// Swallow errors when geth returns errors along with proper responses (but at least log it).
		if err != nil {
			log.Warn(ctx, "Ignoring new_payload_v4 error with proper response", err, "status", resp.Status)
		}

		return resp, nil
	} else if errors.As(err, &rpcErr) {
		// Swallow geth RPC errors, treat them as application errors, ie, invalid payload.
		// Geth server mostly returns status invalid with RPC errors, but the geth client doesn't
		// return errors AND status, it only returns errors OR status.
		log.Warn(ctx, "Converting new_payload_v4 engine rpc.Error to invalid response", err, "code", rpcErr.ErrorCode())
		valErr := err.Error()
		if data := errData(err); data != "" {
			valErr = data
		}

		return engine.PayloadStatusV1{
			Status:          engine.INVALID,
			ValidationError: &valErr,
		}, nil
	} else if err != nil {
		incError(c.chain, endpoint)
		return engine.PayloadStatusV1{}, errors.Wrap(err, "rpc new payload")
	} /* else err==nil && status!=ok */

	incError(c.chain, endpoint)

	return engine.PayloadStatusV1{}, errors.New("nil error and unknown status", "status", resp.Status)
}

func (c engineClient) ForkchoiceUpdatedV3(ctx context.Context, update engine.ForkchoiceStateV1,
	payloadAttributes *engine.PayloadAttributes,
) (engine.ForkChoiceResponse, error) {
	const endpoint = "forkchoice_updated_v3"
	defer latency(c.chain, endpoint)()

	// isStatusOk returns true if the response status is valid.
	isStatusOk := func(resp engine.ForkChoiceResponse) bool {
		return map[string]bool{
			engine.VALID:    true,
			engine.INVALID:  true,
			engine.SYNCING:  true,
			engine.ACCEPTED: false, // Unexpected in ForkchoiceUpdated
		}[resp.PayloadStatus.Status]
	}

	var resp engine.ForkChoiceResponse
	err := c.cl.Client().CallContext(ctx, &resp, forkchoiceUpdatedV3, update, payloadAttributes)
	if isStatusOk(resp) {
		// Swallow errors when geth returns errors along with proper responses (but at least log it).
		if err != nil {
			log.Warn(ctx, "Ignoring forkchoice_updated_v3 error with proper response", err, "status", resp.PayloadStatus.Status)
		}

		return resp, nil
	} else if err != nil {
		incError(c.chain, endpoint)
		return engine.ForkChoiceResponse{}, errors.Wrap(err, "rpc forkchoice updated v3")
	} /* else err==nil && status!=ok */

	incError(c.chain, endpoint)

	return engine.ForkChoiceResponse{}, errors.New("nil error and unknown status", "status", resp.PayloadStatus.Status)
}

func (c engineClient) GetPayloadV4(ctx context.Context, payloadID engine.PayloadID) (
	*engine.ExecutionPayloadEnvelope, error,
) {
	const endpoint = "get_payload_v4"
	defer latency(c.chain, endpoint)()

	var resp engine.ExecutionPayloadEnvelope
	err := c.cl.Client().CallContext(ctx, &resp, getPayloadV4, payloadID)
	if err != nil {
		incError(c.chain, endpoint)
		return nil, errors.Wrap(err, "rpc get payload v4")
	}

	return &resp, nil
}

// errData returns the error data if the error is a rpc.DataError.
func errData(err error) string {
	var dataErr rpc.DataError
	if errors.As(err, &dataErr) {
		return fmt.Sprint(dataErr.ErrorData())
	}

	return ""
}
