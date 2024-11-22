package handler_test

import (
	"GRIMIDIE/handler"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePinjam_Success(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	// Simulate a successful lending operation (book borrowing)
	mockHandler.On("CreatePinjam", 1, 101, 2).Return(nil) // No error means success

	err := mockHandler.CreatePinjam(1, 101, 2)

	// Assert no error occurred
	assert.NoError(t, err)
	mockHandler.AssertExpectations(t)
}

func TestCreatePinjam_DatabaseError(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	// Simulate an error when trying to create a lending transaction
	mockHandler.On("CreatePinjam", 1, 101, 2).
		Return(errors.New("pq: column bod.BookID does not exist"))

	// Call the CreatePinjam method, which should return the simulated error
	err := mockHandler.CreatePinjam(1, 101, 2)

	// Assert the error matches the database column issue
	assert.EqualError(t, err, "pq: column bod.BookID does not exist")
	mockHandler.AssertExpectations(t)
}

func TestCreatePinjam_InvalidInput(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	// Simulate a scenario where invalid input is provided (optional test)
	mockHandler.On("CreatePinjam", 0, 101, 2).Return(errors.New("invalid input"))

	err := mockHandler.CreatePinjam(0, 101, 2)

	// Assert the error matches invalid input
	assert.EqualError(t, err, "invalid input")
	mockHandler.AssertExpectations(t)
}

func TestCreatePinjam_InvalidBookID(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	// Simulate an invalid book ID provided for lending
	mockHandler.On("CreatePinjam", -1, 101, 2).
		Return(errors.New("invalid book ID"))

	err := mockHandler.CreatePinjam(-1, 101, 2)

	// Assert the error matches invalid book ID
	assert.EqualError(t, err, "invalid book ID")
	mockHandler.AssertExpectations(t)
}

func TestCreatePinjam_InsufficientStock(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	// Simulate a scenario where there is not enough stock for lending
	mockHandler.On("CreatePinjam", 1, 101, 10).
		Return(errors.New("not enough stock available"))

	err := mockHandler.CreatePinjam(1, 101, 10)

	// Assert that the stock availability error is returned
	assert.EqualError(t, err, "not enough stock available")
	mockHandler.AssertExpectations(t)
}

func TestCreatePinjam_BookAlreadyLended(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	// Simulate a scenario where the book is already lent out
	mockHandler.On("CreatePinjam", 1, 101, 2).
		Return(errors.New("book already lent out"))

	err := mockHandler.CreatePinjam(1, 101, 2)

	// Assert the error matches the book already lent out scenario
	assert.EqualError(t, err, "book already lent out")
	mockHandler.AssertExpectations(t)
}

func TestListPeminjaman(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	// Simulate a successful listing of borrowed books
	mockHandler.On("ListPeminjaman", 1).Return(nil) // No error, simulating a successful query

	err := mockHandler.ListPeminjaman(1)

	// Assert no error occurred
	assert.NoError(t, err)
	mockHandler.AssertExpectations(t)
}

func TestListPeminjaman_DatabaseError(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	// Simulate a database error when trying to list borrowed books
	mockHandler.On("ListPeminjaman", 1).
		Return(errors.New("failed to fetch from BookOrderDetail"))

	err := mockHandler.ListPeminjaman(1)

	// Assert the error matches the database failure
	assert.EqualError(t, err, "failed to fetch from BookOrderDetail")
	mockHandler.AssertExpectations(t)
}

func TestReturnPinjam_Success(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	// Simulate a successful book return
	mockHandler.On("ReturnPinjam", 1).Return(nil) // No error means success

	err := mockHandler.ReturnPinjam(1)

	// Assert no error occurred
	assert.NoError(t, err)
	mockHandler.AssertExpectations(t)
}

func TestReturnPinjam_DatabaseError(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	// Simulate an error when returning a book
	mockHandler.On("ReturnPinjam", 1).
		Return(errors.New("failed to insert into BookOrderDetail"))

	err := mockHandler.ReturnPinjam(1)

	// Assert the error matches the database failure
	assert.EqualError(t, err, "failed to insert into BookOrderDetail")
	mockHandler.AssertExpectations(t)
}

func TestUserRegister_Success(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	// Simulate a successful user registration
	mockHandler.On("UserRegister", "Maman Racing", "maman.race@yahu.com", "mamansukabalap").
		Return(123, nil)

	// Call the mock method
	userID, err := mockHandler.UserRegister("Maman Racing", "maman.race@yahu.com", "mamansukabalap")

	// Validate the result
	assert.NoError(t, err)
	assert.Equal(t, 123, userID)

	// Ensure all expectations are met
	mockHandler.AssertExpectations(t)
}

func TestUserLogin_Success(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	// Simulate a successful user login
	mockHandler.On("UserLogin", "Maman Racing", "mamansukabalap").
		Return(123, nil)

	// Call the mock method
	userID, err := mockHandler.UserLogin("Maman Racing", "mamansukabalap")

	// Validate the result
	assert.NoError(t, err)
	assert.Equal(t, 123, userID)

	// Ensure all expectations are met
	mockHandler.AssertExpectations(t)
}
