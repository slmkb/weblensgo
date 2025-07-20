package views

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.htmlTpl.Execute(w, data); err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
}

func Parse(filepath ...string) (Template, error) {
	htmlTpl, err := template.ParseFiles(filepath...)
	if err != nil {
		log.Printf("parsing template: %v", err)
		return Template{}, err
	}
	return Template{
		htmlTpl: htmlTpl,
	}, nil
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	htmlTpl, err := template.ParseFS(fs, pattern...)
	if err != nil {
		log.Printf("parsing template: %v", err)
		return Template{}, err
	}
	return Template{
		htmlTpl: htmlTpl,
	}, nil
}
