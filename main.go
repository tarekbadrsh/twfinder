package main

import (
	"flag"
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
	// read Command-Line Flags
	configPath := flag.String("c", "config.json", "configuration file path")
	flag.Parse()

	/* configuration initialize start */
	config.BuildConfiguration(*configPath)
	/* configuration initialize end */

	/* logger initialize start */
	mylogger := logger.NewZapLogger()
	logger.InitializeLogger(&mylogger)
	defer logger.Close()
	/* logger initialize end */

	// Create and start a GUI server (omitting error check)
	server := server.NewServer("", "localhost:8081")
	server.SetText("Twitter Finder App")
	server.AddWin(frontend.ConfigWin())
	server.AddWin(frontend.HomeWin())
	server.AddWin(frontend.FinderWin())
	server.SetDefaultRootWindow(frontend.HomeWin())
	server.Start("home") // Also opens windows list in browser

	/* finder build start */
	finder.BuildSearchCriteria()
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
