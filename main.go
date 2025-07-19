package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func parameterHandler(w http.ResponseWriter, r *http.Request) {
	fetchParam := chi.URLParam(r, "*")
	if fetchParam != "" {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Congratz!!!"))
		fmt.Fprintf(w, "%+v", r)

		return
	}
	fmt.Fprintf(w, "%+v", r)
	w.Write([]byte("parameter required"))
}

func dummyHanlder(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Dummy endpoint"))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tpl, err := template.ParseFiles("templates/home.gohtml")
	if err != nil {
		panic(err)
	}

	user := struct {
		Name      string
		TestSlice []string
		TestMap   map[string]int
	}{
		Name:      "Szu",
		TestSlice: []string{"one", "two", "three"},
		TestMap: map[string]int{
			"One":   1,
			"Two":   2,
			"Three": 3,
		},
	}

	if err := tpl.Execute(w, user); err != nil {
		panic(err)
	}
}

func main() {
	r := chi.NewRouter()

	r.Get("/", homeHandler)

	r.With(middleware.Logger).Get("/param/{asdf}", parameterHandler)

	r.Get("/dummy", dummyHanlder)
	http.ListenAndServe(":3000", r)
}
