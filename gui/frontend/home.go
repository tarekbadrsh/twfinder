package frontend

import (
	"twfinder/gui/server"
)

// HomeWin :
func HomeWin() server.Window {
	// Create and build a window
	win := server.NewWindow("home", "Home - Twitter Finder App")
	win.Style().SetFullWidth()
	win.SetHAlign(server.HACenter)
	win.SetCellPadding(2)
	lbl := server.NewLabel("Welcome Home")
	win.Add(lbl)
	configBtn := server.NewButton("Configuration")
	configBtn.AddEHandlerFunc(func(e server.Event) {
		e.ReloadWin("configuration")
	}, server.ETypeClick)
	win.Add(configBtn)
	finderBtn := server.NewButton("Finder")
	finderBtn.AddEHandlerFunc(func(e server.Event) {
		e.ReloadWin("finder")
	}, server.ETypeClick)
	win.Add(finderBtn)
	return win
}
