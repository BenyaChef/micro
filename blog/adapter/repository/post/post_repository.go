package postrepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	postrepomodel "github.com/BenyaChef/micro/blog/adapter/repository/post/model"
	repositoryinterface "github.com/BenyaChef/micro/blog/boundary/repository"
	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	postgresinterface "github.com/BenyaChef/micro/infrastructure/postgres/interface"
)

const (
	queryInsert = `
		INSERT INTO posts (id, author_id, title, body, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	queryGetByID = `
		SELECT id, author_id, title, body, created_at, updated_at
		FROM posts
		WHERE id = $1`

	queryList = `
		SELECT id, author_id, title, body, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	queryUpdate = `
		UPDATE posts
		SET title = $2, body = $3, updated_at = $4
		WHERE id = $1`

	queryDelete = `
		DELETE FROM posts
		WHERE id = $1`
)

var _ repositoryinterface.PostRepository = (*PostRepository)(nil)

type PostRepository struct {
	executor       postgresinterface.Executor
	logPublisher   loggerinterface.LogPublisher
	errorProcessor *errorProcessor
}

func (r *PostRepository) Insert(ctx context.Context, post *postentity.Post) error {
	if post == nil {
		return r.errorProcessor.LogAndReturn(ctx, ErrPostIsRequired)
	}

	model := postrepomodel.ToModel(post)

	_, err := r.executor.Exec(ctx, queryInsert,
		model.ID, model.AuthorID, model.Title, model.Body, model.CreatedAt, model.UpdatedAt,
	)
	if err != nil {
		return r.errorProcessor.LogAndReturn(ctx, ErrInsertFailed(model.ID, err))
	}

	return nil
}

func (r *PostRepository) GetByID(ctx context.Context, postID postentity.PostID) (*postentity.Post, error) {
	rows, err := r.executor.Query(ctx, queryGetByID, postID.String())
	if err != nil {
		return nil, r.errorProcessor.LogAndReturn(ctx, ErrSelectFailed(err))
	}

	model, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[postrepomodel.Post])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, r.errorProcessor.LogAndReturn(ctx, ErrSelectFailed(err))
	}

	post, err := postrepomodel.ToEntity(model)
	if err != nil {
		return nil, r.errorProcessor.LogAndReturn(ctx, err)
	}

	return post, nil
}

func (r *PostRepository) List(ctx context.Context, limit, offset int) ([]*postentity.Post, error) {
	rows, err := r.executor.Query(ctx, queryList, limit, offset)
	if err != nil {
		return nil, r.errorProcessor.LogAndReturn(ctx, ErrSelectFailed(err))
	}

	models, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[postrepomodel.Post])
	if err != nil {
		return nil, r.errorProcessor.LogAndReturn(ctx, ErrSelectFailed(err))
	}

	posts, err := postrepomodel.ToEntities(models)
	if err != nil {
		return nil, r.errorProcessor.LogAndReturn(ctx, err)
	}

	return posts, nil
}

func (r *PostRepository) Update(ctx context.Context, post *postentity.Post) error {
	if post == nil {
		return r.errorProcessor.LogAndReturn(ctx, ErrPostIsRequired)
	}

	model := postrepomodel.ToModel(post)

	tag, err := r.executor.Exec(ctx, queryUpdate, model.ID, model.Title, model.Body, model.UpdatedAt)
	if err != nil {
		return r.errorProcessor.LogAndReturn(ctx, ErrUpdateFailed(model.ID, err))
	}

	if tag.RowsAffected() == 0 {
		return r.errorProcessor.LogAndReturn(ctx, ErrPostVanished(model.ID))
	}

	return nil
}

func (r *PostRepository) Delete(ctx context.Context, postID postentity.PostID) error {
	tag, err := r.executor.Exec(ctx, queryDelete, postID.String())
	if err != nil {
		return r.errorProcessor.LogAndReturn(ctx, ErrDeleteFailed(postID.String(), err))
	}

	if tag.RowsAffected() == 0 {
		return r.errorProcessor.LogAndReturn(ctx, ErrPostVanished(postID.String()))
	}

	return nil
}
