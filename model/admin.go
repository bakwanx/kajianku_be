package model

type UserAdmin struct {
	IdUser   int    `gorm:"primarykey" json:"id_user" form:"id_user"`
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type UserAdminRegisterResponse struct {
	IdUser   int    `json:"id_user" form:"id_user"`
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email"`
}
