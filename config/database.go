package config

import (
	"kajianku_be/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "root:@tcp(127.0.0.1:3306)/kajianku?charset=utf8mb4&parseTime=True&loc=Local"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		model.User{},
		model.Mosque{},
		model.Kajian{},
		model.UserAdmin{},
	)
}
