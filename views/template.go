package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/slmkb/weblensgo/context"
	"github.com/slmkb/weblensgo/models"
)

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
		},
	)
	if err := tpl.Execute(w, data); err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	htmlTpl := template.New(patterns[0])
	log.Printf("Patterns: %+v", patterns)
	htmlTpl = htmlTpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfFiled not implemented")
			},
			"currentUser": func() (template.HTML, error) {
				return "", fmt.Errorf("loggedIn not implemented")
			},
		},
	)
	htmlTpl, err := htmlTpl.ParseFS(fs, patterns...)
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
