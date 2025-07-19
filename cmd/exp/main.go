package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	u := struct {
		Name   string
		Age    int
		Nested struct {
			NestedMap   map[string]int
			NestedSlice []string
		}
	}{
		Name: "Kabekaes",
		Age:  9001,
		Nested: struct {
			NestedMap   map[string]int
			NestedSlice []string
		}{
			NestedMap: map[string]int{
				"Value1": 33,
				"Value2": 44,
			},
			NestedSlice: []string{
				"SliceString1",
				"SliceString2",
			},
		},
	}
	if err := tpl.Execute(w, u); err != nil {
		panic(err)
	}

}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	u := struct {
		Name   string
		Age    int
		Nested struct {
			NestedMap   map[string]int
			NestedSlice []string
		}
	}{
		Name: "Kabekaes",
		Age:  9001,
		Nested: struct {
			NestedMap   map[string]int
			NestedSlice []string
		}{
			NestedMap: map[string]int{
				"Value1": 33,
				"Value2": 44,
			},
			NestedSlice: []string{
				"SliceString1",
				"SliceString2",
			},
		},
	}

	if err = t.Execute(os.Stdout, u); err != nil {
		panic(err)
	}
	r := chi.NewRouter()

	r.With(middleware.Logger).Get("/param/{asdf}", parameterHandler)
	r.With(middleware.Logger).Get("/", homeHandler)

	r.Get("/dummy", dummyHanlder)

	if err := CreateUser(); err != nil {
		log.Println(err)
	}

	if err := CreateOrg(); err != nil {
		log.Println(err)
	}

	// http.ListenAndServe(":4000", r)

}
