package domain

type User struct {
	Id       int    `gorm:"primarykey" json:"id" form:"id"`
	IdMosque int    `json:"id_mosque" form:"id_mosque"`
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Status   int    `json:"status" form:"status"`
}
