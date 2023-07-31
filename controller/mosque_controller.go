package controller

import (
	"context"
	"kajianku_be/config"
	"kajianku_be/model"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

func RegisterMosque(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		idUser, _ := strconv.Atoi(c.FormValue("id_user"))
		mosqueName := c.FormValue("mosque_name")
		latitude := c.FormValue("latitude")
		longitude := c.FormValue("longitude")
		address := c.FormValue("address")
		file, err := c.FormFile("image")

		checkMosque := model.Mosque{}
		db.Where("id_user  = ?", idUser).First(&checkMosque)
		if checkMosque.MosqueName != "" {
			return c.JSON(http.StatusOK, echo.Map{
				"status":  http.StatusOK,
				"message": "This mosque is already registered",
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

		mosque := model.Mosque{
			IdUser:     idUser,
			MosqueName: mosqueName,
			Latitude:   latitude,
			Longitude:  longitude,
			Address:    address,
			Image:      result.Location,
		}

		err = db.Create(&mosque).Error
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

func GetAllMosque(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		mosque := make([]model.Mosque, 0)
		err := db.Find(&mosque).Error
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

func UpdateMosque(db DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		idUser := c.Param("id_user")
		mosqueName := c.FormValue("mosque_name")
		latitude := c.FormValue("latitude")
		longitude := c.FormValue("longitude")
		address := c.FormValue("address")
		file, err := c.FormFile("image")
		var resultUrl string

		// check mosque
		checkMosque := model.Mosque{}
		db.Where("id_user = ?", idUser).First(&checkMosque)
		if checkMosque.MosqueName == "" {
			return c.JSON(http.StatusNotFound, echo.Map{
				"status":  http.StatusNotFound,
				"message": "Mosque not found",
			})
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

		var mosque model.Mosque = model.Mosque{
			MosqueName: mosqueName,
			Latitude:   latitude,
			Longitude:  longitude,
			Address:    address,
			Image:      resultUrl,
		}

		if err := db.Where("id_user", idUser).Updates(&mosque).Error; err != nil {
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
