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
	UserLogin(email, password string) error
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

func (h *HandlerImpl) UserRegister(Nama, Email, Password string) error {
	_, err := h.DB.Exec(`INSERT INTO "Users" ("Nama", "Email", "Password") VALUES ($1,$2,$3)`, Nama, Email, Password)
	if err != nil {
		log.Print("Error inserting record: ", err)
		return err
	}

	log.Print("Record inserted successfully")
	return nil
}

func (h *HandlerImpl) UserLogin(email, password string) error {
	var storedPassword string

	query := `SELECT "Password" FROM "Users" WHERE "Email" = $1`
	err := h.DB.QueryRow(query, email).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("database error: %v", err)
	}

	if storedPassword != password {
		return fmt.Errorf("invalid password")
	}

	fmt.Printf("Sign in successful! \n")
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
