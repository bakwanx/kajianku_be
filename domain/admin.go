package domain

type UserAdmin struct {
	Id       int    `gorm:"primarykey" json:"id" form:"id"`
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
