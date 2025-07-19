package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func executeTemplate(w http.ResponseWriter, templateFile string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tplPath := filepath.Join("templates", templateFile)
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		log.Printf("parsing template error: %v", err)
		http.Error(w, "parsing error", http.StatusInternalServerError)
		return
	}
	if err := tpl.Execute(w, nil); err != nil {
		log.Printf("execute template error: %v", err)
		http.Error(w, "execute template error", http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "home.gohtml")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "contact.gohtml")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "faq.html")
}

func main() {
	r := chi.NewRouter()

	r.With(middleware.Logger).Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)

	http.ListenAndServe(":3000", r)
}
