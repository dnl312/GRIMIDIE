package handler

import "database/sql"

type Handler interface {
	UserRegister(name, email, password string) error
	UserLogin(email, password string) error
}

type HndlerImpl struct {
	DB *sql.DB
}

func NewHndler(DB *sql.DB) *HndlerImpl {
	return &HndlerImpl{
		DB: DB,
	}
}
