package main

import (
	"flag"
	"twfinder/config"
	"twfinder/gui/frontend"
	"twfinder/gui/server"
	"twfinder/logger"
)

func main() {
	// read Command-Line Flags
	configPath := flag.String("c", "config.json", "configuration file path")
	flag.Parse()

	/* configuration initialize start */
	config.BuildConfiguration(*configPath)
	/* configuration initialize end */

	/* logger initialize start */
	mylogger := logger.NewZapLogger("ERROR", "DEBUG")
	logger.InitializeLogger(&mylogger)
	defer logger.Close()
	/* logger initialize end */

	// Create and start a GUI server (omitting error check)
	server := server.NewServer("", "localhost:8081")
	server.SetText("Twitter Finder App")
	server.AddWin(frontend.HomeWin())
	server.AddWin(frontend.ConfigWin())
	server.AddWin(frontend.FinderWin())
	server.SetDefaultRootWindow(frontend.HomeWin())
	server.Start("home") // Also opens windows list in browser
}
