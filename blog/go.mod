module micro/blog

go 1.27.1

require micro/infrastructure v0.0.0

require (
	github.com/go-chi/chi/v5 v5.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
)

replace micro/infrastructure => ../infrastructure
