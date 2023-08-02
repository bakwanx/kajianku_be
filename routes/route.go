package routes

import (
	"kajianku_be/config"
	"os"

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
	e.GET("/kajian/:distance/:latitude/:longitude", controller.GetKajianByDistance(db))

	eJwt := e.Group("")
	eJwt.Use(mid.JWT([]byte(os.Getenv("SECRET_JWT"))))
	// Users
	eJwt.GET("/users", controller.GetAllUsers(db))
	eJwt.PATCH("/users/:id_user", controller.UpdateUser(db))

	// Mosque
	eJwt.GET("/mosque", controller.GetAllMosque(db))
	eJwt.POST("/mosque", controller.RegisterMosque(db))
	eJwt.GET("/mosque/:id_user", controller.GetMosqueByUserId(db))
	eJwt.PATCH("/mosque/:id_user", controller.UpdateMosque(db))

	// Kajian
	eJwt.POST("/kajian", controller.PostKajian(db))
	eJwt.DELETE("/kajian/:id_kajian", controller.DeleteKajian(db))
	eJwt.PATCH("/kajian/:id_kajian", controller.UpdateKajian(db))

	// Admin
	admin := e.Group("admin")
	admin.POST("/register", controller.RegisterAdmin(db))
	admin.POST("/login", controller.LoginAdmin(db))
	admin.PATCH("/approve", controller.ApproveUser(db))
	admin.PATCH("/:id_user", controller.UpdateUserAdmin(db))
	admin.GET("/users/:status", controller.GetAllUsersByStatus(db))
	return e
}
