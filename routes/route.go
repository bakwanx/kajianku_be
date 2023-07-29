package routes

import (
	"kajianku_be/config"

	"kajianku_be/controller"
	m "kajianku_be/middleware"

	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}

	err = config.MigrateDB(db)
	if err != nil {
		panic(err)
	}
	e := echo.New()
	m.LogMiddleware(e)
	e.POST("/register", controller.Register(db))
	e.POST("/login", controller.Login(db))

	eJwt := e.Group("")
	eJwt.Use(mid.JWT([]byte(config.SECRET_JWT)))
	eJwt.GET("/users", controller.GetAllUsers(db))
	eJwt.POST("/mosque", controller.RegisterMosque(db))
	eJwt.GET("/mosque/:id_user", controller.GetMosqueByUserId(db))
	return e
}
