package request

// Category represents the request body for creating a category.
type Category struct {
	Name string `json:"name" binding:"required" example:"Technology"`
}

// UpdateCategory represents the request body for updating a category.
type UpdateCategory struct {
	Name string `json:"name" binding:"required" example:"Updated Technology"`
}
