package main

import (
	"fmt"
	"os"
	"twfinder/config"
	"twfinder/finder"
	"twfinder/pipeline"
	"twfinder/request"
)

func main() {
	/* configuration initialize start */
	c := config.Configuration()
	/* configuration initialize end */

	/* finder build start */
	finder.BuildSearchCriteria()
	/* finder build end */

	/* build TwitterAPI start */
	request.TwitterAPI()
	/* build TwitterAPI end */

	/* start Pipline */
	pip := pipeline.NewPipeline(c.SearchUser)
	pip.Start()
	/* start Pipline */

	// shutdown the application gracefully
	cancelChan := make(chan os.Signal, 1)
	sig := <-cancelChan
	fmt.Printf("Caught SIGTERM %v\n", sig)
	pip.Close()
}
