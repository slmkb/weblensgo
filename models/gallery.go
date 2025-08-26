package models

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/slmkb/weblensgo/rand"
)

type Gallery struct {
	UserID      uuid.UUID
	Title       string
	GalleryHash string
}

type GalleryService struct {
	DB *sql.DB
}

type Image struct {
	Path string
}

var (
	ErrGalleryExists = errors.New("models: gallery already exists")
	ErrDirExists     = errors.New("models: target directory already exists")
	extensions       = []string{".png", ".jpg", ".jpeg", ".gif"}
)

const (
	galleriesDir = "galleries"
)

func (gs *GalleryService) Create(usedID uuid.UUID, title string) (*Gallery, error) {
	galleryHash, err := rand.GalleryHash()
	if err != nil {
		return nil, fmt.Errorf("gs: %w", err)
	}
	gallery := Gallery{
		UserID:      usedID,
		Title:       title,
		GalleryHash: galleryHash,
	}

	galleryDir := galleryDir(galleryHash)
	if _, err = os.Stat(galleryDir); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("gs: %w", err)
		}
	}

	_, err = gs.DB.Exec(`
	INSERT INTO galleries
	VALUES($1, $2, $3)`, gallery.UserID, gallery.Title, gallery.GalleryHash)

	var pgErr pgx.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.ConstraintName {
		case "galleries_user_id_title_key":
			return nil, ErrGalleryExists
		case "galleries_hash_key":
			log.Printf("gs: gallery hash collision")
			fallthrough
		default:
			return nil, err
		}
	}

	err = os.Mkdir(galleryDir, 0770)
	if err != nil {
		return nil, fmt.Errorf("gs: %w", err)
	}

	return &gallery, nil
}

func (gs *GalleryService) GetGalleries(userID uuid.UUID) (*[]Gallery, error) {

	rows, err := gs.DB.Query(`
	SELECT title, hash 
	FROM galleries
	WHERE user_id = $1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("gs: %w", err)
	}
	defer rows.Close()

	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{
			UserID: userID,
		}
		if err := rows.Scan(&gallery.Title, &gallery.GalleryHash); err != nil {
			return nil, fmt.Errorf("gs: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gs: %w", err)
	}

	return &galleries, nil
}

func (gs *GalleryService) Edit(userID uuid.UUID, galleryHash, newTitle string) error {

	_, err := gs.DB.Exec(`
	UPDATE galleries
	SET title = $3
	WHERE user_id = $1 AND hash = $2`, userID, galleryHash, newTitle)
	var pgErr pgx.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.ConstraintName {
		case "galleries_user_id_title_key":
			return ErrGalleryExists
		default:
			return err
		}
	}

	return nil
}

func (gs *GalleryService) GetGallery(userID uuid.UUID, galleryHash string) (string, error) {
	row := gs.DB.QueryRow(`
	SELECT title
	FROM galleries
	WHERE user_id = $1 AND hash = $2`, userID, galleryHash)

	var title string
	err := row.Scan(&title)
	if err != nil {
		return "", fmt.Errorf("gs: %w", err)
	}
	return title, nil
}

func (gs *GalleryService) Delete(userID uuid.UUID, galleryHash string) error {

	galleryDir := galleryDir(galleryHash)

	_, err := gs.DB.Exec(`
	DELETE FROM galleries
	WHERE user_id = $1 AND hash = $2`, userID, galleryHash)
	if err != nil {
		return fmt.Errorf("gs: %w", err)
	}

	if _, err := os.Stat(galleryDir); err != nil {
		return fmt.Errorf("gs: %w", err)
	}

	err = os.RemoveAll(galleryDir)
	if err != nil {
		return fmt.Errorf("gs: %w", err)
	}

	return nil
}

func (gs *GalleryService) Images(galleryHash string) ([]Image, error) {
	dir := galleryDir(galleryHash)
	parent := filepath.Dir(dir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("gs: read dir %q: %w", dir, err)
	}

	var images []Image
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		full := filepath.Join(dir, e.Name())
		if !isImage(full) {
			continue
		}

		rel, err := filepath.Rel(parent, full)
		if err != nil {
			continue
		}

		fileName := url.PathEscape(filepath.Base(rel))
		fileDir := filepath.Clean(filepath.Dir(rel))
		fullPath := filepath.Join(fileDir, fileName)
		images = append(images, Image{
			Path: fullPath,
		})

	}
	return images, nil
}

func (gs *GalleryService) DeleteImage(userID uuid.UUID, galleryHash, filepath string) error {
	_, err := gs.GetGallery(userID, galleryHash)
	if err != nil {
		return fmt.Errorf("gs: %w", err)
	}
	err = os.Remove(filepath)
	if err != nil {
		return fmt.Errorf("gs: %w", err)
	}
	return nil
}

func (gs *GalleryService) CreateImage(galleryHash, filename string, contents io.ReadSeeker) error {
	if err := checkContentType(contents, []string{"image/jpeg"}); err != nil {
		return fmt.Errorf("ci: %w", err)
	}
	dir := galleryDir(galleryHash)
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf("ci: %w", err)
	}
	filePath := filepath.Join(dir, filename)
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists %q", filePath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("ci: %w", err)
	}
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("ci: %w", err)
	}
	if _, err := io.Copy(file, contents); err != nil {
		return fmt.Errorf("ci: %w", err)
	}

	return nil
}

func galleryDir(galleryHash string) string {
	return filepath.Join(galleriesDir, galleryHash)
}

func isImage(file string) bool {
	fileExtension := strings.ToLower(filepath.Ext(file))
	for _, extension := range extensions {
		if extension == fileExtension {
			return true
		}
	}
	return false
}
