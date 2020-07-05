package storage

import (
	"os"
	"text/template"
)

// Template :
func Template(temp string) (*template.Template, error) {
	tmpl, err := template.New("model").Parse(temp)

	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// Store :
func Store(filepath string, tmp *template.Template, data interface{}) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmp.Execute(f, data)
}
