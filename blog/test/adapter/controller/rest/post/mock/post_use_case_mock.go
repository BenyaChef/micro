package postrestcontrollermock

import (
	"context"

	boundarymodel "github.com/BenyaChef/micro/blog/boundary/model"
	usecaseinterface "github.com/BenyaChef/micro/blog/boundary/usecase"
	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
)

var _ usecaseinterface.PostUseCase = (*PostUseCaseMock)(nil)

type PostUseCaseMock struct {
	CreateFunc  func(ctx context.Context, input *boundarymodel.CreatePostInput) (*postentity.Post, error)
	GetByIDFunc func(ctx context.Context, postID postentity.PostID) (*postentity.Post, error)
	ListFunc    func(ctx context.Context, input *boundarymodel.ListPostsInput) ([]*postentity.Post, error)
	UpdateFunc  func(ctx context.Context, input *boundarymodel.UpdatePostInput) (*postentity.Post, error)
	DeleteFunc  func(ctx context.Context, postID postentity.PostID, authorID commonuserentity.UserID) error

	CreateInput *boundarymodel.CreatePostInput
	ListInput   *boundarymodel.ListPostsInput
}

func NewPostUseCaseMock() *PostUseCaseMock {
	return &PostUseCaseMock{
		CreateFunc: func(context.Context, *boundarymodel.CreatePostInput) (*postentity.Post, error) {
			return nil, nil
		},
		GetByIDFunc: func(context.Context, postentity.PostID) (*postentity.Post, error) { return nil, nil },
		ListFunc: func(context.Context, *boundarymodel.ListPostsInput) ([]*postentity.Post, error) {
			return nil, nil
		},
		UpdateFunc: func(context.Context, *boundarymodel.UpdatePostInput) (*postentity.Post, error) {
			return nil, nil
		},
		DeleteFunc: func(context.Context, postentity.PostID, commonuserentity.UserID) error { return nil },
	}
}

func (m *PostUseCaseMock) Create(ctx context.Context, input *boundarymodel.CreatePostInput) (*postentity.Post, error) {
	m.CreateInput = input

	return m.CreateFunc(ctx, input)
}

func (m *PostUseCaseMock) GetByID(ctx context.Context, postID postentity.PostID) (*postentity.Post, error) {
	return m.GetByIDFunc(ctx, postID)
}

func (m *PostUseCaseMock) List(ctx context.Context, input *boundarymodel.ListPostsInput) ([]*postentity.Post, error) {
	m.ListInput = input

	return m.ListFunc(ctx, input)
}

func (m *PostUseCaseMock) Update(ctx context.Context, input *boundarymodel.UpdatePostInput) (*postentity.Post, error) {
	return m.UpdateFunc(ctx, input)
}

func (m *PostUseCaseMock) Delete(ctx context.Context, postID postentity.PostID, authorID commonuserentity.UserID) error {
	return m.DeleteFunc(ctx, postID, authorID)
}
