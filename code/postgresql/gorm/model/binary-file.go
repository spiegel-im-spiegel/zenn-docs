package model

import "gorm.io/gorm"

type BinaryFile struct {
	gorm.Model
	UserId   uint   `gorm:"not null"`
	Filename string `gorm:"not null"`
	Body     []byte
}
