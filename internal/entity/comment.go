package entity

import "time"

type Comment struct {
	ID        string    `json:"id"`
	NewsID    string    `json:"news_id"`
	Name      string    `json:"name"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
