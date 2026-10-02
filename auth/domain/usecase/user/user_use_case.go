package userusecase

import (
	"context"

	boundarymodel "github.com/BenyaChef/micro/auth/boundary/model"
	repositoryinterface "github.com/BenyaChef/micro/auth/boundary/repository"
	serviceinterface "github.com/BenyaChef/micro/auth/boundary/service"
	usecaseinterface "github.com/BenyaChef/micro/auth/boundary/usecase"
	userentity "github.com/BenyaChef/micro/auth/domain/entity/user"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	emailprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/email"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

const MinPasswordLength = 8

var _ usecaseinterface.UserUseCase = (*UserUseCase)(nil)

type UserUseCase struct {
	userRepository repositoryinterface.UserRepository
	tokenUseCase   usecaseinterface.TokenUseCase
	passwordHasher serviceinterface.PasswordHasher
	logPublisher   loggerinterface.LogPublisher
	errorProcessor *errorProcessor
}

func (uc *UserUseCase) Register(ctx context.Context, data *boundarymodel.RegisterData) (*userentity.User, error) {
	if data == nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrDataIsRequired)
	}

	email, err := uc.parseCredentials(ctx, data.Email, data.Password)
	if err != nil {
		return nil, err
	}

	existing, err := uc.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, uc.errorProcessor.LogAndReturnOperationError(ctx, err)
	}

	if existing != nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrEmailAlreadyTaken(email.String()))
	}

	user, err := uc.buildUser(ctx, email, data.Password)
	if err != nil {
		return nil, err
	}

	if err := uc.userRepository.Insert(ctx, user); err != nil {
		return nil, uc.errorProcessor.LogAndReturnOperationError(ctx, err)
	}

	return user, nil
}

func (uc *UserUseCase) Login(ctx context.Context, data *boundarymodel.LoginData) (*boundarymodel.LoginResult, error) {
	if data == nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrDataIsRequired)
	}

	email, err := emailprimitive.EmailFrom(data.Email)
	if err != nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrInvalidCredentials)
	}

	user, err := uc.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, uc.errorProcessor.LogAndReturnOperationError(ctx, err)
	}

	if user == nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrInvalidCredentials)
	}

	if err := uc.passwordHasher.Compare(user.PasswordHash().String(), data.Password); err != nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrInvalidCredentials)
	}

	return uc.tokenUseCase.Issue(ctx, user.ID())
}

func (uc *UserUseCase) GetByID(ctx context.Context, userID commonuserentity.UserID) (*userentity.User, error) {
	if userID.IsZero() {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrDataIsRequired)
	}

	user, err := uc.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, uc.errorProcessor.LogAndReturnOperationError(ctx, err)
	}

	if user == nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, ErrUserNotFound(userID.String()))
	}

	return user, nil
}

func (uc *UserUseCase) parseCredentials(ctx context.Context, rawEmail, rawPassword string) (emailprimitive.Email, error) {
	email, err := emailprimitive.EmailFrom(rawEmail)
	if err != nil {
		return "", uc.errorProcessor.LogAndReturn(ctx, err)
	}

	if rawPassword == "" {
		return "", uc.errorProcessor.LogAndReturn(ctx, ErrPasswordIsRequired)
	}

	if len(rawPassword) < MinPasswordLength {
		return "", uc.errorProcessor.LogAndReturn(ctx, ErrPasswordTooShort(len(rawPassword), MinPasswordLength))
	}

	return email, nil
}

func (uc *UserUseCase) buildUser(ctx context.Context, email emailprimitive.Email, rawPassword string) (*userentity.User, error) {
	rawHash, err := uc.passwordHasher.Hash(rawPassword)
	if err != nil {
		return nil, uc.errorProcessor.LogAndReturnOperationError(ctx, err)
	}

	passwordHash, err := userentity.PasswordHashFrom(rawHash)
	if err != nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, err)
	}

	user, err := userentity.NewBuilder().
		ID(commonuserentity.NewUserID()).
		Email(email).
		PasswordHash(passwordHash).
		Build()
	if err != nil {
		return nil, uc.errorProcessor.LogAndReturn(ctx, err)
	}

	return user, nil
}
