package eventListener

import "github.com/r3labs/sse/v2"

//Struct to hold event listener variables
type Listener struct {
	Url          string
	Path         string
	EventChannel chan *sse.Event
	Finished     chan bool
	//function to handle event changes in the database reference
	HandleEvents       func(chan *sse.Event, chan bool)
	Client             *sse.Client
	AuthRefreshChannel chan bool
}
