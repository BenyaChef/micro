package postusecasemock

import (
	"context"

	repositoryinterface "github.com/BenyaChef/micro/blog/boundary/repository"
	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
)

var _ repositoryinterface.PostRepository = (*PostRepositoryMock)(nil)

type PostRepositoryMock struct {
	InsertFunc  func(ctx context.Context, post *postentity.Post) error
	GetByIDFunc func(ctx context.Context, postID postentity.PostID) (*postentity.Post, error)
	ListFunc    func(ctx context.Context, limit, offset int) ([]*postentity.Post, error)
	UpdateFunc  func(ctx context.Context, post *postentity.Post) error
	DeleteFunc  func(ctx context.Context, postID postentity.PostID) error

	Inserted     *postentity.Post
	Updated      *postentity.Post
	Deleted      postentity.PostID
	CalledLimit  int
	CalledOffset int
}

func NewPostRepositoryMock() *PostRepositoryMock {
	return &PostRepositoryMock{
		InsertFunc:  func(context.Context, *postentity.Post) error { return nil },
		GetByIDFunc: func(context.Context, postentity.PostID) (*postentity.Post, error) { return nil, nil },
		ListFunc:    func(context.Context, int, int) ([]*postentity.Post, error) { return nil, nil },
		UpdateFunc:  func(context.Context, *postentity.Post) error { return nil },
		DeleteFunc:  func(context.Context, postentity.PostID) error { return nil },
	}
}

func (m *PostRepositoryMock) Insert(ctx context.Context, post *postentity.Post) error {
	m.Inserted = post

	return m.InsertFunc(ctx, post)
}

func (m *PostRepositoryMock) GetByID(ctx context.Context, postID postentity.PostID) (*postentity.Post, error) {
	return m.GetByIDFunc(ctx, postID)
}

func (m *PostRepositoryMock) List(ctx context.Context, limit, offset int) ([]*postentity.Post, error) {
	m.CalledLimit = limit
	m.CalledOffset = offset

	return m.ListFunc(ctx, limit, offset)
}

func (m *PostRepositoryMock) Update(ctx context.Context, post *postentity.Post) error {
	m.Updated = post

	return m.UpdateFunc(ctx, post)
}

func (m *PostRepositoryMock) Delete(ctx context.Context, postID postentity.PostID) error {
	m.Deleted = postID

	return m.DeleteFunc(ctx, postID)
}
