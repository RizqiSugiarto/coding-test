package request

// News represents the request body for creating news.
type News struct {
	CategoryID string `json:"category_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title      string `json:"title" binding:"required" example:"Breaking News: Technology Advances"`
	Content    string `json:"content" binding:"required" example:"This is the full content of the news article..."`
}

// UpdateNews represents the request body for updating news.
type UpdateNews struct {
	CategoryID string `json:"category_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title      string `json:"title" binding:"required" example:"Updated News Title"`
	Content    string `json:"content" binding:"required" example:"This is the updated content..."`
}
