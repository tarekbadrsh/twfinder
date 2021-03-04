package main

import (
	"os"
	"twfinder/config"
	"twfinder/finder"
	"twfinder/gui/frontend"
	"twfinder/gui/server"
	"twfinder/logger"
	"twfinder/pipeline"
	"twfinder/request"
)

func main() {
	/* logger initialize start */
	mylogger := logger.NewZapLogger()
	logger.InitializeLogger(&mylogger)
	defer logger.Close()
	/* logger initialize end */

	win := frontend.ConfigWin()
	// Create and start a GUI server (omitting error check)
	server := server.NewServer("guitest", "localhost:8081")
	server.SetText("Test GUI App")
	server.AddWin(win)
	server.Start("Configuration") // Also opens windows list in browser

	/* configuration initialize start */
	c := config.Configuration()
	/* configuration initialize end */

	/* finder build start */
	finder.BuildSearchCriteria(c)
	/* finder build end */

	/* build TwitterAPI start */
	request.TwitterAPI()
	/* build TwitterAPI end */

	/* start Pipline */
	pip := pipeline.NewPipeline()
	pip.Start()
	/* start Pipline */

	// shutdown the application gracefully
	cancelChan := make(chan os.Signal, 1)
	sig := <-cancelChan
	logger.Infof("Caught SIGTERM %v\n", sig)
	pip.Close()
}
