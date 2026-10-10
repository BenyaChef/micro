package boundarymodel

import commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"

type CreatePostInput struct {
	AuthorID commonuserentity.UserID
	Title    string
	Body     string
}
