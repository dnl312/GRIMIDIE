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

	var orderDtlId int64
	
	err := h.DB.QueryRow(`INSERT INTO BookOrderDetail (BookID, Quantity, TanggalPinjam, TanggalBalik, Denda) 
						VALUES ($1, $2, NOW(), NULL, 0) 
						RETURNING ID;`,
						BookID, Qty).Scan(&orderDtlId)

	//_ , err = h.DB.Query("INSERT INTO BookOrderDetail (BookID, Quantity, TanggalPinjam, TanggalBalik, Denda) VALUES ($1, $2, NOW(), NULL, 0) ", BookID, Qty)
	
	if err != nil {
		log.Print("Error creating Book Order Detail transaction: ", err)
	}else{
		_, err = h.DB.Exec("INSERT INTO BookOrders(UserID, BookOrderDetailID) VALUES($1, $2)", UserID, orderDtlId)

		if err != nil {
			log.Print("Error creating Book Orders transaction: ", err)
		}
	}

	log.Print("Transaction inserted successfully")
	return nil
}

func (m *MockHandler) CreatePinjam(UserID, BookID, Qty int) error {
	args := m.Called(UserID, BookID, Qty)
	return args.Error(0)
}
