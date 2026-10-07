package initservice

import postusecase "github.com/BenyaChef/micro/blog/domain/usecase/post"

func (c *DependencyContainer) initUseCases() error {
	postUseCase, err := postusecase.NewBuilder().
		PostRepository(c.postRepository).
		LogPublisher(c.infra.LogPublisher()).
		Build()
	if err != nil {
		return err
	}

	c.postUseCase = postUseCase

	return nil
}
