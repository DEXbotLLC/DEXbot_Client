package terminalUI

import (
	"dexbot/internal/handler"

	"github.com/pterm/pterm"
)

//Variables to display UI
var terminalArea *pterm.AreaPrinter
var terminalCenterPrinter *pterm.CenterPrinter
var terminalPrinter pterm.CenterPrinter

//Initialize the terminal area to display the UI
func initializeTerminalArea() {
	area, err := pterm.DefaultArea.WithRemoveWhenDone().Start()
	handler.HandleError("Error when initializing terminal area", err)
	terminalArea = area
}

//Initalize the terminal printer to display the UI
func initializeTerminalCenterPrinter() {
	printer := pterm.DefaultCenter.WithCenterEachLineSeparately()
	terminalCenterPrinter = printer
}
