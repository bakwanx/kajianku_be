package controller

import (
	"kajianku_be/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterMosque(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		mosque := model.Mosque{}
		c.Bind(&mosque)

		checkMosque := model.Mosque{}
		db.Where("id_user  = ?", mosque.IdUser).First(&checkMosque)
		if checkMosque.MosqueName != "" {
			return c.JSON(http.StatusOK, echo.Map{
				"status":  http.StatusOK,
				"message": "This mosque is already registered",
			})
		}

		err := db.Create(&mosque).Error
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"status":  http.StatusOK,
			"message": "success",
			"data":    mosque,
		})
	}
}

func GetMosqueByUserId(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		idUser := c.Param("id_user")
		mosque := model.Mosque{}

		errMosque := db.Where("id_user  = ?", idUser).First(&mosque).Error
		if errMosque != nil {
			return c.JSON(http.StatusNotFound, echo.Map{
				"status":  http.StatusNotFound,
				"message": "mosque not found",
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status":  http.StatusOK,
			"message": "success",
			"data":    mosque,
		})
	}
}
