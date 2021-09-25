package model

import "gorm.io/gorm"

type BinaryFile struct {
	gorm.Model
	UserId   uint   `gorm:"not null"`
	Filename string `gorm:"unique;not null"`
	Body     []byte
}
