package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string
}

type BinaryFiles struct {
	UserId   string
	Filename string
	Body     []byte
}
