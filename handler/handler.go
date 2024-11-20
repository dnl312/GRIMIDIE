package handler

import (
	"database/sql"
	"fmt"
)

type Handler interface {
	UserRegister(name, email, password string) error
	//UserLogin(email, password string) error
}

type HandlerImpl struct {
	DB *sql.DB
}

func NewHandler(DB *sql.DB) *HandlerImpl {
	return &HandlerImpl{
		DB: DB,
	}
}


func (h *HandlerImpl) UserRegister(name, email, password string) error{
	fmt.Println("\n")
	fmt.Printf("%s, %s, %s\n", name, email, password)

	return nil
}