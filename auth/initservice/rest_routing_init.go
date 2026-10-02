package initservice

import "github.com/go-chi/chi/v5"

func (c *DependencyContainer) initRESTRouting() error {
	c.restServer.Router().Route("/auth", func(router chi.Router) {
		router.Post("/register", c.authController.Register)
		router.Post("/login", c.authController.Login)
		router.Get("/me", c.authController.Me)
	})

	return nil
}
