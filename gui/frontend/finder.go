package frontend

import (
	"twfinder/gui/server"
)

// FinderWin :
func FinderWin() server.Window {
	// Create and build a window
	win := server.NewWindow("finder", "Finder - Twitter Finder App")
	win.Style().SetFullWidth()
	win.SetHAlign(server.HACenter)
	win.SetCellPadding(2)
	lbl := server.NewLabel("Welcome finder")
	win.Add(lbl)
	return win
}
