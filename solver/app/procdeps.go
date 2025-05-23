package app

import (
	"context"
	"math/big"

	"github.com/omni-network/omni/contracts/bindings"
	"github.com/omni-network/omni/lib/contracts/solvernet"
	"github.com/omni-network/omni/lib/errors"
	"github.com/omni-network/omni/lib/ethclient/ethbackend"
	"github.com/omni-network/omni/lib/log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// procDeps abstracts dependencies for the event processor allowed simplified testing.
type procDeps struct {
	ParseID      func(chainID uint64, log types.Log) (OrderID, error)
	GetOrder     func(ctx context.Context, chainID uint64, id OrderID) (Order, bool, error)
	SetCursor    func(ctx context.Context, chainID uint64, height uint64) error
	ShouldReject func(ctx context.Context, order Order) (rejectReason, bool, error)
	DidFill      func(ctx context.Context, order Order) (bool, error)

	Accept func(ctx context.Context, order Order) error
	Reject func(ctx context.Context, order Order, reason rejectReason) error
	Fill   func(ctx context.Context, order Order) error
	Claim  func(ctx context.Context, order Order) error

	// Monitoring helpers
	TargetName func(Order) string
	ChainName  func(chainID uint64) string
}

func newClaimer(
	inboxContracts map[uint64]*bindings.SolverNetInbox,
	backends ethbackend.Backends,
	solverAddr common.Address,
) func(ctx context.Context, order Order) error {
	return func(ctx context.Context, order Order) error {
		inbox, ok := inboxContracts[order.SourceChainID]
		if !ok {
			return errors.New("unknown chain")
		}

		backend, err := backends.Backend(order.SourceChainID)
		if err != nil {
			return err
		}

		txOpts, err := backend.BindOpts(ctx, solverAddr)
		if err != nil {
			return err
		}

		// Claim to solver address for now
		// TODO: consider claiming to hot / cold funding wallet
		tx, err := inbox.Claim(txOpts, order.ID, solverAddr)
		if err != nil {
			return errors.Wrap(err, "claim order")
		} else if _, err := backend.WaitMined(ctx, tx); err != nil {
			return errors.Wrap(err, "wait mined")
		}

		return nil
	}
}

func newFiller(
	outboxContracts map[uint64]*bindings.SolverNetOutbox,
	backends ethbackend.Backends,
	solverAddr, outboxAddr common.Address,
) func(ctx context.Context, order Order) error {
	return func(ctx context.Context, order Order) error {
		if order.DestinationSettler != outboxAddr {
			return errors.New("destination settler mismatch [BUG] ", "got", order.DestinationSettler.Hex(), "expected", outboxAddr.Hex())
		}

		destChainID := order.DestinationChainID
		outbox, ok := outboxContracts[destChainID]
		if !ok {
			return errors.New("unknown chain")
		}

		backend, err := backends.Backend(destChainID)
		if err != nil {
			return err
		}

		callOpts := &bind.CallOpts{Context: ctx}
		txOpts, err := backend.BindOpts(ctx, solverAddr)
		if err != nil {
			return err
		}

		if ok, err := outbox.DidFill(callOpts, order.ID, order.FillOriginData); err != nil {
			return errors.Wrap(err, "did fill")
		} else if ok {
			log.Info(ctx, "Skipping already filled order", "order_id", order.ID)
			return nil
		}

		nativeValue := big.NewInt(0)
		for _, output := range order.MaxSpent {
			if output.ChainId.Uint64() != destChainID {
				// We error on this case for now, as our contracts only allow single dest chain orders
				// ERC7683 allows for orders with multiple destination chains, so continue-ing here
				// would also be appropriate.
				return errors.New("destination chain mismatch [BUG] ")
			}

			// zero token address means native token
			if output.Token == [32]byte{} {
				nativeValue.Add(nativeValue, output.Amount)
				continue
			}

			tknAddr := toEthAddr(output.Token)
			tkn, ok := tokens.Find(destChainID, tknAddr)
			if !ok {
				return errors.New("unsupported token, should have been rejected [BUG]", "addr", tknAddr.Hex(), "chain_id", destChainID)
			}

			isAppproved, err := isAppproved(ctx, tknAddr, backend, solverAddr, outboxAddr, output.Amount)
			if err != nil {
				return errors.Wrap(err, "is approved")
			}

			if !isAppproved {
				return errors.New("outbox not approved to spend token",
					"token", tkn.Symbol,
					"chain_id", destChainID,
					"addr", tknAddr.Hex(),
					"amount", output.Amount,
				)
			}
		}

		// xcall fee
		fee, err := outbox.FillFee(callOpts, order.FillOriginData)
		if err != nil {
			return errors.Wrap(err, "get fulfill fee")
		}

		txOpts.Value = new(big.Int).Add(nativeValue, fee)
		fillerData := []byte{} // fillerData is optional ERC7683 custom filler specific data, unused in our contracts
		tx, err := outbox.Fill(txOpts, order.ID, order.FillOriginData, fillerData)
		if err != nil {
			return errors.Wrap(err, "fill order", "custom", solvernet.DetectCustomError(err))
		} else if _, err := backend.WaitMined(ctx, tx); err != nil {
			return errors.Wrap(err, "wait mined")
		}

		if ok, err := outbox.DidFill(callOpts, order.ID, order.FillOriginData); err != nil {
			return errors.Wrap(err, "did fill")
		} else if !ok {
			return errors.New("fill failed [BUG]")
		}

		return nil
	}
}

func newRejector(
	inboxContracts map[uint64]*bindings.SolverNetInbox,
	backends ethbackend.Backends,
	solverAddr common.Address,
) func(ctx context.Context, order Order, reason rejectReason) error {
	return func(ctx context.Context, order Order, reason rejectReason) error {
		inbox, ok := inboxContracts[order.SourceChainID]
		if !ok {
			return errors.New("unknown chain")
		}

		backend, err := backends.Backend(order.SourceChainID)
		if err != nil {
			return err
		}

		// Ensure latest on-chain order is still pending
		if latest, err := inbox.GetOrder(&bind.CallOpts{Context: ctx}, order.ID); err != nil {
			return errors.Wrap(err, "get order")
		} else if latest.State.Status != statusPending.Uint8() {
			return errors.New("order status not pending anymore")
		}

		txOpts, err := backend.BindOpts(ctx, solverAddr)
		if err != nil {
			return err
		}

		tx, err := inbox.Reject(txOpts, order.ID, uint8(reason))
		if err != nil {
			return errors.Wrap(err, "reject order", "custom", solvernet.DetectCustomError(err))
		} else if _, err := backend.WaitMined(ctx, tx); err != nil {
			return errors.Wrap(err, "wait mined")
		}

		return nil
	}
}

func newDidFiller(outboxContracts map[uint64]*bindings.SolverNetOutbox) func(ctx context.Context, order Order) (bool, error) {
	return func(ctx context.Context, order Order) (bool, error) {
		outbox, ok := outboxContracts[order.DestinationChainID]
		if !ok {
			return false, errors.New("unknown chain")
		}

		filled, err := outbox.DidFill(&bind.CallOpts{Context: ctx}, order.ID, order.FillOriginData)
		if err != nil {
			return false, errors.Wrap(err, "did fill")
		}

		return filled, nil
	}
}

func newIDParser(inboxContracts map[uint64]*bindings.SolverNetInbox) func(chainID uint64, log types.Log) (OrderID, error) {
	return func(chainID uint64, log types.Log) (OrderID, error) {
		inbox, ok := inboxContracts[chainID]
		if !ok {
			return OrderID{}, errors.New("unknown chain")
		}

		event, ok := solvernet.EventByTopic(log.Topics[0])
		if !ok {
			return OrderID{}, errors.New("unknown event")
		}

		return event.ParseID(inbox.SolverNetInboxFilterer, log)
	}
}

func newOrderGetter(inboxContracts map[uint64]*bindings.SolverNetInbox) func(ctx context.Context, chainID uint64, id OrderID) (Order, bool, error) {
	return func(ctx context.Context, chainID uint64, id OrderID) (Order, bool, error) {
		inbox, ok := inboxContracts[chainID]
		if !ok {
			return Order{}, false, errors.New("unknown chain")
		}

		o, err := inbox.GetOrder(&bind.CallOpts{Context: ctx}, id)
		if err != nil {
			return Order{}, false, errors.Wrap(err, "get order")
		}

		// not found
		if o.Resolved.OrderId == [32]byte{} {
			return Order{}, false, nil
		}

		if o.Resolved.OrderId != id {
			return Order{}, false, errors.New("[BUG] order ID mismatch")
		}

		order, err := newOrder(o.Resolved, o.State)
		if err != nil {
			return Order{}, false, errors.Wrap(err, "new order")
		}

		return order, true, nil
	}
}
