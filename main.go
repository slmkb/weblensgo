package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/slmkb/weblensgo/views"
)

func executeTemplate(w http.ResponseWriter, templateFile string) {

	// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// 	tplPath := filepath.Join("templates", templateFile)
	tpl, err := views.Parse(templateFile)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was a problem parsing the template", http.StatusInternalServerError)
		return
	}
	tpl.Execute(w, nil)
	// tpl, err := template.ParseFiles(tplPath)
	//
	//	if err != nil {
	//		log.Printf("parsing template error: %v", err)
	//		http.Error(w, "parsing error", http.StatusInternalServerError)
	//		return
	//	}
	//
	//	if err := tpl.Execute(w, nil); err != nil {
	//		log.Printf("execute template error: %v", err)
	//		http.Error(w, "execute template error", http.StatusInternalServerError)
	//	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tplPath)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tplPath)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "faq.gohtml")
	executeTemplate(w, tplPath)
}

func main() {
	r := chi.NewRouter()
	// tplPath := filepath.Join("templates", "home.gohtml")
	// faq := template.Must(views.Parse(tplPath))
	r.With(middleware.Logger).Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)

	http.ListenAndServe(":3000", r)
}
