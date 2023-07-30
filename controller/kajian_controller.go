package controller

import (
	"kajianku_be/model"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func PostKajian(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := model.User{}
		kajian := model.Kajian{}
		c.Bind(&kajian)

		// check user status
		err := db.Where("id_kajian  = ? AND status = ?", kajian.IdUser, 1).First(&user).Error
		if user.Email == "" {
			return c.JSON(http.StatusNotFound, echo.Map{
				"status":  http.StatusNotFound,
				"message": "user not found",
			})
		}
		if err != nil {
			return c.JSON(http.StatusNotAcceptable, echo.Map{
				"status":  http.StatusNotAcceptable,
				"message": "user unverified",
			})
		}

		// parse format time
		_, err = time.Parse("02/01/2006", kajian.Date)
		if err != nil {
			return c.JSON(http.StatusNotAcceptable, echo.Map{
				"status":  http.StatusNotAcceptable,
				"message": "accepted format 'dd/MM/yyyy'",
			})
		}

		// create kajian
		err = db.Create(&kajian).Error
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"status":  http.StatusOK,
			"message": "success upload",
			"data":    kajian,
		})

	}
}

func GetKajianByDistance(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		distance := c.Param("distance")
		latitude := c.Param("latitude")
		longitude := c.Param("longitude")
		kajianResponse := []model.KajianByDistanceResponse{}
		if distance == "" {
			distance = "5"
		}

		sqlRow := `SELECT 
			kajians.id_kajian,
			kajians.title,
			kajians.description,
			kajians.date,
			kajians.poster,
			mosques.id_mosque,
			mosques.mosque_name,
			mosques.image,
			mosques.latitude,
			mosques.longitude, 
			( 6371 * acos ( 
				cos ( radians( ? ) ) 
				* cos( radians( latitude ) ) 
				* cos( radians( longitude ) - radians( ? ) ) 
				+ sin ( radians( ? ) ) 
				* sin( radians( latitude ) ) 
				) 
			) AS distance FROM mosques JOIN kajians ON mosques.id_user = kajians.id_user HAVING distance <= ?;`

		db.Raw(sqlRow, latitude, longitude, latitude, distance).Scan(&kajianResponse)

		return c.JSON(http.StatusOK, echo.Map{
			"status":  http.StatusOK,
			"message": "success",
			"data":    kajianResponse,
		})
	}
}
