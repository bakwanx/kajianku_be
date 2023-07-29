package model

type User struct {
	Id       int    `gorm:"primarykey" json:"id" form:"id"`
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Status   int    `json:"status" form:"status"`
}

type UserLoginResponse struct {
	Id       int    `json:"id" form:"id"`
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email"`
	Token    string `json:"token" form:"token"`
}

type UserRegisterResponse struct {
	Id       int    `json:"id" form:"id"`
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email"`
}

// type UserRepository interface {
// 	GetAll(dest interface{}, conds ...interface{}) *gorm.DB
// 	GetByStatus(query interface{}, args ...interface{}) *gorm.DB
// 	GetById(query interface{}, args ...interface{}) *gorm.DB
// 	Create(value interface{}) *gorm.DB
// 	Delete(value interface{}, conds ...interface{}) *gorm.DB
// }
