package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string       `gorm:"size:63;not null"`
	BinaryFiles []BinaryFile // has many (0..N)
}
