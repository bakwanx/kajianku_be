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

		idUser, _ := strconv.Atoi(c.FormValue("id_user"))
		title := c.FormValue("title")
		description := c.FormValue("description")
		date := c.FormValue("date")
		file, err := c.FormFile("poster")

		// check user status
		err = db.Where("id_user  = ? AND status = ?", idUser, 1).First(&user).Error
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
			IdUser:      idUser,
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

func DeleteKajian(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id_kajian := c.Param("id_kajian")
		var kajian model.Kajian

		// check kajian
		checkKajian := model.Kajian{}
		db.Where("id_kajian = ?", id_kajian).First(&checkKajian)
		if checkKajian.Title == "" {
			return c.JSON(http.StatusNotFound, echo.Map{
				"status":  http.StatusNotFound,
				"message": "Kajian not found",
			})
		}

		// delete kajian
		if err := db.Delete(&kajian, id_kajian).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  http.StatusOK,
			"message": "success delete",
			"data":    checkKajian,
		})
	}
}

func UpdateKajian(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		idKajian := c.Param("id_kajian")
		title := c.FormValue("title")
		description := c.FormValue("description")
		date := c.FormValue("date")
		file, err := c.FormFile("poster")
		var resultUrl string

		// check kajian
		checkKajian := model.Kajian{}
		db.Where("id_kajian = ?", idKajian).First(&checkKajian)
		if checkKajian.Title == "" {
			return c.JSON(http.StatusNotFound, echo.Map{
				"status":  http.StatusNotFound,
				"message": "Kajian not found",
			})
		}

		// check format time
		if date != "" {
			_, err := time.Parse("02/01/2006", date)
			if err != nil {
				return c.JSON(http.StatusNotAcceptable, echo.Map{
					"status":  http.StatusNotAcceptable,
					"message": "accepted format 'dd/MM/yyyy'",
				})
			}
		}

		if file != nil && err == nil {
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

			resultUrl = result.Location

		}

		var kajian model.Kajian = model.Kajian{
			Title:       title,
			Description: description,
			Date:        date,
			Poster:      resultUrl,
		}

		if err := db.Where("id_kajian", idKajian).Updates(&kajian).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  http.StatusOK,
			"message": "success update",
		})
	}

}
