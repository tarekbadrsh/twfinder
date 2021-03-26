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
	pagecount int
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
			h.pagecount = i
			break
		}
	}
	return h, nil
}

// Store :
func (h *html) Store(users []anaconda.User) {
	str := StorageObj{
		PreviousPage: h.pagecount - 1,
		NextPage:     h.pagecount + 1,
		Users:        users,
	}
	fName := fmt.Sprintf("%v/html/%v.html", static.STORAGEDIR, h.pagecount)
	f, err := os.Create(fName)
	if err != nil {
		logger.Error(err)
		f.Close()
		return
	}
	defer f.Close()
	err = h.tmpl.Execute(f, str)
	if err != nil {
		logger.Error(err)
	}
	h.pagecount = h.pagecount + 1
}
