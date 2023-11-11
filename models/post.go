package models

import "time"

type Post struct {
	ID          string    `json:"id"`
	PostContent string    `json:"post_content"`
	UserId      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"createdAt"`
}
