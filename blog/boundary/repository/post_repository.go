package repositoryinterface

import (
	"context"

	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
)

type PostRepository interface {
	Insert(ctx context.Context, post *postentity.Post) error
	GetByID(ctx context.Context, postID postentity.PostID) (*postentity.Post, error)
	List(ctx context.Context, limit, offset int) ([]*postentity.Post, error)
	Update(ctx context.Context, post *postentity.Post) error
	Delete(ctx context.Context, postID postentity.PostID) error
}
