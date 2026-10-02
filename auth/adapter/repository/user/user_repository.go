package userrepository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	userrepomodel "github.com/BenyaChef/micro/auth/adapter/repository/user/model"
	repositoryinterface "github.com/BenyaChef/micro/auth/boundary/repository"
	userentity "github.com/BenyaChef/micro/auth/domain/entity/user"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
	postgresinterface "github.com/BenyaChef/micro/infrastructure/postgres/interface"
)

const uniqueViolationCode = "23505"

const (
	queryInsert = `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	queryGetByID = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1`

	queryGetByEmail = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1`
)

var _ repositoryinterface.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	executor       postgresinterface.Executor
	logPublisher   loggerinterface.LogPublisher
	errorProcessor *errorProcessor
}

func (r *UserRepository) Insert(ctx context.Context, user *userentity.User) error {
	if user == nil {
		return r.errorProcessor.LogAndReturn(ctx, ErrUserIsRequired)
	}

	model := userrepomodel.ToModel(user)

	_, err := r.executor.Exec(ctx, queryInsert,
		model.ID, model.Email, model.PasswordHash, model.CreatedAt, model.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return r.errorProcessor.LogAndReturn(ctx, ErrEmailConflict(model.Email, err))
		}

		return r.errorProcessor.LogAndReturn(ctx, ErrInsertFailed(model.ID, err))
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID commonuserentity.UserID) (*userentity.User, error) {
	return r.selectOne(ctx, queryGetByID, userID.String())
}

func (r *UserRepository) GetByEmail(ctx context.Context, email emailprimitive.Email) (*userentity.User, error) {
	return r.selectOne(ctx, queryGetByEmail, email.String())
}

func (r *UserRepository) selectOne(ctx context.Context, query, arg string) (*userentity.User, error) {
	rows, err := r.executor.Query(ctx, query, arg)
	if err != nil {
		return nil, r.errorProcessor.LogAndReturn(ctx, ErrSelectFailed(err))
	}

	model, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[userrepomodel.User])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, r.errorProcessor.LogAndReturn(ctx, ErrSelectFailed(err))
	}

	user, err := userrepomodel.ToEntity(model)
	if err != nil {
		return nil, r.errorProcessor.LogAndReturn(ctx, err)
	}

	return user, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError

	return errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode
}
