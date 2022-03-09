package userConfig

//Initialize the userConfig package
func Initialize() {

	//Initialize the user configuration
	initializeUserConfig()

	//Send the version checksum to DEXbot
	SendVersionChecksum()

}
