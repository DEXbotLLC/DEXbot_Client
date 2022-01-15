package eventListener

import (
	"dexbot/internal/authentication"
	"dexbot/internal/database"
	"dexbot/internal/handler"
	"fmt"

	"github.com/r3labs/sse/v2"
)

//Initalize listener to receive unsigned swap transactions from the database
var UnsignedSwapTransactionListener *Listener

//Initialize database listeners
func initalizeListeners() {

	//Initialize UnsignedSwapTransactionListener
	UnsignedSwapTransactionListener = &Listener{
		Url:                fmt.Sprintf("%s/unsigned_swap_tx/%s.json", database.DatabaseURL, authentication.FirebaseAuthToken.LocalId),
		Path:               fmt.Sprintf("unsigned_swap_tx/%s", authentication.FirebaseAuthToken.LocalId),
		EventChannel:       make(chan *sse.Event),
		Finished:           make(chan bool),
		HandleEvents:       handleUnsignedSwapTransactionEvents,
		AuthRefreshChannel: make(chan bool),
	}

	//Listen for event changes in unsigned_swap_tx/{firebase_user_id}
	go UnsignedSwapTransactionListener.Listen()
}

//Initialize function to wait for token to refresh, signaling the unsigned swap transaction listener to refresh
func initializeListenerRefresh() {
	for {
		<-authentication.AuthRefreshChannel
		UnsignedSwapTransactionListener.AuthRefreshChannel <- true
	}
}

//Start listening to changes in the database reference
func (l Listener) Listen() {

	//Sreate a new sse client with authenticated URL
	l.Client = sse.NewClient(l.Url + "?auth=" + authentication.FirebaseAuthToken.IdToken)

	//Subscribe to events changes
	err := l.Client.SubscribeChanRaw(l.EventChannel)
	handler.HandleError("listener: Listen: SubscribeChanRaw", err)

	//Handle event changes in the database reference as a goroutine
	go l.HandleEvents(l.EventChannel, l.Finished)

	//Wait for listener to refresh
	<-l.AuthRefreshChannel
	go l.RefreshAuthConnection()
}

//Stop listening to changes in the database reference
func (l Listener) StopListening() {

	//Unsubscribe to event changes in the database reference
	l.Client.Unsubscribe(l.EventChannel)

	//Send a boolean through the finished channel to stop the event handler
	l.Finished <- true
}

//Refresh authenticated connections
func (l Listener) RefreshAuthConnection() {

	//Stop listening with the previous authentication token
	l.StopListening()

	//Listen with the updated authenticated token
	l.Listen()
}
