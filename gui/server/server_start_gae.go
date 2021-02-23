// +build appengine

// Implementation of the GUI server Start on Google App Engine.

package server

import (
	"log"
	"net/http"
)

func (s *serverImpl) Start(openWins ...string) error {
	http.HandleFunc(s.appPath, func(w http.ResponseWriter, r *http.Request) {
		s.serveHTTP(w, r)
	})

	http.HandleFunc(s.appPath+pathStatic, func(w http.ResponseWriter, r *http.Request) {
		s.serveStatic(w, r)
	})

	log.Println("GAE - Starting GUI server on path:", s.appPath)
	if s.logger != nil {
		s.logger.Println("GAE - Starting GUI server on path:", s.appPath)
	}

	go s.sessCleaner()

	return nil
}
