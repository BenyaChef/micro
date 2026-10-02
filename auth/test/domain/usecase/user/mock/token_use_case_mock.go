package userusecasemock

import (
	"context"

	boundarymodel "github.com/BenyaChef/micro/auth/boundary/model"
	usecaseinterface "github.com/BenyaChef/micro/auth/boundary/usecase"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
)

var _ usecaseinterface.TokenUseCase = (*TokenUseCaseMock)(nil)

const (
	StubAccessToken = "stub-access-token"
	StubExpiresIn   = 3600
)

type TokenUseCaseMock struct {
	IssueFunc    func(ctx context.Context, userID commonuserentity.UserID) (*boundarymodel.LoginResult, error)
	ValidateFunc func(ctx context.Context, rawToken string) (commonuserentity.UserID, error)
}

func NewTokenUseCaseMock() *TokenUseCaseMock {
	return &TokenUseCaseMock{
		IssueFunc: func(context.Context, commonuserentity.UserID) (*boundarymodel.LoginResult, error) {
			return &boundarymodel.LoginResult{AccessToken: StubAccessToken, ExpiresIn: StubExpiresIn}, nil
		},
		ValidateFunc: func(context.Context, string) (commonuserentity.UserID, error) {
			return commonuserentity.NewUserID(), nil
		},
	}
}

func (m *TokenUseCaseMock) Issue(ctx context.Context, userID commonuserentity.UserID) (*boundarymodel.LoginResult, error) {
	return m.IssueFunc(ctx, userID)
}

func (m *TokenUseCaseMock) Validate(ctx context.Context, rawToken string) (commonuserentity.UserID, error) {
	return m.ValidateFunc(ctx, rawToken)
}
