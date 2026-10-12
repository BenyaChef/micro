package postusecase

import (
	"context"

	boundarymodel "github.com/BenyaChef/micro/blog/boundary/model"
	repositoryinterface "github.com/BenyaChef/micro/blog/boundary/repository"
	usecaseinterface "github.com/BenyaChef/micro/blog/boundary/usecase"
	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	timeprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/time"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

const (
	DefaultListLimit = 20
	MaxListLimit     = 100
)

var _ usecaseinterface.PostUseCase = (*PostUseCase)(nil)

type PostUseCase struct {
	postRepository repositoryinterface.PostRepository
	logPublisher   loggerinterface.LogPublisher
	errorProcessor *errorProcessor
}

func (uc *PostUseCase) Create(ctx context.Context, input *boundarymodel.CreatePostInput) (*postentity.Post, error) {
	if input == nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrInputIsRequired)
	}

	if input.AuthorID.IsZero() {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrAuthorIDIsRequired)
	}

	title, body, err := uc.parseContent(ctx, input.Title, input.Body)
	if err != nil {
		return nil, err
	}

	post, err := postentity.NewBuilder().
		ID(postentity.NewPostID()).
		AuthorID(input.AuthorID).
		Title(title).
		Body(body).
		Build()
	if err != nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, err)
	}

	if err := uc.postRepository.Insert(ctx, post); err != nil {
		return nil, uc.errorProcessor.WrapOperationError(err)
	}

	return post, nil
}

func (uc *PostUseCase) GetByID(ctx context.Context, postID postentity.PostID) (*postentity.Post, error) {
	if postID.IsZero() {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrPostIDIsRequired)
	}

	post, err := uc.postRepository.GetByID(ctx, postID)
	if err != nil {
		return nil, uc.errorProcessor.WrapOperationError(err)
	}

	if post == nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrPostNotFound(postID.String()))
	}

	return post, nil
}

func (uc *PostUseCase) List(ctx context.Context, input *boundarymodel.ListPostsInput) ([]*postentity.Post, error) {
	if input == nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrInputIsRequired)
	}

	limit, offset := normalizePaging(input.Limit, input.Offset)

	posts, err := uc.postRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, uc.errorProcessor.WrapOperationError(err)
	}

	return posts, nil
}

func (uc *PostUseCase) Update(ctx context.Context, input *boundarymodel.UpdatePostInput) (*postentity.Post, error) {
	if input == nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrInputIsRequired)
	}

	existing, err := uc.ownedPost(ctx, input.PostID, input.AuthorID)
	if err != nil {
		return nil, err
	}

	title, body, err := uc.parseContent(ctx, input.Title, input.Body)
	if err != nil {
		return nil, err
	}

	updated, err := postentity.NewBuilder().
		ID(existing.ID()).
		AuthorID(existing.AuthorID()).
		Title(title).
		Body(body).
		CreatedAt(existing.CreatedAt()).
		UpdatedAt(timeprimitive.Now()).
		Build()
	if err != nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, err)
	}

	if err := uc.postRepository.Update(ctx, updated); err != nil {
		return nil, uc.errorProcessor.WrapOperationError(err)
	}

	return updated, nil
}

func (uc *PostUseCase) Delete(ctx context.Context, postID postentity.PostID, authorID commonuserentity.UserID) error {
	if _, err := uc.ownedPost(ctx, postID, authorID); err != nil {
		return err
	}

	if err := uc.postRepository.Delete(ctx, postID); err != nil {
		return uc.errorProcessor.WrapOperationError(err)
	}

	return nil
}

func (uc *PostUseCase) ownedPost(
	ctx context.Context,
	postID postentity.PostID,
	authorID commonuserentity.UserID,
) (*postentity.Post, error) {
	if authorID.IsZero() {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrAuthorIDIsRequired)
	}

	post, err := uc.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if !post.IsAuthoredBy(authorID) {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrNotAuthor(postID.String(), authorID.String()))
	}

	return post, nil
}

func (uc *PostUseCase) parseContent(ctx context.Context, rawTitle, rawBody string) (postentity.Title, postentity.Body, error) {
	title, err := postentity.TitleFrom(rawTitle)
	if err != nil {
		return "", "", uc.errorProcessor.LogAndReturn(ctx, err)
	}

	body, err := postentity.BodyFrom(rawBody)
	if err != nil {
		return "", "", uc.errorProcessor.LogAndReturn(ctx, err)
	}

	return title, body, nil
}

func normalizePaging(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = DefaultListLimit
	}

	if limit > MaxListLimit {
		limit = MaxListLimit
	}

	if offset < 0 {
		offset = 0
	}

	return limit, offset
}
