package controllers

import (
	"fmt"
	"net/http"

	"github.com/slmkb/weblensgo/models"
)

type Users struct {
	Template struct {
		New Templater
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Template.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	// if err := r.ParseForm(); err != nil {
	// 	log.Printf("error parsing form: %v", err)
	// 	fmt.Fprintf(w, "form error")
	// 	return
	// }
	// fmt.Fprintf(w, "Parsed form %+v\n", r.PostForm)
	// fmt.Fprintf(w, "Testing endpoint")
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		http.Error(w, "User creation failed", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "user created successfully: %+v", user)
}

// func (u Users) Signin(w http.ResponseWriter, r *http.Request) {
// 	email := r.FormValue("email")
// 	password := r.FormValue("password")
// 	fmt.Fprintf(w, "%s:%s", email, password)
// }

// func (u Users) Signin1(w http.ResponseWriter, r *http.Request) {
// 	u.Template.New.Execute(w, nil)
// 	// fmt.Fprintf(w, "%s:%s", email, password)
// }
