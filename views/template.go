package views

import (
	"errors"
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

type public interface {
	Public() string
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data any, errs ...error) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}
	errMsgs := errorMessages(errs...)
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				return errMsgs
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
			"errors": func() []string {
				return nil
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

func errorMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			msgs = append(msgs, "Something went wrong.")
		}
	}
	return msgs
}
