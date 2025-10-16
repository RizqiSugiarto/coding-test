package request

// CustomPage represents the request body for creating custom page.
type CustomPage struct {
	CustomURL string `json:"custom_url" binding:"required" example:"/about-us"`
	Content   string `json:"content" binding:"required" example:"<h1>About Us</h1><p>This is our about page...</p>"`
}

// UpdateCustomPage represents the request body for updating custom page.
type UpdateCustomPage struct {
	CustomURL string `json:"custom_url" binding:"required" example:"/about-company"`
	Content   string `json:"content" binding:"required" example:"<h1>About Company</h1><p>Updated content...</p>"`
}
