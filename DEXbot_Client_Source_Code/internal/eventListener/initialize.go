package eventListener

//Initialize the eventListener package
func Initialize() {

	//Initialize database listeners
	initalizeListeners()
	//Initialize refersh for authenticated connection
	go initializeListenerRefresh()
}
