package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

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
		Show   Templater
	}
	GalleryService *models.GalleryService
}

var (
	ErrEmptyString = errors.New("controllers: empty gallery title")
)

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
	if title == "" {
		g.Templates.New.Execute(w, r, nil, apperrors.Public(ErrEmptyString, "Gallery name cannot be empty."))
		return
	}
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
	galleryHash := chi.URLParam(r, "hash")
	err := g.GalleryService.Edit(user.ID, galleryHash, title)
	if err != nil {
		if errors.Is(err, models.ErrGalleryExists) {
			err = apperrors.Public(err, "Gallery already exists.")
		} else {
			log.Printf("galleries: %+v", err)
		}
		data := struct {
			Title       string
			GalleryHash string
		}{
			Title:       title,
			GalleryHash: galleryHash,
		}
		g.Templates.Edit.Execute(w, r, data, err)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	galleryHash := chi.URLParam(r, "hash")
	user := context.User(r.Context())

	title, err := g.GalleryService.GetGallery(user.ID, galleryHash)
	if err != nil {
		log.Printf("galleries: %+v", err)
		g.Templates.Index.Execute(w, r, nil, err)
	}
	data := struct {
		Title       string
		GalleryHash string
		ImagePaths  []string
	}{
		Title:       title,
		GalleryHash: galleryHash,
	}
	images, err := g.GalleryService.Images(galleryHash)
	if err != nil {
		log.Printf("galleries: %+v", err)
	}
	for _, image := range images {
		imagePath := filepath.Base(image.Path)
		data.ImagePaths = append(data.ImagePaths, imagePath)
	}
	g.Templates.Edit.Execute(w, r, data)
}

func (g Galleries) ConfirmDelete(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	confirm := r.FormValue("title.confirm")
	galleryHash := chi.URLParam(r, "hash")
	if title != confirm {
		err := apperrors.Public(errors.New("strings not equal"), "Gallery name validation failed.")
		data := struct {
			Title       string
			GalleryHash string
		}{
			Title:       confirm,
			GalleryHash: galleryHash,
		}
		g.Templates.Delete.Execute(w, r, data, err)
		return
	}

	user := context.User(r.Context())

	err := g.GalleryService.Delete(user.ID, galleryHash)
	if err != nil {
		log.Printf("galleries: %+v", err)
		g.Templates.Delete.Execute(w, r, nil, err)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	galleryHash := chi.URLParam(r, "hash")
	user := context.User(r.Context())

	var data struct {
		Title       string
		GalleryHash string
	}

	title, err := g.GalleryService.GetGallery(user.ID, galleryHash)
	if err != nil {
		log.Printf("galleries: %+v", err)

	}
	data.Title = title
	data.GalleryHash = galleryHash
	g.Templates.Delete.Execute(w, r, data)
}

func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	galleryHash := chi.URLParam(r, "hash")
	user := context.User(r.Context())

	title, err := g.GalleryService.GetGallery(user.ID, galleryHash)
	if err != nil {
		log.Printf("galleries: %+v", err)
	}

	data := struct {
		Title      string
		ImagePaths []string
	}{
		Title: title,
	}

	images, err := g.GalleryService.Images(galleryHash)
	if err != nil {
		log.Printf("galleries: %+v", err)
	}
	for _, image := range images {
		data.ImagePaths = append(data.ImagePaths, image.Path)
	}
	g.Templates.Show.Execute(w, r, data)
}

func (g Galleries) Image(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.EscapedPath()
	// dir := filepath.Dir(urlPath)
	// filename := chi.URLParam(r, "filename")
	// fullPath := filepath.Join(dir, filename)
	rel := strings.TrimPrefix(urlPath, "/")
	// log.Printf("%q %q %q %q %q", urlPath, dir, filename, fullPath, rel)
	http.ServeFile(w, r, rel)
}

func (g Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleryHash := chi.URLParam(r, "hash")
	filename := chi.URLParam(r, "filename")
	filepath := fmt.Sprintf("galleries/%s/%s", galleryHash, filename)
	log.Printf("%q %q %q", galleryHash, filename, filepath)
	err := g.GalleryService.DeleteImage(user.ID, galleryHash, filepath)
	if err != nil {
		log.Printf("di: %+v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	editURL := fmt.Sprintf("/galleries/%s/edit", galleryHash)
	http.Redirect(w, r, editURL, http.StatusFound)
}

func (g Galleries) UploadImage(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleryHash := chi.URLParam(r, "hash")
	_, err := g.GalleryService.GetGallery(user.ID, galleryHash)
	if err != nil {
		log.Printf("ui: %+v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	err = r.ParseMultipartForm(5 * 1024 * 1024)
	if err != nil {
		log.Printf("ui: %+v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	fileHeaders := r.MultipartForm.File["images"]
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("ui: %+v", err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		if err := g.GalleryService.CreateImage(galleryHash, fileHeader.Filename, file); err != nil {
			log.Printf("ui: %+v", err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		log.Printf("successfully uploaded: %q", fileHeader.Filename)
	}
	url := fmt.Sprintf("/galleries/%s/edit", galleryHash)
	http.Redirect(w, r, url, http.StatusFound)
}
