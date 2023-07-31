package controller

import (
	"kajianku_be/middleware"
	"kajianku_be/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func LoginAdmin(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userAdmin := model.UserAdmin{}
		dbUser := model.User{}
		c.Bind(&userAdmin)

		err := db.Where("email = ?", userAdmin.Email).First(&dbUser).Error

		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"status":  http.StatusNotFound,
				"message": "user not found",
			})
		}

		userPass := []byte(userAdmin.Password)
		dbPass := []byte(dbUser.Password)

		passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

		if passErr != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"status":  http.StatusUnauthorized,
				"message": "unauthorized",
			})
		}

		token, _ := middleware.CreateToken(userAdmin.IdUser, userAdmin.Fullname)

		userResponse := model.UserLoginResponse{
			IdUser:   dbUser.IdUser,
			Email:    dbUser.Email,
			Fullname: dbUser.Fullname,
			Token:    token,
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status":  http.StatusOK,
			"message": "success login",
			"data":    userResponse,
		})
	}
}

func RegisterAdmin(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		userAdmin := model.UserAdmin{}
		c.Bind(&userAdmin)

		// check user
		checkUser := model.UserAdmin{}
		db.Where("email = ?", userAdmin.Email).First(&checkUser)
		if checkUser.Fullname != "" {
			return c.JSON(http.StatusOK, echo.Map{
				"status":  http.StatusOK,
				"message": "This email is already registered",
			})
		}

		// encrypt password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userAdmin.Password), 10)
		userAdmin.Password = string(hashedPassword)

		// create user
		err = db.Create(&userAdmin).Error
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}

		userResponse := model.UserAdminRegisterResponse{
			IdUser:   userAdmin.IdUser,
			Email:    userAdmin.Email,
			Fullname: userAdmin.Fullname,
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status":  http.StatusOK,
			"message": "success register",
			"data":    userResponse,
		})

	}
}
