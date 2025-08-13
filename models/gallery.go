package models

import (
	"database/sql"
	"errors"
	"fmt"

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

var (
	ErrGalleryExists = errors.New("models: gallery already exists")
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

	_, err = gs.DB.Exec(`
	INSERT INTO galleries
	VALUES($1, $2, $3)`, gallery.UserID, gallery.Title, gallery.GalleryHash)

	var pgErr pgx.PgError
	if errors.As(err, &pgErr) {
		if pgErr.ConstraintName == "galleries_user_id_title_key" {
			return nil, ErrGalleryExists
		}
		return nil, err
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

func (gs *GalleryService) Rename(userID uuid.UUID, hash, newTitle string) error {

	_, err := gs.DB.Exec(`
	UPDATE galleries
	SET title = $3
	WHERE user_id = $1 AND hash = $2`, userID, hash, newTitle)
	var pgErr pgx.PgError
	if errors.As(err, &pgErr) {
		if pgErr.ConstraintName == "galleries_user_id_title_key" {
			return ErrGalleryExists
		}
		return err
	}

	return nil
}

func (gs *GalleryService) GetGallery(userID uuid.UUID, hash string) (string, error) {
	row := gs.DB.QueryRow(`
	SELECT title
	FROM galleries
	WHERE user_id = $1 AND hash = $2`, userID, hash)

	var title string
	err := row.Scan(&title)
	if err != nil {
		return "", fmt.Errorf("gs: %w", err)
	}
	return title, nil
}

func (gs *GalleryService) Delete(userID uuid.UUID, hash string) error {
	_, err := gs.DB.Exec(`
	DELETE FROM galleries
	WHERE user_id = $1 AND hash = $2`, userID, hash)
	if err != nil {
		return fmt.Errorf("gs: %w", err)
	}
	return nil
}
