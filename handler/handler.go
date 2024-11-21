package handler

import (
	"database/sql"
	"log"

	"github.com/stretchr/testify/mock"
)

type MockHandler struct {
	mock.Mock
}

type Handler interface {
	UserRegister(name, email, password string) error
	CreatePinjam(UserID, BookID, Qty int) error
}

type HandlerImpl struct {
	DB *sql.DB
}

func NewHandler(DB *sql.DB) *HandlerImpl {
	return &HandlerImpl{
		DB: DB,
	}
}

func (h *HandlerImpl) UserRegister(name, email, password string) error {
	_, err := h.DB.Exec("INSERT INTO Users (name, email, password) VALUES ($1,$2,$3)", name, email, password)
	if err != nil {
		log.Print("Error inserting record: ", err)
		return err
	}

	log.Print("Record inserted successfully")
	return nil
}

func (h *HandlerImpl) CreatePinjam(UserID, BookID, Qty int) error {

	orderDetail, err := h.DB.Exec("INSERT INTO BookOrderDetail (BookID, Quantity, TanggalPinjam, TanggalBalik, Denda) VALUES(?, ?, NOW(), '', 0);", BookID, Qty)

	if err != nil {
		log.Print("Error creating transaction: ", err)
	} else {
		orderDtlId, _ := orderDetail.LastInsertId()
		_, err = h.DB.Exec("INSERT INTO BookOrders(UserID, BookOrderDetailID) VALUES(?, ?)", UserID, orderDtlId)
		if err != nil {
			log.Print("Error creating transaction: ", err)
		}
	}

	log.Print("Transaction inserted successfully")
	return nil
}

func (m *MockHandler) CreatePinjam(UserID, BookID, Qty int) error {
	args := m.Called(UserID, BookID, Qty)
	return args.Error(0)
}
