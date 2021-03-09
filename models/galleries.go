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

type galleryGorm struct {
	db *gorm.DB
}
