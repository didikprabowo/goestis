package models

import (
	"time"
)

// Article
type Article struct {
	ID      int       `json:"id"`
	Author  string    `json:"author"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

// NewArticle
func NewArticle() *Article {
	return &Article{
		Created: time.Now(),
	}
}

type ArticleElastic struct {
	ID      int64     `json:"id,omitempty"`
	Author  string    `json:"author,omitempty"`
	Body    string    `json:"body,omitempty"`
	Created time.Time `json:"created,omitempty"`
}
