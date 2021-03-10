package models

import "github.com/jinzhu/gorm"

// Gallery is our image contatiner resource that visitors view.
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

// GalleryService ...
type GalleryService interface {
	GalleryDB
}

// GalleryDB ...
type GalleryDB interface {
	Create(gallery *Gallery) error
}

type galleryValidator struct {
	GalleryDB
}

var _ GalleryDB = &galleryGorm{}

// NewGalleryService ...
func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{&galleryGorm{db}},
	}
}

type galleryService struct {
	GalleryDB
}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}
