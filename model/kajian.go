package model

type Kajian struct {
	IdKajian    int    `gorm:"primarykey" json:"id_kajian" form:"id_kajian"`
	IdUser      int    `json:"id_user" form:"id_user"`
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Date        string `json:"date" form:"date"`
	Poster      string `json:"poster" form:"poster"`
}

type KajianByDistanceResponse struct {
	IdKajian    int    `json:"id_kajian" form:"id_kajian"`
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Date        string `json:"date" form:"date"`
	Poster      string `json:"poster" form:"poster"`
	IdMosque    int    `json:"id_mosque" form:"id_mosque"`
	MosqueName  string `json:"mosque_name" form:"mosque_name"`
	Latitude    string `json:"latitude" form:"latitude"`
	Longitude   string `json:"longitude" form:"longitude"`
	Address     string `json:"address" form:"address"`
	Image       string `json:"image" form:"image"`
	Distance    string `form:"distance"`
}
