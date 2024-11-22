package handler_test

import (
	"GRIMIDIE/handler"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePinjam(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	mockHandler.On("CreatePinjam", 1, 101, 2).Return(errors.New("some error"))

	err := mockHandler.CreatePinjam(1, 101, 2)

	assert.EqualError(t, err, "some error")
	mockHandler.AssertExpectations(t)
}

func TestCreatePinjam_Negative(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	mockHandler.On("CreatePinjam", 1, 101, 2).
		Return(errors.New("failed to insert into BookOrderDetail"))

	err := mockHandler.CreatePinjam(1, 101, 2)

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to insert into BookOrderDetail")

	mockHandler.AssertExpectations(t)
}

func TestListPeminjaman(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	mockHandler.On("ListPeminjaman", 1).Return(errors.New("some error"))

	err := mockHandler.ListPeminjaman(1)

	assert.EqualError(t, err, "some error")
	mockHandler.AssertExpectations(t)
}

func TestListPeminjaman_Negative(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	mockHandler.On("ListPeminjaman", 1).
		Return(errors.New("failed to insert into BookOrderDetail"))

	err := mockHandler.ListPeminjaman(1)

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to insert into BookOrderDetail")

	mockHandler.AssertExpectations(t)
}

func TestReturnPinjam(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	mockHandler.On("ReturnPinjam", 1).Return(errors.New("some error"))

	err := mockHandler.ReturnPinjam(1)

	assert.EqualError(t, err, "some error")
	mockHandler.AssertExpectations(t)
}

func TestReturnPinjam_Negative(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	mockHandler.On("ReturnPinjam", 1).
		Return(errors.New("failed to insert into BookOrderDetail"))

	err := mockHandler.ReturnPinjam(1)

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to insert into BookOrderDetail")

	mockHandler.AssertExpectations(t)
}

func TestUserRegister(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	// Atur return value untuk method UserRegister
	mockHandler.On("UserRegister", "Maman Racing", "maman.race@yahu.com", "mamansukabalap").
		Return(123, nil)

	// Panggil fungsi mock
	userID, err := mockHandler.UserRegister("Maman Racing", "maman.race@yahu.com", "mamansukabalap")

	// Validasi hasil
	assert.NoError(t, err)
	assert.Equal(t, 123, userID)

	// Memastikan semua expectation terpenuhi
	mockHandler.AssertExpectations(t)
}

func TestUserLogin(t *testing.T) {
	mockHandler := new(handler.MockHandler)

	// Atur return value untuk method UserLogin
	mockHandler.On("UserLogin", "Maman Racing", "mamansukabalap").
		Return(123, nil)

	// Panggil fungsi mock
	userID, err := mockHandler.UserLogin("Maman Racing", "mamansukabalap")

	// Validasi hasil
	assert.NoError(t, err)
	assert.Equal(t, 123, userID)

	// Memastikan semua expectation terpenuhi
	mockHandler.AssertExpectations(t)
}
