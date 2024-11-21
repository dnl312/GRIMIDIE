package handler

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockHandler struct {
	mock.Mock
}

type Handler interface {
	UserRegister(name, email, password string) error
	UserLogin(email, password string) error
	ListBooks() error
	CreatePinjam(UserID, BookID, Qty int) error
	ReturnPinjam(BookOrderID int) error
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

	fmt.Printf("Sign In Successful! \n")
	return nil
}

func (h *HandlerImpl) ListBooks() error {
	rows, err := h.DB.Query(`SELECT * FROM "Books"`)
	if err != nil {
		log.Print("Error listing books: ", err)
		return err
	}
	defer rows.Close()

	// Adjusted separator length
	fmt.Println(strings.Repeat("-", 89))
	fmt.Printf("| %-3s | %-25s | %-20s | %-15s | %-10s |\n", "ID", "BOOK TITLE", "AUTHOR", "PUBLISH DATE", "STOCK")
	fmt.Println(strings.Repeat("-", 89))

	for rows.Next() {
		var id int
		var name string
		var author string
		var publishDate time.Time
		var stock int

		if err := rows.Scan(&id, &name, &author, &publishDate, &stock); err != nil {
			return fmt.Errorf("database scanning rows: %v", err)
		}
		fmt.Printf("| %-3d | %-25s | %-20s | %-15s | %-10d |\n", id, name, author, publishDate.Format("2006-01-02"), stock)
	}
	fmt.Println(strings.Repeat("-", 89))

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error scanning rows: %v", err)
	}
	return nil
}

func (h *HandlerImpl) CreatePinjam(UserID, BookID, Qty int) error {

	var orderDtlId int64

	err := h.DB.QueryRow(`INSERT INTO "BookOrderDetail" ("BookID", "Quantity", "TanggalPinjam", "TanggalBalik", "Denda") 
						VALUES ($1, $2, NOW(), NULL, 0) 
						RETURNING ID`,
		BookID, Qty).Scan(&orderDtlId)

	if err != nil {
		log.Print("Error creating Book Order Detail transaction: ", err)

	} else {
		_, err = h.DB.Exec(`INSERT INTO "BookOrders" ("UserID", "BookOrderDetailID") VALUES($1, $2)`, UserID, orderDtlId)

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

func (h *HandlerImpl) ListPeminjaman(UserID int) error {
	rows, err := h.DB.Query(`SELECT bo."OrderID", b."JudulBuku" , bod."TanggalPinjam" FROM "BookOrders" bo
							left join "BookOrderDetail" bod on bod."BookOrderDetailID" = bo."BookOrderDetailID" 
							left join "Books" b ON b."BookID" = bod."BookID" where bod."UserID" = $1`, UserID)
	if err != nil {
		log.Print("Error fetching records: ", err)
		return err
	}
	defer rows.Close()

	fmt.Println("ID\tJudul Buku\tTanggal Pinjam")
	for rows.Next() {
		var OrderID int
		var JudulBuku, TanggalPinjam string
		err = rows.Scan(&OrderID, &JudulBuku, &TanggalPinjam)
		if err != nil {
			log.Print("Error scanning record: ", err)
			return err
		}

		fmt.Printf("%d\t%s\t%s\n", OrderID, JudulBuku, TanggalPinjam)
	}

	return nil
}

func (m *MockHandler) ListPeminjaman(UserID int) error {
	args := m.Called(UserID)
	return args.Error(0)
}

func (h *HandlerImpl) ReturnPinjam(BookOrderID int) error {

	rows, err := h.DB.Query(`SELECT (CAST(NOW() AS date) - CAST("TanggalPinjam" AS date)) as DateDifference FROM "BookOrders" bo
							left join "BookOrderDetail" bod on bod."BookOrderDetailID" = bo."BookOrderDetailID" 
							where bo."OrderID" = $1`, BookOrderID)
	if err != nil {
		log.Print("Error Fetch Book Order transaction: ", err)
	}

	for rows.Next() {
		var BookOrderDetailID, DateDifference int
		err = rows.Scan(&BookOrderDetailID, &DateDifference)
		if err != nil {
			log.Print("Error scanning record: ", err)
			return err
		}

		if DateDifference > 7 {
			_, err := h.DB.Query(`UPDATE public."BookOrderDetail"
							SET  "TanggalBalik"=NOW(), "Denda"= $1
							WHERE "BookOrderDetailID"= $2`, DateDifference*5000, BookOrderDetailID)
			if err != nil {
				log.Print("Error scanning record: ", err)
				return err
			}
		} else {
			//delete query fahri
		}
	}

	log.Print("Transaction inserted successfully")
	return nil
}

func (m *MockHandler) ReturnPinjam(BookOrderID int) error {
	args := m.Called(BookOrderID)
	return args.Error(0)
}
