package request

// Comment represents the request body for creating a comment.
type Comment struct {
	Name    string `json:"name" binding:"required" example:"John Doe"`
	Comment string `json:"comment" binding:"required" example:"This is a great article!"`
}
