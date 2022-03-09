package web3Logic

import (
	"dexbot/internal/handler"

	"github.com/umbracle/go-web3/jsonrpc"
)

func InitRPCClient(nodeURL string) *jsonrpc.Client {
	client, err := jsonrpc.NewClient(nodeURL)
	handler.HandleError("web3: initWeb3HTTPClient: NewClient", err)
	return client
}
