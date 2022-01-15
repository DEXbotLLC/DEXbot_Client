package authentication

import "time"

//Struct to store firebase username and password
type firebaseCredentials struct {
	firebaseEmail    string
	firebasePassword string
}

//Struct to employ the Token() method and act as a token source for firebase auth
type Authenticator struct {
}

//Struct to store the Authentication token
type authToken struct {
	Kind         string
	LocalId      string
	Email        string
	DisplayName  string
	IdToken      string
	Registered   bool
	RefreshToken string
	ExpiresIn    time.Time
}

//Struct to store the refresh token
type refreshToken struct {
	Id_token string
}
