package main

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	r := chi.NewRouter()

	r.With(middleware.Logger).Get("/", homeHandler)

	http.ListenAndServe(":3000", r)
}
