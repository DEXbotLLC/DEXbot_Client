package web3Fork

import (
	"crypto/ecdsa"

	"github.com/btcsuite/btcd/btcec"
	"github.com/umbracle/go-web3"
)

//Struct to store wallet key details
type Key struct {
	priv *ecdsa.PrivateKey
	pub  *ecdsa.PublicKey
	addr web3.Address
}

//Struct used to sign transactions
type EIP1155Signer struct {
	chainID uint64
}

// S256 is the secp256k1 elliptic curve
var S256 = btcec.S256()
