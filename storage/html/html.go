package html

import (
	"fmt"
	"os"
	"text/template"
	"time"
	"twfinder/static"
	"twfinder/storage"

	"github.com/tarekbadrshalaan/anaconda"
)

type html struct {
	tmpl       *template.Template
	storageDir string
	usersChan  chan anaconda.User
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

// Store :
// func (h *html) Store(filepath string, tmp *template.Template, data interface{}) error {
func (h *html) Store() {
	err := os.Mkdir(h.storageDir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	for {
		fName := fmt.Sprintf("%v/%v.html", h.storageDir, time.Now().Unix())
		f, err := os.Create(fName)
		if err != nil {
			fmt.Println(err)
			f.Close()
			continue
		}
		res := []anaconda.User{}
		for i := 0; i < static.RESULTPATCHSIZE; i++ {
			u := <-h.usersChan
			res = append(res, u)
		}
		h.tmpl.Execute(f, res)
		f.Close()
	}
}
