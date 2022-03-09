package web3Logic

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"dexbot/internal/handler"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/umbracle/fastrlp"
	"github.com/umbracle/go-web3"
	"github.com/valyala/fastjson"
	"golang.org/x/crypto/sha3"
)

//Sign a transaction with a specified wallet key
func SignTransaction(chainID uint64, txn *web3.Transaction, privateKey *Key) *web3.Transaction {
	//create a signer object with the bsc mainnet chainID
	signer := NewEIP155Signer(chainID)

	//sign the transaction
	signedTx, err := signer.SignTx(txn, privateKey)
	handler.HandleError("Error: SignTransaction, SignTx", err)

	//return the signed transaction
	return signedTx
}

//Key Functions--------------------------------------------

func (k *Key) Address() web3.Address {
	return k.addr
}

func (k *Key) MarshallPrivateKey() ([]byte, error) {
	return (*btcec.PrivateKey)(k.priv).Serialize(), nil
}

func (k *Key) SignMsg(msg []byte) ([]byte, error) {
	return k.Sign(keccak256(msg))
}

func (k *Key) Sign(hash []byte) ([]byte, error) {
	sig, err := btcec.SignCompact(S256, (*btcec.PrivateKey)(k.priv), hash, false)
	if err != nil {
		return nil, err
	}
	term := byte(0)
	if sig[0] == 28 {
		term = 1
	}
	return append(sig, term)[1:], nil
}

func NewKey(priv *ecdsa.PrivateKey) *Key {
	return &Key{
		priv: priv,
		pub:  &priv.PublicKey,
		addr: pubKeyToAddress(&priv.PublicKey),
	}
}

func pubKeyToAddress(pub *ecdsa.PublicKey) (addr web3.Address) {
	b := keccak256(elliptic.Marshal(S256, pub.X, pub.Y)[1:])
	copy(addr[:], b[12:])
	return
}

// GenerateKey generates a new key based on the secp256k1 elliptic curve.
func GenerateKey() (*Key, error) {
	priv, err := ecdsa.GenerateKey(S256, rand.Reader)
	if err != nil {
		return nil, err
	}
	return NewKey(priv), nil
}

func EcrecoverMsg(msg, signature []byte) (web3.Address, error) {
	return Ecrecover(keccak256(msg), signature)
}

func Ecrecover(hash, signature []byte) (web3.Address, error) {
	pub, err := RecoverPubkey(signature, hash)
	if err != nil {
		return web3.Address{}, err
	}
	return pubKeyToAddress(pub), nil
}

func RecoverPubkey(signature, hash []byte) (*ecdsa.PublicKey, error) {
	size := len(signature)
	term := byte(27)
	if signature[size-1] == 1 {
		term = 28
	}

	sig := append([]byte{term}, signature[:size-1]...)
	pub, _, err := btcec.RecoverCompact(S256, sig, hash)
	if err != nil {
		return nil, err
	}
	return pub.ToECDSA(), nil
}

func keccak256(buf []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(buf)
	b := h.Sum(nil)
	return b
}

//Signer Functions-----------------------------------------
func NewEIP155Signer(chainID uint64) *EIP1155Signer {
	return &EIP1155Signer{chainID: chainID}
}

func (e *EIP1155Signer) SignTx(tx *web3.Transaction, key *Key) (*web3.Transaction, error) {
	hash := signHash(tx, e.chainID)

	sig, err := key.Sign(hash)
	if err != nil {
		return nil, err
	}

	vv := uint64(sig[64]) + 35 + e.chainID*2

	tx.R = sig[:32]
	tx.S = sig[32:64]
	tx.V = new(big.Int).SetUint64(vv).Bytes()
	return tx, nil
}

func signHash(tx *web3.Transaction, chainID uint64) []byte {
	a := fastrlp.DefaultArenaPool.Get()

	v := a.NewArray()
	v.Set(a.NewUint(tx.Nonce))
	v.Set(a.NewUint(tx.GasPrice))
	v.Set(a.NewUint(tx.Gas))
	if tx.To == nil {
		v.Set(a.NewNull())
	} else {
		v.Set(a.NewCopyBytes((*tx.To)[:]))
	}
	v.Set(a.NewBigInt(tx.Value))
	v.Set(a.NewCopyBytes(tx.Input))

	// EIP155
	if chainID != 0 {
		v.Set(a.NewUint(chainID))
		v.Set(a.NewUint(0))
		v.Set(a.NewUint(0))
	}

	hash := keccak256(v.MarshalTo(nil))
	fastrlp.DefaultArenaPool.Put(a)
	return hash
}

//Marshal Transaction to JSON
var defaultArena fastjson.ArenaPool

func MarshalJSON(t *web3.Transaction) ([]byte, error) {
	a := defaultArena.Get()

	o := a.NewObject()
	o.Set("hash", a.NewString(t.Hash.String()))
	o.Set("from", a.NewString(t.From.String()))
	if len(t.Input) != 0 {
		o.Set("input", a.NewString("0x"+hex.EncodeToString(t.Input)))
	}
	if t.Value != nil {
		o.Set("value", a.NewString(fmt.Sprintf("0x%x", t.Value)))
	}
	o.Set("gasPrice", a.NewString(fmt.Sprintf("0x%x", t.GasPrice)))
	o.Set("gas", a.NewString(fmt.Sprintf("0x%x", t.Gas)))
	if t.Nonce != 0 {
		// we can remove this once we include support for custom nonces
		o.Set("nonce", a.NewString(fmt.Sprintf("0x%x", t.Nonce)))
	}
	if t.To == nil {
		o.Set("to", a.NewNull())
	} else {
		o.Set("to", a.NewString(t.To.String()))
	}
	o.Set("v", a.NewString("0x"+hex.EncodeToString(t.V)))
	o.Set("r", a.NewString("0x"+hex.EncodeToString(t.R)))
	o.Set("s", a.NewString("0x"+hex.EncodeToString(t.S)))

	o.Set("blockHash", a.NewString(t.BlockHash.String()))
	o.Set("blockNumber", a.NewString(fmt.Sprintf("0x%x", t.BlockNumber)))
	o.Set("transactionIndex", a.NewString(fmt.Sprintf("0x%x", t.TxnIndex)))

	res := o.MarshalTo(nil)
	defaultArena.Put(a)
	return res, nil
}
