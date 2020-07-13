package html

import (
	"fmt"
	"os"
	"text/template"
	"twfinder/static"
	"twfinder/storage"

	"github.com/tarekbadrshalaan/anaconda"
)

type html struct {
	tmpl       *template.Template
	storageDir string
	usersChan  chan anaconda.User
	startFile  int
}

// StorageObj :
type StorageObj struct {
	Users        []anaconda.User
	PreviousPage int
	NextPage     int
}

// BuildHTMLStore :
func BuildHTMLStore(storageDir string, usersChan chan anaconda.User) (storage.IStorage, error) {
	tmpl, err := template.New("model").Parse(timelineTmpl)
	if err != nil {
		return nil, err
	}
	h := &html{tmpl: tmpl, storageDir: storageDir, usersChan: usersChan}
	return h, nil
}

func (h *html) MkStoragedir() {
	err := os.Mkdir(h.storageDir, os.ModePerm)
	if err != nil {
		// todo add logger ...
		fmt.Println(err)
	}
	for i := 1; ; i++ {
		fname := fmt.Sprintf("result/%v.html", i)
		if _, err := os.Stat(fname); err != nil {
			h.startFile = i
			return
		}
	}
}

// Store :
func (h *html) Store() {
	h.MkStoragedir()
	counter := h.startFile
	for {
		str := StorageObj{
			PreviousPage: counter - 1,
			NextPage:     counter + 1,
			Users:        []anaconda.User{},
		}
		fName := fmt.Sprintf("%v/%v.html", h.storageDir, counter)
		f, err := os.Create(fName)
		if err != nil {
			fmt.Println(err)
			f.Close()
			continue
		}
		for i := 0; i < static.RESULTPATCHSIZE; i++ {
			u := <-h.usersChan
			str.Users = append(str.Users, u)
		}
		h.tmpl.Execute(f, str)
		f.Close()
		counter++
	}
}
