package html

import (
	"fmt"
	"os"
	"text/template"
	"twfinder/logger"
	"twfinder/static"
	"twfinder/storage"

	"github.com/tarekbadrshalaan/anaconda"
)

type html struct {
	tmpl      *template.Template
	startFile int
}

// StorageObj :
type StorageObj struct {
	Users        []anaconda.User
	PreviousPage int
	NextPage     int
}

// BuildHTMLStore :
func BuildHTMLStore() (storage.IStorage, error) {
	tmpl, err := template.New("model").Parse(timelineTmpl)
	if err != nil {
		return nil, err
	}
	h := &html{tmpl: tmpl}

	// create storage directory
	htmldir := fmt.Sprintf("%v/html", static.STORAGEDIR)
	err = os.Mkdir(htmldir, os.ModePerm)
	if err != nil {
		logger.Error(err)
	}

	for i := 1; ; i++ {
		fname := fmt.Sprintf("%v/html/%v.html", static.STORAGEDIR, i)
		if _, err := os.Stat(fname); err != nil {
			h.startFile = i
			break
		}
	}
	return h, nil
}

// Store :
func (h *html) Store(usersChan <-chan anaconda.User) {
	counter := h.startFile
	for {
		str := StorageObj{
			PreviousPage: counter - 1,
			NextPage:     counter + 1,
			Users:        []anaconda.User{},
		}
		fName := fmt.Sprintf("%v/html/%v.html", static.STORAGEDIR, counter)
		f, err := os.Create(fName)
		if err != nil {
			logger.Error(err)
			f.Close()
			continue
		}
		for i := 0; i < static.RESULTPATCHSIZE; i++ {
			u := <-usersChan
			str.Users = append(str.Users, u)
		}
		h.tmpl.Execute(f, str)
		f.Close()
		counter++
	}
}
