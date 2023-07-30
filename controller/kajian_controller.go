package controller

import (
	"context"
	"kajianku_be/config"
	"kajianku_be/model"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

func PostKajian(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := model.User{}

		IdUser, _ := strconv.Atoi(c.FormValue("id_user"))
		title := c.FormValue("title")
		description := c.FormValue("description")
		date := c.FormValue("date")
		file, err := c.FormFile("poster")

		// check user status
		err = db.Where("id_user  = ? AND status = ?", IdUser, 1).First(&user).Error
		if err != nil {
			return c.JSON(http.StatusNotAcceptable, echo.Map{
				"status":  http.StatusNotAcceptable,
				"message": "user unverified",
			})
		}

		// check format time
		_, err = time.Parse("02/01/2006", date)
		if err != nil {
			return c.JSON(http.StatusNotAcceptable, echo.Map{
				"status":  http.StatusNotAcceptable,
				"message": "accepted format 'dd/MM/yyyy'",
			})
		}

		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}
		defer src.Close()

		result, err := config.Uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("kajianku-bucket"),
			Key:    aws.String(file.Filename),
			Body:   src,
			ACL:    "public-read",
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}

		kajian := model.Kajian{
			IdUser:      IdUser,
			Title:       title,
			Description: description,
			Date:        date,
			Poster:      result.Location,
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
