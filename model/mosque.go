package model

type Mosque struct {
	IdMosque   int    `gorm:"primarykey" json:"id_mosque" form:"id_mosque"`
	IdUser     int    `json:"id_user" form:"id_user"`
	MosqueName string `json:"mosque_name" form:"mosque_name"`
	Latitude   string `json:"latitude" form:"latitude"`
	Longitude  string `json:"longitude" form:"longitude"`
	Address    string `json:"address" form:"address"`
	Image      string `json:"image" form:"image"`
}
