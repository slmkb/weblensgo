package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/slmkb/weblensgo/apperrors"
	"github.com/slmkb/weblensgo/context"
	"github.com/slmkb/weblensgo/models"
)

type Galleries struct {
	Templates struct {
		New    Templater
		Edit   Templater
		Index  Templater
		Delete Templater
	}
	GalleryService *models.GalleryService
}

func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := g.GalleryService.GetGalleries(user.ID)
	if err != nil {
		g.Templates.Index.Execute(w, r, nil, err)
		log.Printf("%+v", err)
		return
	}
	g.Templates.Index.Execute(w, r, galleries)
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	g.Templates.New.Execute(w, r, nil)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	user := context.User(r.Context())
	_, err := g.GalleryService.Create(user.ID, title)
	if err != nil {
		if errors.Is(err, models.ErrGalleryExists) {
			err = apperrors.Public(err, "Gallery already exists.")
		} else {
			log.Printf("galleries: %+v", err)
		}
		g.Templates.New.Execute(w, r, nil, err)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	user := context.User(r.Context())
	hash := chi.URLParam(r, "hash")
	err := g.GalleryService.Rename(user.ID, hash, title)
	if err != nil {
		if errors.Is(err, models.ErrGalleryExists) {
			err = apperrors.Public(err, "Gallery already exists.")
		} else {
			log.Printf("galleries: %+v", err)
		}
		data := struct {
			Title string
			Hash  string
		}{
			Title: title,
			Hash:  hash,
		}
		g.Templates.Edit.Execute(w, r, data, err)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	user := context.User(r.Context())

	var data struct {
		Title string
		Hash  string
	}

	title, err := g.GalleryService.GetGallery(user.ID, hash)
	if err != nil {
		log.Printf("galleries: %+v", err)

	}
	data.Title = title
	data.Hash = hash
	g.Templates.Edit.Execute(w, r, data)
}

func (g Galleries) ConfirmDelete(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	confirm := r.FormValue("title.confirm")
	hash := chi.URLParam(r, "hash")
	if title != confirm {
		err := apperrors.Public(errors.New("strings not equal"), "Gallery name validation failed.")
		data := struct {
			Title string
			Hash  string
		}{
			Title: confirm,
			Hash:  hash,
		}
		g.Templates.Delete.Execute(w, r, data, err)
		return
	}

	user := context.User(r.Context())

	err := g.GalleryService.Delete(user.ID, hash)
	if err != nil {
		log.Printf("galleries: %+v", err)
		g.Templates.Index.Execute(w, r, nil, err)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	user := context.User(r.Context())

	var data struct {
		Title string
		Hash  string
	}

	title, err := g.GalleryService.GetGallery(user.ID, hash)
	if err != nil {
		log.Printf("galleries: %+v", err)

	}
	data.Title = title
	data.Hash = hash
	g.Templates.Delete.Execute(w, r, data)
}
