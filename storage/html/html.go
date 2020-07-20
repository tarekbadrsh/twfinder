package html

import (
	"fmt"
	"os"
	"text/template"
	"twfinder/static"
	"twfinder/storage"

	"github.com/tarekbadrshalaan/anaconda"
)

const storageDir = "result"

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
	err = os.Mkdir(storageDir, os.ModePerm)
	if err != nil {
		// todo add logger ...
		fmt.Println(err)
	}
	for i := 1; ; i++ {
		fname := fmt.Sprintf("%v/%v.html", storageDir, i)
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
		fName := fmt.Sprintf("%v/%v.html", storageDir, counter)
		f, err := os.Create(fName)
		if err != nil {
			fmt.Println(err)
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
