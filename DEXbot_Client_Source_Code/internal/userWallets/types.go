package userWallets

import "dexbot/internal/web3Fork"

//Struct to hold wallet address and private key to sign incoming unsigned transactions
type UserWallet struct {
	WalletAddress string
	PrivateKey    *web3Fork.Key
}
