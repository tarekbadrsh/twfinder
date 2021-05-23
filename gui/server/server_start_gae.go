// +build appengine

// Implementation of the GUI server Start on Google App Engine.

package server

import (
	"net/http"
	"twfinder/logger"
)

func (s *serverImpl) Start(openWins ...string) error {
	http.HandleFunc(s.appPath, func(w http.ResponseWriter, r *http.Request) {
		s.serveHTTP(w, r)
	})

	http.HandleFunc(s.appPath+pathStatic, func(w http.ResponseWriter, r *http.Request) {
		s.serveStatic(w, r)
	})

	logger.Infof("Starting GUI server on path:", s.appPath)

	go s.sessCleaner()

	return nil
}
