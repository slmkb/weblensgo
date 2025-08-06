package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/slmkb/weblensgo/controllers"
	"github.com/slmkb/weblensgo/models"
	sqlFS "github.com/slmkb/weblensgo/models/sql"
	"github.com/slmkb/weblensgo/templates"
	"github.com/slmkb/weblensgo/views"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file")
	}

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		log.Fatalf("database open: %v", err)
	}
	defer db.Close()

	smtpServer, err := models.NewSMTPClient(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)
	if err != nil {
		log.Fatalf("smtpserver: %v", err)
	}

	usersCtrl := controllers.Users{
		UserService: &models.UserService{
			DB: db,
		},
		SessionService: &models.SessionService{
			DB: db,
		},
		EmailService: &models.EmailService{
			SMTPServer: smtpServer,
		},
		PasswordResetService: &models.PasswordResetService{
			DB: db,
		},
	}

	// if err := usersCtrl.EmailService.ForgotPassword(models.User{}); err != nil {
	// 	log.Fatalf("email send: %+v", err)
	// }

	userMw := controllers.UserMiddleware{
		SessionService: usersCtrl.SessionService,
	}

	err = models.MigrateFS(db, sqlFS.FS, "")
	if err != nil {
		log.Fatalf("migrate: %v", err)
	}

	homeTpl := views.Must(views.ParseFS(templates.FS, "base.gohtml", "home.gohtml"))
	contactTpl := views.Must(views.ParseFS(templates.FS, "base.gohtml", "contact.gohtml"))
	usersCtrl.Template.SignUp = views.Must(views.ParseFS(templates.FS, "base.gohtml", "signup.gohtml"))
	usersCtrl.Template.SignIn = views.Must(views.ParseFS(templates.FS, "base.gohtml", "signin.gohtml"))
	usersCtrl.Template.PasswordReset = views.Must(views.ParseFS(templates.FS, "base.gohtml", "password-reset.gohtml"))
	usersCtrl.Template.UpdatePassword = views.Must(views.ParseFS(templates.FS, "base.gohtml", "update-password.gohtml"))
	faqTpl := views.Must(views.ParseFS(templates.FS, "base.gohtml", "faq.gohtml"))

	csrfMw := csrf.Protect(
		[]byte(os.Getenv("CSRF_KEY")),
		csrf.Secure(false),
		csrf.TrustedOrigins([]string{"localhost:3000"}),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[CSRF BLOCKED] Host=%s, Method=%s, URL=%s, Origin=%s, Referer=%s, Reason=%v",
				r.Host, r.Method, r.URL.String(), r.Header.Get("Origin"), r.Header.Get("Referer"), csrf.FailureReason(r))

			http.Error(w, "Forbidden - CSRF check failed", http.StatusForbidden)
		})),
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger, csrfMw, userMw.SetUser)

	r.Get("/", controllers.StaticHandler(homeTpl))
	r.Get("/contact", controllers.StaticHandler(contactTpl))
	r.Get("/faq", controllers.FAQ(faqTpl))
	r.Get("/signup", usersCtrl.New)
	r.Post("/signup", usersCtrl.Create)
	r.Get("/signin", usersCtrl.Signin)
	r.Get("/signout", usersCtrl.SignOut)
	r.Get("/password-reset", usersCtrl.PasswordReset)
	r.Post("/password-reset", usersCtrl.SendPasswordResetEmail)
	r.Get("/update-password", usersCtrl.UpdatePassword)
	r.Post("/update-password", usersCtrl.ExecutePasswordReset)
	r.With( /*middleware.Logger*/ ).Post("/signin", usersCtrl.ExecuteSignIn)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(userMw.RequireUser)
		r.Get("/", usersCtrl.CurrentUser)
	})
	// r.With(userMw.RequireUser).Get("/users/me", usersCtrl.CurrentUser)

	fmt.Println("Starting server...")
	http.ListenAndServe(":3000", r)
	fmt.Println("Exiting...")
}
