package initservice

import (
	tokenusecase "github.com/BenyaChef/micro/auth/domain/usecase/token"
	userusecase "github.com/BenyaChef/micro/auth/domain/usecase/user"
	"github.com/BenyaChef/micro/infrastructure/jwt"
	"github.com/BenyaChef/micro/infrastructure/passwordhasher"
)

func (c *DependencyContainer) initServices() error {
	tokenService, err := jwt.NewBuilder().
		Secret(c.environments.jwtSecret).
		Issuer(c.serviceName).
		TTL(c.environments.jwtTTL).
		Build()
	if err != nil {
		return err
	}

	c.tokenService = tokenService

	hasherBuilder := passwordhasher.NewBuilder()
	if c.environments.passwordHashCost > 0 {
		hasherBuilder = hasherBuilder.Cost(c.environments.passwordHashCost)
	}

	passwordHasher, err := hasherBuilder.Build()
	if err != nil {
		return err
	}

	c.passwordHasher = passwordHasher

	return nil
}

func (c *DependencyContainer) initUseCases() error {
	tokenUseCase, err := tokenusecase.NewBuilder().
		TokenService(c.tokenService).
		LogPublisher(c.infra.LogPublisher()).
		Build()
	if err != nil {
		return err
	}

	c.tokenUseCase = tokenUseCase

	userUseCase, err := userusecase.NewBuilder().
		UserRepository(c.userRepository).
		TokenUseCase(tokenUseCase).
		PasswordHasher(c.passwordHasher).
		LogPublisher(c.infra.LogPublisher()).
		Build()
	if err != nil {
		return err
	}

	c.userUseCase = userUseCase

	return nil
}
