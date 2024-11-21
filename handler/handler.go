package handler

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	"github.com/stretchr/testify/mock"
)

// MockHandler untuk testing
type MockHandler struct {
	mock.Mock
}

// Handler interface
type Handler interface {
	UserRegister(name, email, password string) error
	UserLogin(email, password string) error
	CreatePinjam(UserID, BookID, Qty int) error
	DeleteReturnedTransactions() error
}

// HandlerImpl implementasi handler
type HandlerImpl struct {
	DB *sql.DB
}

// Buat instance HandlerImpl
func NewHandler(DB *sql.DB) *HandlerImpl {
	return &HandlerImpl{
		DB: DB,
	}
}

// Helper function to validate password
func isValidPassword(password string) bool {
	// Example: password must be at least 8 characters long
	if len(password) < 8 {
		return false
	}
	return true
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	// Simple regex for email validation
	const emailRegex = `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// Helper function to validate name (alphabetic only)
func isValidName(name string) bool {
	// Simple regex for name validation (only alphabetic characters and spaces)
	const nameRegex = `^[a-zA-Z\s]+$`
	re := regexp.MustCompile(nameRegex)
	return re.MatchString(name)
}

// Daftarkan user baru
func (h *HandlerImpl) UserRegister(Nama, Email, Password string) error {
	// Validate Name
	if Nama == "" {
		return fmt.Errorf("nama tidak boleh kosong")
	}

	if !isValidName(Nama) {
		return fmt.Errorf("nama hanya boleh berisi huruf dan spasi")
	}

	// Validate Email
	if !isValidEmail(Email) {
		fmt.Println("Invalid email format:", Email) // Debugging log
		return fmt.Errorf("email tidak valid")
	}

	// Validate Password
	if !isValidPassword(Password) {
		return fmt.Errorf("password harus minimal 8 karakter")
	}

	// Check if email is already registered
	var existingUserCount int
	err := h.DB.QueryRow(`SELECT COUNT(*) FROM "Users" WHERE "Email" = $1`, Email).Scan(&existingUserCount)
	if err != nil {
		log.Print("Gagal memeriksa email di database: ", err)
		return err
	}

	if existingUserCount > 0 {
		return fmt.Errorf("email sudah terdaftar")
	}

	// Insert new user into database
	_, err = h.DB.Exec(`INSERT INTO "Users" ("Nama", "Email", "Password") VALUES ($1, $2, $3)`, Nama, Email, Password)
	if err != nil {
		log.Print("Gagal menambahkan user: ", err)
		return err
	}

	log.Print("User berhasil ditambahkan")
	return nil
}

// Prompt for user details with validation as they are entered
func (h *HandlerImpl) PromptUserInput() (string, string, string, error) {
	var name, email, password string

	// Prompt user for Name and validate immediately
	fmt.Print("Enter Name: ")
	fmt.Scanln(&name)
	if !isValidName(name) {
		return "", "", "", fmt.Errorf("nama hanya boleh berisi huruf dan spasi")
	}

	// Prompt user for Email and validate immediately
	fmt.Print("Enter Email: ")
	fmt.Scanln(&email)
	if !isValidEmail(email) {
		return "", "", "", fmt.Errorf("email tidak valid")
	}

	// Prompt user for Password and validate immediately
	fmt.Print("Enter Password: ")
	fmt.Scanln(&password)
	if !isValidPassword(password) {
		return "", "", "", fmt.Errorf("password harus minimal 8 karakter")
	}

	return name, email, password, nil
}

// Login user
func (h *HandlerImpl) UserLogin(email, password string) error {
	var storedPassword string

	query := `SELECT "Password" FROM "Users" WHERE "Email" = $1`
	err := h.DB.QueryRow(query, email).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user tidak ditemukan")
		}
		return fmt.Errorf("error database: %v", err)
	}

	if storedPassword != password {
		return fmt.Errorf("password salah")
	}

	fmt.Printf("Login berhasil\n")
	return nil
}

// Tambah transaksi peminjaman buku
func (h *HandlerImpl) CreatePinjam(UserID, BookID, Qty int) error {
	var orderDtlId int64

	err := h.DB.QueryRow(`INSERT INTO "BookOrderDetail" ("BookID", "Quantity", "TanggalPinjam", "TanggalBalik", "Denda") 
						VALUES ($1, $2, NOW(), NULL, 0) 
						RETURNING "ID"`,
		BookID, Qty).Scan(&orderDtlId)

	if err != nil {
		log.Print("Gagal membuat detail transaksi: ", err)
		return err
	}

	_, err = h.DB.Exec(`INSERT INTO "BookOrders" ("UserID", "BookOrderDetailID") VALUES($1, $2)`, UserID, orderDtlId)
	if err != nil {
		log.Print("Gagal membuat transaksi peminjaman: ", err)
		return err
	}

	log.Print("Transaksi berhasil ditambahkan")
	return nil
}

// Hapus transaksi yang buku-nya sudah dikembalikan
func (h *HandlerImpl) DeleteReturnedTransactions() error {
	result, err := h.DB.Exec(`DELETE FROM "BookOrders" 
		WHERE "BookOrderDetailID" IN (
			SELECT "ID" 
			FROM "BookOrderDetail" 
			WHERE "TanggalBalik" IS NOT NULL
		)`)

	if err != nil {
		log.Printf("Gagal menghapus transaksi: %v", err)
		return err
	}

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error retrieving rows affected: %v", err)
		return err
	}

	log.Printf("Berhasil menghapus %d transaksi yang sudah selesai.", rowsDeleted)
	return nil
}

// Mock untuk testing CreatePinjam
func (m *MockHandler) CreatePinjam(UserID, BookID, Qty int) error {
	args := m.Called(UserID, BookID, Qty)
	return args.Error(0)
}

// Mock untuk testing DeleteReturnedTransactions
func (m *MockHandler) DeleteReturnedTransactions() error {
	args := m.Called()
	return args.Error(0)
}
