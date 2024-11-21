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
