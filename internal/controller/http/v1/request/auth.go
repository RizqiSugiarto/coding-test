package request

type Auth struct {
	Username string `json:"username" binding:"required" example:"Naruto"`
	Password string `json:"password" binding:"required" example:"Uuk2019Tyu"`
}
