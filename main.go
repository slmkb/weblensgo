package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/slmkb/weblensgo/controllers"
	"github.com/slmkb/weblensgo/models"
	"github.com/slmkb/weblensgo/templates"
	"github.com/slmkb/weblensgo/views"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	homeTpl := views.Must(views.ParseFS(templates.FS, "base.gohtml", "home.gohtml"))
	r.Get("/", controllers.StaticHandler(homeTpl))

	contactTpl := views.Must(views.ParseFS(templates.FS, "base.gohtml", "contact.gohtml"))
	r.Get("/contact", controllers.StaticHandler(contactTpl))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "base.gohtml", "faq.gohtml"))))

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

	usersCtrl.Template.SignUp = views.Must(views.ParseFS(templates.FS, "base.gohtml", "signup.gohtml"))
	usersCtrl.Template.SignIn = views.Must(views.ParseFS(templates.FS, "base.gohtml", "signin.gohtml"))
	r.Get("/signup", usersCtrl.New)

	r.Post("/signup", usersCtrl.Create)

	r.Get("/signin", usersCtrl.Signin)
	r.With(middleware.Logger).Post("/signin", usersCtrl.ExecuteSignIn)

	csrfKey := "abcdefghijklmnopqrstuvwxyz"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(false),
		csrf.TrustedOrigins([]string{"localhost:3000"}),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[CSRF BLOCKED] Host=%s, Method=%s, URL=%s, Origin=%s, Referer=%s, Reason=%v",
				r.Host, r.Method, r.URL.String(), r.Header.Get("Origin"), r.Header.Get("Referer"), csrf.FailureReason(r))

			// Respond with default 403 Forbidden
			http.Error(w, "Forbidden - CSRF check failed", http.StatusForbidden)
		})),
	)

	fmt.Println("Starting server...")

	// r.Use(csrfMw)
	http.ListenAndServe(":3000", csrfMw(r))
	fmt.Println("Exiting...")
}
