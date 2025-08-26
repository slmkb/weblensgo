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

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	if err := godotenv.Load(); err != nil {
		return config{}, err
	}
	cfg := config{
		PSQL: models.PostgresConfig{
			Host:     os.Getenv("PSQL_HOST"),
			Port:     os.Getenv("PSQL_PORT"),
			User:     os.Getenv("PSQL_USERNAME"),
			Password: os.Getenv("PSQL_PASSWORD"),
			Database: os.Getenv("PSQL_DATABASE"),
			SSLMode:  os.Getenv("PSQL_SSLMODE"),
		},
		SMTP: models.SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			User:     os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
		},
		CSRF: struct {
			Key    string
			Secure bool
		}{
			Key:    os.Getenv("CSRF_KEY"),
			Secure: os.Getenv("CSRF_SECURE") == "true",
		},
		Server: struct{ Address string }{
			Address: os.Getenv("HTTP_ADDRESS"),
		},
	}
	return cfg, nil
}

func main() {

	cfg, err := loadEnvConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	if err := run(cfg); err != nil {
		log.Fatalf("error run %v", err)
	}
}

func run(cfg config) error {
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return fmt.Errorf("database open: %v", err)
	}
	defer db.Close()

	smtpServer, err := models.NewSMTPClient(cfg.SMTP)
	if err != nil {
		return fmt.Errorf("smtpserver: %v", err)
	}

	usersCtrl := controllers.Users{
		UserService: &models.UserService{
			DB: db,
		},
		SessionService: &models.SessionService{
			DB: db,
		},
		PasswordResetService: &models.PasswordResetService{
			DB: db,
		},
		EmailService: &models.EmailService{
			SMTPServer: smtpServer,
		},
	}

	galleryCtrl := controllers.Galleries{
		GalleryService: &models.GalleryService{
			DB: db,
		},
	}
	userMw := controllers.UserMiddleware{
		SessionService: usersCtrl.SessionService,
	}

	err = models.MigrateFS(db, sqlFS.FS, "")
	if err != nil {
		return fmt.Errorf("migrate: %v", err)
	}

	homeTpl := views.Must(views.ParseFS(templates.FS, "base.gohtml", "home.gohtml"))
	contactTpl := views.Must(views.ParseFS(templates.FS, "base.gohtml", "contact.gohtml"))
	faqTpl := views.Must(views.ParseFS(templates.FS, "base.gohtml", "faq.gohtml"))
	usersCtrl.Template.Register = views.Must(views.ParseFS(templates.FS, "base.gohtml", "signup.gohtml"))
	usersCtrl.Template.SignIn = views.Must(views.ParseFS(templates.FS, "base.gohtml", "signin.gohtml"))
	usersCtrl.Template.ForgotPassword = views.Must(views.ParseFS(templates.FS, "base.gohtml", "forgot-password.gohtml"))
	usersCtrl.Template.ResetPassword = views.Must(views.ParseFS(templates.FS, "base.gohtml", "reset-password.gohtml"))
	usersCtrl.Template.SendResetLink = views.Must(views.ParseFS(templates.FS, "base.gohtml", "send-reset-link.gohtml"))
	galleryCtrl.Templates.Index = views.Must(views.ParseFS(templates.FS, "base.gohtml", "galleries/index.gohtml"))
	galleryCtrl.Templates.New = views.Must(views.ParseFS(templates.FS, "base.gohtml", "galleries/new.gohtml"))
	galleryCtrl.Templates.Edit = views.Must(views.ParseFS(templates.FS, "base.gohtml", "galleries/edit.gohtml"))
	galleryCtrl.Templates.Delete = views.Must(views.ParseFS(templates.FS, "base.gohtml", "galleries/delete.gohtml"))
	galleryCtrl.Templates.Show = views.Must(views.ParseFS(templates.FS, "base.gohtml", "galleries/show.gohtml"))
	// usersCtrl.Template.Galleries = views.Must(views.ParseFS(templates.FS, "base.gohtml", "galleries/galleries.gohtml"))

	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
		csrf.TrustedOrigins([]string{"localhost:3000", "localhost"}),
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
	r.Get("/signup", usersCtrl.RegistrationForm)
	r.Post("/signup", usersCtrl.ProcessRegistration)
	r.Get("/signin", usersCtrl.SignInForm)
	r.With( /*middleware.Logger*/ ).Post("/signin", usersCtrl.ProcessSignIn)
	r.Get("/signout", usersCtrl.ProcessSignOut)
	r.Get("/forgot-password", usersCtrl.ForgotPassword)
	r.Post("/forgot-password", usersCtrl.SendPasswordResetEmail)
	r.Get("/reset-password", usersCtrl.ResetPassword)
	r.Post("/reset-password", usersCtrl.ProcessResetPassword)
	r.Route("/galleries", func(r chi.Router) {

		r.Group(func(r chi.Router) {
			r.Use(userMw.RequireUser)
			// r.Get("/", usersCtrl.CurrentUser)
			r.Get("/", galleryCtrl.Index)
			r.Get("/new", galleryCtrl.New)
			r.Post("/new", galleryCtrl.Create)
			r.Get("/{hash}/edit", galleryCtrl.Edit)
			r.Post("/{hash}", galleryCtrl.Update)
			r.Get("/{hash}/delete", galleryCtrl.Delete)
			r.Post("/{hash}/delete", galleryCtrl.ConfirmDelete)
			r.Get("/{hash}", galleryCtrl.Show)
			r.Get("/{hash}/{filename}", galleryCtrl.Image)
			r.Post("/{hash}/{filename}/delete", galleryCtrl.DeleteImage)
			r.Post("/{hash}/images", galleryCtrl.UploadImage)
		})
	})
	// r.With(userMw.RequireUser).Get("/users/me", usersCtrl.CurrentUser)

	fmt.Println("Starting server...")
	http.ListenAndServe(cfg.Server.Address, r)
	fmt.Println("Exiting...")
	return nil
}
