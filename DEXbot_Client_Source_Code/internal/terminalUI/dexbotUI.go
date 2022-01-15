package terminalUI

import "dexbot/internal/dexbot"

//Clear the UI display and start DEXbot
func startDEXbot() {
	clearTerminal()
	dexbot.StartDEXbot()
}
