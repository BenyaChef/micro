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

	const query = `INSERT INTO users (` + userrepomodel.Columns + `) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.executor.Exec(ctx, query, model.InsertValues()...)
	if err != nil {
		if isUniqueViolation(err) {
			return r.errorProcessor.LogAndReturn(ctx, ErrEmailConflict(model.Email, err))
		}

		return r.errorProcessor.LogAndReturn(ctx, ErrInsertFailed(model.ID, err))
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID commonuserentity.UserID) (*userentity.User, error) {
	return r.selectOne(ctx, "SELECT "+userrepomodel.Columns+" FROM users WHERE id = $1", userID.String())
}

func (r *UserRepository) GetByEmail(ctx context.Context, email emailprimitive.Email) (*userentity.User, error) {
	return r.selectOne(ctx, "SELECT "+userrepomodel.Columns+" FROM users WHERE email = $1", email.String())
}

func (r *UserRepository) selectOne(ctx context.Context, query string, arg any) (*userentity.User, error) {
	model := &userrepomodel.User{}

	err := r.executor.QueryRow(ctx, query, arg).Scan(model.ScanTargets()...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

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
