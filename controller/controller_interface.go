package controller

import "gorm.io/gorm"

type DB interface {
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Raw(sql string, values ...interface{}) *gorm.DB
}
