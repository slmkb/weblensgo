package main

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/slmkb/weblensgo/controllers"
	"github.com/slmkb/weblensgo/views"
)

func main() {
	r := chi.NewRouter()

	homeTpl, err := views.Parse(filepath.Join("templates", "home.gohtml"))
	if err != nil {
		panic(err)
	}
	r.With(middleware.Logger).Get("/", controllers.StaticHandler(homeTpl))

	contactTpl, err := views.Parse(filepath.Join("templates", "contact.gohtml"))
	if err != nil {
		panic(err)
	}

	r.Get("/contact", controllers.StaticHandler(contactTpl))
	faqTpl, err := views.Parse(filepath.Join("templates", "faq.gohtml"))
	if err != nil {
		panic(err)
	}
	r.Get("/faq", controllers.StaticHandler(faqTpl))

	http.ListenAndServe(":3000", r)
}
