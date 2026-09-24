module micro/common

go 1.27.1

require (
	github.com/google/uuid v1.6.0
	micro/infrastructure v0.0.0
)

replace micro/infrastructure => ../infrastructure
