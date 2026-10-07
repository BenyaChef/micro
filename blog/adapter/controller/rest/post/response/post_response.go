package responserestmodel

import (
	"time"

	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
)

type PostResponse struct {
	PostID    string `json:"post_id"`
	AuthorID  string `json:"author_id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func FromPost(post *postentity.Post) *PostResponse {
	return &PostResponse{
		PostID:    post.ID().String(),
		AuthorID:  post.AuthorID().String(),
		Title:     post.Title().String(),
		Body:      post.Body().String(),
		CreatedAt: post.CreatedAt().Time().Format(time.RFC3339Nano),
		UpdatedAt: post.UpdatedAt().Time().Format(time.RFC3339Nano),
	}
}

func FromPosts(posts []*postentity.Post) []*PostResponse {
	responses := make([]*PostResponse, 0, len(posts))
	for _, post := range posts {
		responses = append(responses, FromPost(post))
	}

	return responses
}
