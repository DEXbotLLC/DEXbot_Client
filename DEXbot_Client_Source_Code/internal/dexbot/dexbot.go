package dexbot

import (
	"dexbot/internal/dexbotABI"
	"dexbot/internal/dexbotUtils"
	"dexbot/internal/eventListener"
	"dexbot/internal/handler"
	"dexbot/internal/userConfig"
	"dexbot/internal/userWallets"
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

//Start listening for transactions to sign and send back to DEXbot
func StartDEXbot() {
	//Initialize DEXbot packages
	fmt.Println("Intializing DEXbot Client...")
	initalizeDEXbotClient()

	//Notify that DEXbot is listening for transactions to sign
	dexbotUtils.WhitePrinter.Println("Listening for transactions...")

	//Arrow animation to show that the program is alive while waiting for transactions to sign
	s := spinner.New(spinner.CharSets[31], 100*time.Millisecond)
	s.Start()

	//Wait for the wait group to finish to keep the program running forever
	handler.WaitForWaitGroupToFinish()
}

//Initialize all DEXbot packages
func initalizeDEXbotClient() {

	//Initialize application binary interfaces for the Router and ERC20 contracts
	dexbotABI.Initialize()

	//Add a task to the wait group to keep the program running forever
	handler.AddGoroutineToWaitGroup()

	//Initialize userWallets package
	userWallets.Initialize()

	//Initialize alive signal
	userConfig.InitializePulse()

	//Initialize eventListener package
	eventListener.Initialize()
}
