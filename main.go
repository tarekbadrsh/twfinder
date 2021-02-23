package main

import (
	"os"
	"twfinder/config"
	"twfinder/finder"
	"twfinder/logger"
	"twfinder/pipeline"
	"twfinder/request"
)

func main() {
	/* configuration initialize start */
	c := config.Configuration()
	/* configuration initialize end */

	/* logger initialize start */
	mylogger := logger.NewZapLogger()
	logger.InitializeLogger(&mylogger)
	defer logger.Close()
	/* logger initialize end */

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
