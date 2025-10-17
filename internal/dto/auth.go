package dto

type LoginRequestDTO struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type RefreshRequestDTO struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
