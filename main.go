package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/slmkb/weblensgo/controllers"
	"github.com/slmkb/weblensgo/models"
	"github.com/slmkb/weblensgo/templates"
	"github.com/slmkb/weblensgo/views"
)

func main() {
	r := chi.NewRouter()

	homeTpl := views.Must(views.Parse("templates/base.gohtml", "templates/home.gohtml"))
	r.With(middleware.Logger).Get("/", controllers.StaticHandler(homeTpl))

	contactTpl := views.Must(views.ParseFS(templates.FS, "base.gohtml", "contact.gohtml"))
	r.Get("/contact", controllers.StaticHandler(contactTpl))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "base.gohtml", "faq.gohtml"))))

	r.Get("/admin", controllers.StaticHandler(
		views.Must(views.Parse(filepath.Join("templates", "admin.gohtml")))))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		log.Fatalf("database open: %v", err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}

	usersCtrl := controllers.Users{
		UserService: &userService,
	}

	usersCtrl.Template.New = views.Must(views.Parse("templates/base.gohtml", "templates/signup.gohtml"))
	r.With(middleware.Logger).Get("/signup", usersCtrl.New)

	r.With(middleware.Logger).Post("/signup", usersCtrl.Create)

	// r.With(middleware.Logger).Get("/signin", usersCtrl.Signin1)

	fmt.Println("Starting server...")
	http.ListenAndServe(":3000", r)
	fmt.Println("Exiting...")
}
