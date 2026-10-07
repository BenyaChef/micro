package usecaseinterface

import (
	"context"

	boundarymodel "github.com/BenyaChef/micro/blog/boundary/model"
	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
)

type PostUseCase interface {
	Create(ctx context.Context, input *boundarymodel.CreatePostInput) (*postentity.Post, error)
	GetByID(ctx context.Context, postID postentity.PostID) (*postentity.Post, error)
	List(ctx context.Context, input *boundarymodel.ListPostsInput) ([]*postentity.Post, error)
	Update(ctx context.Context, input *boundarymodel.UpdatePostInput) (*postentity.Post, error)
	Delete(ctx context.Context, postID postentity.PostID, authorID commonuserentity.UserID) error
}
