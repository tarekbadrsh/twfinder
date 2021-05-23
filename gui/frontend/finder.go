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
	lbl := server.NewLabel("Collecting Data")
	win.Add(lbl)
	lodImg := server.NewHTML(`<iframe src="https://giphy.com/embed/VX7yEoXAFf8as" width="480" height="480" frameBorder="0" class="giphy-embed" allowFullScreen></iframe><p><a href="https://giphy.com/gifs/today-loading-icon-VX7yEoXAFf8as">via GIPHY</a></p>`)
	win.Add(lodImg)
	bckhomBtn := server.NewButton("back to home")
	bckhomBtn.AddEHandlerFunc(func(e server.Event) {
		e.ReloadWin("home")
	}, server.ETypeClick)
	win.Add(bckhomBtn)
	return win
}
