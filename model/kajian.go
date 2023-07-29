package model

import "time"

type Kajian struct {
	Id          int       `gorm:"primarykey" json:"id" form:"id"`
	IdUser      int       `json:"id_user" form:"id_user"`
	Title       string    `json:"title" form:"title"`
	Description int       `json:"description" form:"description"`
	Date        time.Time `json:"date" form:"date"`
	Poster      string    `json:"poster" form:"poster"`
}
