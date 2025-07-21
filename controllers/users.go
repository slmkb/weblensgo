package controllers

import (
	"fmt"
	"log"
	"net/http"
)

type Users struct {
	Template struct {
		New Templater
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Template.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Printf("error parsing form: %v", err)
		fmt.Fprintf(w, "form error")
		return
	}
	fmt.Fprintf(w, "Parsed form %+v\n", r.PostForm)
	fmt.Fprintf(w, "Testing endpoint")
}
