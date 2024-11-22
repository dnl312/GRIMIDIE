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
	UserRegister(name, email, password string) (int, error)
	UserLogin(email, password string) (int, bool, error)
	ListBooks() error
	CreatePinjam(UserID, BookID, Qty int) error
	ReturnPinjam(UserID, BookOrderID int) (float64, error)
	ListPeminjaman(UserID int) error
	ReportPeminjaman() error
	ReportStock() error
	ReportPopularBooks() error
}

type HandlerImpl struct {
	DB *sql.DB
}

func NewHandler(DB *sql.DB) *HandlerImpl {
	return &HandlerImpl{
		DB: DB,
	}
}

func (h *HandlerImpl) UserRegister(Nama, Email, Password string) (int, error) {
	var UserID int
	err := h.DB.QueryRow(`INSERT INTO "Users" ("Nama", "Email", "Password") VALUES ($1,$2,$3) RETURNING "UserID"`, Nama, Email, Password).Scan(&UserID)
	if err != nil {
		log.Print("Error inserting record: ", err)
		return 0, err
	}

	log.Print("Record inserted successfully")
	return UserID, nil
}

func (m *MockHandler) UserRegister(Nama, Email, Password string) (int, error) {
	args := m.Called(Nama, Email, Password)
	return args.Int(0), args.Error(1)
}

func (h *HandlerImpl) UserLogin(email, password string) (int, bool, error) {
	var storedPassword string
	var UserID int
	var IsAdmin bool

	query := `SELECT "UserID", "Password", "IsAdmin" FROM "Users" WHERE "Email" = $1`
	err := h.DB.QueryRow(query, email).Scan(&UserID, &storedPassword, &IsAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false, fmt.Errorf("user not found")
		}
		return 0, false, fmt.Errorf("database error: %v", err)
	}

	if storedPassword != password {
		return 0, false, fmt.Errorf("invalid password")
	}

	fmt.Printf("Sign In Successful! \n")
	return UserID, IsAdmin, nil
}

func (m *MockHandler) UserLogin(email, password string) (int, error) {
	args := m.Called(email, password)
	return args.Int(0), args.Error(1)
}
func (h *HandlerImpl) ListBooks() error {
	rows, err := h.DB.Query(`
SELECT
    b."BookID",
    b."JudulBuku",
    b."Pengarang",
    b."PublishDate",
    (b."StokBuku" - COALESCE(COUNT(ub."BookID"), 0))
FROM
    public."Books" b
LEFT JOIN
    public."UserBooks" ub
ON
    b."BookID" = ub."BookID" AND ub."TanggalBalik" IS NULL
GROUP BY
    b."BookID", b."JudulBuku", b."Pengarang", b."PublishDate", b."StokBuku"
ORDER BY
    b."BookID" ASC;
`)
	if err != nil {
		log.Print("Error listing books: ", err)
		return err
	}
	defer rows.Close()

	// Adjusted separator length
	fmt.Println(strings.Repeat("-", 81))
	fmt.Printf("| %-3s | %-25s | %-20s | %-12s | %-5s |\n", "ID", "BOOK TITLE", "AUTHOR", "PUBLISH DATE", "STOCK")
	fmt.Println(strings.Repeat("-", 81))

	for rows.Next() {
		var id int
		var name string
		var author string
		var publishDate time.Time
		var stock int

		if err := rows.Scan(&id, &name, &author, &publishDate, &stock); err != nil {
			return fmt.Errorf("database scanning rows: %v", err)
		}
		fmt.Printf("| %-3d | %-25s | %-20s | %-12s | %-5d |\n", id, name, author, publishDate.Format("2006-01-02"), stock)
	}
	fmt.Println(strings.Repeat("-", 81))

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error scanning rows: %v", err)
	}
	return nil
}

func (h *HandlerImpl) CreatePinjam(UserID, BookID, Qty int) error {
	// Check if the requested quantity is available in stock
	var stock int
	err := h.DB.QueryRow(`SELECT "StokBuku" FROM "Books" WHERE "BookID" = $1`, BookID).Scan(&stock)
	if err != nil {
		log.Print("Error checking stock: ", err)
		return fmt.Errorf("error checking book stock")
	}

	// If there aren't enough books in stock
	if stock < Qty {
		log.Print("Not enough stock available")
		return fmt.Errorf("not enough stock available")
	}

	// Begin a transaction to ensure both inserts are atomic
	tx, err := h.DB.Begin()
	if err != nil {
		log.Print("Error starting transaction: ", err)
		return err
	}
	defer tx.Rollback() // Ensure rollback if anything fails

	// Insert into BookOrderDetail
	var orderDtlId int64
	err = tx.QueryRow(`INSERT INTO "BookOrderDetail" ("BookID", "Quantity", "TanggalPinjam", "TanggalBalik", "Denda") 
						VALUES ($1, $2, NOW(), NULL, 0) 
						RETURNING "BookOrderDetailID"`, BookID, Qty).Scan(&orderDtlId)

	if err != nil {
		log.Print("Error creating Book Order Detail transaction: ", err)
		return err
	}

	// Insert into BookOrders
	_, err = tx.Exec(`INSERT INTO "BookOrders" ("UserID", "BookOrderDetailID") VALUES($1, $2)`, UserID, orderDtlId)
	if err != nil {
		log.Print("Error creating Book Orders transaction: ", err)
		return err
	}

	// Decrease the stock of books after the loan
	_, err = tx.Exec(`UPDATE "Books" SET "StokBuku" = "StokBuku" - $1 WHERE "BookID" = $2`, Qty, BookID)
	if err != nil {
		log.Print("Error updating book stock: ", err)
		return err
	}

	// Commit the transaction if everything went well
	if err := tx.Commit(); err != nil {
		log.Print("Error committing transaction: ", err)
		return err
	}

	log.Print("Book loan transaction completed successfully")
	return nil
}

func (m *MockHandler) CreatePinjam(UserID, BookID, Qty int) error {
	args := m.Called(UserID, BookID, Qty)
	return args.Error(0)
}

func (h *HandlerImpl) ListPeminjaman(UserID int) error {
	rows, err := h.DB.Query(`SELECT bo."OrderID", b."JudulBuku" , bod."TanggalPinjam" FROM "BookOrders" bo
							left join "BookOrderDetail" bod on bod."BookOrderDetailID" = bo."BookOrderDetailID" 
							left join "Books" b ON b."BookID" = bod."BookID" where bo."UserID" = $1 AND bod."TanggalBalik" IS NULL`, UserID)
	if err != nil {
		log.Print("Error fetching records: ", err)
		return err
	}
	defer rows.Close()

	fmt.Println("ID\tJudul Buku\tTanggal Pinjam")
	for rows.Next() {
		var OrderID int
		var JudulBuku string
		var TanggalPinjam time.Time
		err = rows.Scan(&OrderID, &JudulBuku, &TanggalPinjam)
		if err != nil {
			log.Print("Error scanning record: ", err)
			return err
		}

		dateOnly := TanggalPinjam.Format("2006-01-02")
		fmt.Printf("%d\t%s\t%s\n", OrderID, JudulBuku, dateOnly)
	}
	fmt.Println()

	return nil
}

func (m *MockHandler) ListPeminjaman(UserID int) error {
	args := m.Called(UserID)
	return args.Error(0)
}

func (h *HandlerImpl) ReportPeminjaman() error {
	rows, err := h.DB.Query(`SELECT bo."OrderID", b."JudulBuku" , u."Nama", bod."TanggalPinjam" FROM "BookOrders" bo
							left join "BookOrderDetail" bod on bod."BookOrderDetailID" = bo."BookOrderDetailID" 
							left join "Books" b ON b."BookID" = bod."BookID"
							left join "Users" u ON u."UserID" = bo."UserID"`)
	if err != nil {
		log.Print("Error fetching records: ", err)
		return err
	}
	defer rows.Close()

	fmt.Println("ID\tNama User\tTanggal Pinjam\tJudul Buku")
	fmt.Println(strings.Repeat("-", 80))
	for rows.Next() {
		var OrderID int
		var JudulBuku, Nama string
		var TanggalPinjam time.Time
		err = rows.Scan(&OrderID, &JudulBuku, &Nama, &TanggalPinjam)
		if err != nil {
			log.Print("Error scanning record: ", err)
			return err
		}

		dateOnly := TanggalPinjam.Format("2006-01-02")
		fmt.Printf("%d\t%s\t\t%s\t%s\n", OrderID, Nama, dateOnly, JudulBuku)
	}
	fmt.Println()

	return nil
}

func (h *HandlerImpl) ReturnPinjam(UserID, BookOrderID int) (float64, error) {
	var Denda float64

	// Query to get the book order details
	rows, err := h.DB.Query(`
		SELECT bo."BookOrderDetailID", bod."BookID", 
		       (CAST(NOW() AS date) - CAST(bod."TanggalPinjam" AS date)) as DateDifference, 
		       bod."Quantity" 
		FROM "BookOrders" bo
		LEFT JOIN "BookOrderDetail" bod ON bod."BookOrderDetailID" = bo."BookOrderDetailID" 
		WHERE bo."UserID" = $1 and bo."OrderID" = $2`, UserID, BookOrderID)

	if err != nil {
		log.Print("Error fetching book order transaction: ", err)
		return 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var BookOrderDetailID, BookID, DateDifference, Quantity int
		err = rows.Scan(&BookOrderDetailID, &BookID, &DateDifference, &Quantity)

		if err != nil {
			log.Print("Error scanning record: ", err)
			return 0, err
		}

		// Calculate the fine if the return is overdue
		if DateDifference > 7 {
			Denda = float64(DateDifference * 5000)

			// Update the return date and apply the fine
			_, err = h.DB.Exec(`UPDATE "BookOrderDetail"
				SET "TanggalBalik" = NOW(), "Denda" = $1
				WHERE "BookOrderDetailID" = $2`, Denda, BookOrderDetailID)
			if err != nil {
				log.Print("Error updating BookOrderDetail: ", err)
				return 0, err
			}
		}

		// Update the stock of the returned book
		_, err = h.DB.Exec(`
			UPDATE "Books" 
			SET "StokBuku" = "StokBuku" + $1 
			WHERE "BookID" = $2`, Quantity, BookID)
		if err != nil {
			log.Print("Error updating Books table: ", err)
			return 0, err
		}

		// Delete the record from BookOrderDetail after the return is processed
		if DateDifference < 7 {
			_, err = h.DB.Exec(`DELETE FROM "BookOrders" WHERE "BookOrderDetailID" = $1`, BookOrderDetailID)
			if err != nil {
				log.Print("Error deleting from BookOrderDetail: ", err)
				return 0, err
			}

			_, err = h.DB.Exec(`DELETE FROM "BookOrderDetail" WHERE "BookOrderDetailID" = $1`, BookOrderDetailID)
			if err != nil {
				log.Print("Error deleting from BookOrderDetail: ", err)
				return 0, err
			}

			log.Print("Returning book successfully and deleting transaction details")
			return Denda, nil
		}
	}

	log.Print("Returning book successfully")
	return Denda, nil
}

func (m *MockHandler) ReturnPinjam(BookOrderID int) error {
	args := m.Called(BookOrderID)
	return args.Error(0)
}

func (h *HandlerImpl) ReportStock() error {
	rows, err := h.DB.Query(`
SELECT
    b."JudulBuku",
    b."StokBuku",
    COALESCE(COUNT(ub."BookID"), 0) ,
    (b."StokBuku" - COALESCE(COUNT(ub."BookID"), 0))
FROM
    public."Books" b
LEFT JOIN
    public."UserBooks" ub
ON
    b."BookID" = ub."BookID" AND ub."TanggalBalik" IS NULL
GROUP BY
    b."BookID", b."JudulBuku", b."StokBuku"
ORDER BY
    b."JudulBuku";
`)
	if err != nil {
		log.Print("Error listing books: ", err)
		return err
	}
	defer rows.Close()

	// Adjusted separator length
	fmt.Println(strings.Repeat("-", 66))
	fmt.Printf("| %-25s | %-6s | %-10s | %-12s |\n", "BOOK TITLE", "STOCK", "LOAN BOOK", "ACTUAL STOCK")
	fmt.Println(strings.Repeat("-", 66))

	for rows.Next() {
		var title string
		var Stock, loanBook, finalStock int

		if err := rows.Scan(&title, &Stock, &loanBook, &finalStock); err != nil {
			return fmt.Errorf("database scanning rows: %v", err)
		}
		fmt.Printf("| %-25s | %-7d| %-10d | %-12d |\n", title, Stock, loanBook, finalStock)
	}
	fmt.Println(strings.Repeat("-", 66))

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error scanning rows: %v", err)
	}
	return nil
}

func (h *HandlerImpl) ReportPopularBooks() error {
	rows, err := h.DB.Query(`
SELECT b."JudulBuku",
       COUNT(ub."UserBookID") AS "TotalDipinjam"
FROM "UserBooks" ub
JOIN "Books" b ON b."BookID" = ub."BookID"
GROUP BY b."JudulBuku"
ORDER BY "TotalDipinjam" DESC;
`)
	if err != nil {
		log.Print("Error listing books: ", err)
		return err
	}
	defer rows.Close()

	// Adjusted separator length
	fmt.Println(strings.Repeat("-", 38))
	fmt.Printf("| %-25s | %-6s |\n", "BOOK TITLE", "TOTAL")
	fmt.Println(strings.Repeat("-", 38))

	for rows.Next() {
		var title string
		var loanBook int

		if err := rows.Scan(&title, &loanBook); err != nil {
			return fmt.Errorf("database scanning rows: %v", err)
		}
		fmt.Printf("| %-25s | %-7d|\n", title, loanBook)
	}
	fmt.Println(strings.Repeat("-", 38))

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error scanning rows: %v", err)
	}
	return nil
}
