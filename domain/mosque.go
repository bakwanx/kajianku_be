package domain

type Mosque struct {
	Id         int    `gorm:"primarykey" json:"id" form:"id"`
	MosqueName string `json:"mosque_name" form:"mosque_name"`
	Latitude   int    `json:"latitude" form:"latitude"`
	Longitude  int    `json:"longitude" form:"longitude"`
	Address    string `json:"address" form:"address"`
	Image      string `json:"image" form:"image"`
}
