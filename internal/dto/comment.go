package dto

type CreateCommentRequestDTO struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
	NewsID  string `json:"news_id"`
}
