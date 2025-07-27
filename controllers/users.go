package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/slmkb/weblensgo/models"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Template struct {
		SignUp Templater
		SignIn Templater
	}
	UserService *models.UserService
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
		return
	}
	fmt.Fprintf(w, "user created successfully: %+v", user)
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
	// fmt.Fprintf(w, "HEADERD: %+v", r.Header)
	u.Template.SignIn.Execute(w, r, data)
	// fmt.Fprintf(w, "%s:%s", email, password)
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
		// fmt.Println(err)
		return
	}

	fmt.Fprintf(w, "user logged in successfully: %+v", user)
}
