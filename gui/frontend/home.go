package frontend

import (
	"twfinder/gui/server"
)

// HomeWin :
func HomeWin() server.Window {

	// Create and build a window
	win := server.NewWindow("Home", "Twitter Finder")
	win.Style().SetFullWidth()
	win.SetHAlign(server.HACenter)
	win.SetCellPadding(2)
	lbl := server.NewLabel("Welcome Home")
	win.Add(lbl)
	return win
}
