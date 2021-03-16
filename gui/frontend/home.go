package frontend

import (
	"twfinder/finder"
	"twfinder/gui/server"
	"twfinder/pipeline"
	"twfinder/request"
)

func buildPipeline() {
	/* finder build start */
	finder.BuildSearchCriteria()
	/* finder build end */

	/* build TwitterAPI start */
	request.TwitterAPI()
	/* build TwitterAPI end */
}

// HomeWin :
func HomeWin() server.Window {
	/* start Pipline */
	pip := pipeline.NewPipeline()
	/* start Pipline */

	// Create and build a window
	win := server.NewWindow("home", "Home - Twitter Finder App")
	win.Style().SetFullWidth()
	win.SetHAlign(server.HACenter)
	win.SetCellPadding(2)
	configBtn := server.NewButton("Configuration")
	configBtn.AddEHandlerFunc(func(e server.Event) {
		e.ReloadWin("configuration")
	}, server.ETypeClick)
	win.Add(configBtn)
	lodImg := server.NewHTML(`<iframe src="https://giphy.com/embed/VX7yEoXAFf8as" width="480" height="480" frameBorder="0" class="giphy-embed" allowFullScreen></iframe><p><a href="https://giphy.com/gifs/today-loading-icon-VX7yEoXAFf8as">via GIPHY</a></p> </div>`)
	lblTitle := server.NewLabel("")

	startBtn := server.NewButton("Start")
	startBtn.AddEHandlerFunc(func(e server.Event) {
		lblTitle.SetText("Collecting Data ... ")
		win.Add(lblTitle)
		win.Add(lodImg)
		//
		buildPipeline()
		pip.Start()
		//
		e.MarkDirty(win)
	}, server.ETypeClick)
	win.Add(startBtn)

	stopBtn := server.NewButton("stop")
	stopBtn.AddEHandlerFunc(func(e server.Event) {
		lblTitle.SetText("Stop collection data !")
		win.Remove(lodImg)
		//
		pip.Close()
		//
		e.MarkDirty(win)
	}, server.ETypeClick)
	win.Add(stopBtn)
	return win
}
