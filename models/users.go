package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User ...
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
