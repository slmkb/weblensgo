package controllers

import "net/http"

type Templater interface {
	Execute(w http.ResponseWriter, r *http.Request, data any, errs ...error)
}
