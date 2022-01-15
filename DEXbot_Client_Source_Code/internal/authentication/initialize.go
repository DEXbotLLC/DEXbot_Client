package authentication

//Initialize the authentication package
func Initalize(username string, password string) {

	//Initialize the firebase username and password
	initalizeFirebaseCreds(username, password)

	//Initialize the authentication token to connect to the database
	initializeAuthTokenWithRefresh()
}
