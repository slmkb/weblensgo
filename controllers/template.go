package controllers

import "net/http"

type Templater interface {
	Execute(w http.ResponseWriter, data any)
}
