package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/slmkb/weblensgo/context"
	"github.com/slmkb/weblensgo/models"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Template struct {
		SignUp         Templater
		SignIn         Templater
		PasswordReset  Templater
		UpdatePassword Templater
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	EmailService         *models.EmailService
	PasswordResetService *models.PasswordResetService
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie(r, CookieSession)
		if err != nil {
			log.Printf("set user: %+v", err)
			next.ServeHTTP(w, r)
			return
		}
		if len(token) != 44 {
			log.Printf("token length violation")
			next.ServeHTTP(w, r)
			return
		}
		user, err := umw.SessionService.GetUser(token)
		if err != nil {
			log.Printf("set user: %+v", err)
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithUser(r.Context(), user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email     string
		CSRFField template.HTML
	}
	data.CSRFField = csrf.TemplateField(r)
	data.Email = r.FormValue("email")
	u.Template.SignUp.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		http.Error(w, "User creation failed", http.StatusInternalServerError)
		log.Printf("user create: %+v", err)
		return
	}

	session, err := u.SessionService.Create(user)
	if err != nil {
		log.Printf("session creation: %+v", err)
		http.Error(w, "Session creation failed", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) Signin(w http.ResponseWriter, r *http.Request) {
	var data struct {
		AuthError string
		Email     string
		CSRFField template.HTML
	}
	data.CSRFField = csrf.TemplateField(r)
	data.Email = r.FormValue("email")
	data.AuthError = r.FormValue("error")
	u.Template.SignIn.Execute(w, r, data)
}

func (u Users) ExecuteSignIn(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.GetUser(email, password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || errors.Is(err, sql.ErrNoRows) {
			url := fmt.Sprintf("/signin?error=1&email=%s", email)
			http.Redirect(w, r, url, http.StatusSeeOther)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user)
	if err != nil {
		http.Error(w, "Session creation failed", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Current user: %+v\n", context.User(r.Context()))
}

func (u Users) SignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		log.Printf("sign out: %+v", err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = u.SessionService.DeleteSession(token)
	if err != nil {
		log.Printf("sign out: %+v", err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	deleteCookie(w, CookieSession)
	log.Printf("logged out")
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (u Users) PasswordReset(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email     string
		CSRFField template.HTML
	}
	data.CSRFField = csrf.TemplateField(r)
	data.Email = r.FormValue("email")
	u.Template.PasswordReset.Execute(w, r, data)
}

func (u Users) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	pr, err := u.PasswordResetService.Create(email)
	if err != nil {
		log.Printf("send password reset email: %+v", err)
		return
	}
	err = u.EmailService.ForgotPassword(pr)
	if err != nil {
		log.Printf("send password reset email: %+v", err)
	}
	fmt.Fprintf(w, "password reset: %#v", pr)
}

func (u Users) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token     string
		CSRFField template.HTML
	}
	data.CSRFField = csrf.TemplateField(r)
	data.Token = r.FormValue("token")
	u.Template.UpdatePassword.Execute(w, r, data)
}

func (u Users) ExecutePasswordReset(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	password := r.FormValue("password")

	userID, err := u.PasswordResetService.Consume(token)
	if err != nil {
		log.Printf("execute password reset: %v", err)
		return
	}

	err = u.UserService.UpdatePassword(userID, password)
	if err != nil {
		log.Printf("execute password reset: %v", err)
		return
	}

	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}
