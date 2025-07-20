package main

import (
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/slmkb/weblensgo/controllers"
	"github.com/slmkb/weblensgo/templates"
	"github.com/slmkb/weblensgo/views"
)

func main() {
	r := chi.NewRouter()

	homeTpl := views.Must(views.Parse(filepath.Join("templates", "home.gohtml")))
	r.With(middleware.Logger).Get("/", controllers.StaticHandler(homeTpl))

	contactTpl := views.Must(views.Parse("templates/base.gohtml", "templates/contact.gohtml"))
	r.Get("/contact", controllers.StaticHandler(contactTpl))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "base.gohtml", "faq.gohtml"))))

	r.Get("/admin", controllers.StaticHandler(
		views.Must(views.Parse(filepath.Join("templates", "admin.gohtml")))))

	http.ListenAndServe(":3000", r)
}
