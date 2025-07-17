package main

import (
	"fmt"
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

func main() {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Joł"))
	})
	// r.Route("/param", func(r chi.Router) {
	// 	r.Use(middleware.Logger)
	// 	r.Get("/param/{asdf}", parameterHandler)
	// })

	r.With(middleware.Logger).Get("/param/{asdf}", parameterHandler)

	r.Get("/dummy", dummyHanlder)
	http.ListenAndServe(":3000", r)
}
