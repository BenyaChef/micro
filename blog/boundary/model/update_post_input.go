package boundarymodel

import (
	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
)

type UpdatePostInput struct {
	PostID   postentity.PostID
	AuthorID commonuserentity.UserID
	Title    string
	Body     string
}
