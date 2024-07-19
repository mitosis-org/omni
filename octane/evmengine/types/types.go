package types

import (
	"github.com/cometbft/cometbft/crypto"
	"github.com/ethereum/go-ethereum/common"
)

// TODO(thai): temporary customization before ethos will support secp256k1.
type AddressProvider interface {
	//// LocalAddress returns the local validator's ethereum address.
	//LocalAddress() common.Address

	// PubKey returns the local validator's public key. (ecdsa)
	PubKey() crypto.PubKey
}

type FeeRecipientProvider interface {
	// LocalFeeRecipient returns the local validator's fee recipient address.
	LocalFeeRecipient() common.Address
	// VerifyFeeRecipient returns true if the given address is a valid fee recipient
	VerifyFeeRecipient(proposedFeeRecipient common.Address) error
}
