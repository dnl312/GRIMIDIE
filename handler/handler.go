package handler

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/stretchr/testify/mock"
)

type MockHandler struct {
	mock.Mock
}

type Handler interface {
	UserRegister(name, email, password string) error
	CreatePinjam(UserID, BookID, Qty int) error
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
  	fmt.Println()
	fmt.Printf("%s, %s, %s\n", name, email, password)

	return nil
}

func (h *HandlerImpl) CreatePinjam (UserID, BookID, Qty int) error {

	orderDetail, err := h.DB.Exec("INSERT INTO BookOrderDetail (BookID, Quantity, TanggalPinjam, TanggalBalik, Denda) VALUES(?, ?, NOW(), '', 0);", BookID, Qty )
	
	if err != nil {
		log.Print("Error creating transaction: ", err)
	}else{
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