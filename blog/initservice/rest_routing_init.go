package initservice

import (
	"github.com/go-chi/chi/v5"

	restservermiddleware "github.com/BenyaChef/micro/infrastructure/restserver/middleware"
)

func (c *DependencyContainer) initRESTRouting() error {
	authMiddleware := restservermiddleware.Auth(c.tokenValidator, c.restServer.ResponseService())

	c.restServer.Router().Route("/posts", func(router chi.Router) {
		router.Get("/", c.postController.List)
		router.Get("/{postID}", c.postController.GetByID)

		router.Group(func(protected chi.Router) {
			protected.Use(authMiddleware)

			protected.Post("/", c.postController.Create)
			protected.Put("/{postID}", c.postController.Update)
			protected.Delete("/{postID}", c.postController.Delete)
		})
	})

	return nil
}
