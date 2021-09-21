package model

import "gorm.io/gorm"

type BinaryFile struct {
	gorm.Model
	UserId   string
	Filename string
	Body     []byte
}
