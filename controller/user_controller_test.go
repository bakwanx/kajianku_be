package controller

import (
	"encoding/json"
	"errors"
	"kajianku_be/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetAllUsers(t *testing.T) {
	e := echo.New()
	mockDb := new(model.MockDB)

	mockDb.On("Find", mock.Anything, mock.Anything).Return(&gorm.DB{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	c.SetPath("/users")

	handler := GetAllUsers(mockDb)
	err := handler(c)

	assert.Nil(t, err)

	expectedResponse := echo.Map{
		"status":  http.StatusOK,
		"message": "success",
		"data":    []model.User{},
	}

	jsonString, err := ToJSONString(expectedResponse)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.JSONEq(t, res.Body.String(), string(jsonString))

	mockDb.AssertCalled(t, "Find", mock.Anything, mock.Anything)

	e = echo.New()

	mockDb = new(model.MockDB)

	expectedError := errors.New("database error")
	mockDb.On("Find", mock.Anything, mock.Anything).Return(&gorm.DB{
		Error: expectedError,
	})

}

func ToJSONString(v interface{}) (string, error) {

	jsonString, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}
