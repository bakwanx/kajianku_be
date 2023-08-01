package config

import (
	"kajianku_be/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "golang_db:golang123@tcp(golang-db.cbuoaypgqh0v.us-east-1.rds.amazonaws.com:3306)/golang_db?charset=utf8mb4&parseTime=True&loc=Local"
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
