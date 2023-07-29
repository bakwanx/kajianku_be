package controller

import (
	"kajianku_be/middleware"
	"kajianku_be/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		users := make([]model.User, 0)
		err := db.Find(&users).Error
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, echo.Map{
			"status":  http.StatusOK,
			"message": "success",
			"data":    users,
		})
	}
}

func Delete(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		var user model.User
		var mosque model.Mosque

		// delete user
		if err := db.Delete(&user, id).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}

		// delete mosque
		if err := db.Where("id_user = ?", id).Delete(&mosque, id).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  http.StatusOK,
			"message": "success delete",
			"data":    user,
		})
	}
}

func Register(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := model.User{}
		c.Bind(&user)

		// check user
		checkUser := model.User{}
		db.Where("email = ?", user.Email).First(&checkUser)
		if checkUser.Fullname != "" {
			return c.JSON(http.StatusOK, echo.Map{
				"status":  http.StatusOK,
				"message": "This email is already registered",
			})
		}

		// encrypt password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		user.Password = string(hashedPassword)

		// create user
		err = db.Create(&user).Error
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}

		userResponse := model.UserRegisterResponse{
			Id:       user.Id,
			Email:    user.Email,
			Fullname: user.Fullname,
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status":  http.StatusOK,
			"message": "success register",
			"data":    userResponse,
		})

	}
}

func Login(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := model.User{}
		dbUser := model.User{}
		c.Bind(&user)

		err := db.Where("email = ?", user.Email).First(&dbUser).Error

		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"status":  http.StatusNotFound,
				"message": "user not found",
			})
		}

		userPass := []byte(user.Password)
		dbPass := []byte(dbUser.Password)

		passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

		if passErr != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"status":  http.StatusUnauthorized,
				"message": "unauthorized",
			})
		}

		token, _ := middleware.CreateToken(user.Id, user.Fullname)

		userResponse := model.UserLoginResponse{
			Id:       user.Id,
			Email:    user.Email,
			Fullname: user.Fullname,
			Token:    token,
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status":  http.StatusOK,
			"message": "success login",
			"data":    userResponse,
		})
	}
}
