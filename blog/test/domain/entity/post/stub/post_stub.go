package postentitystub

import (
	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
)

const (
	RawTitle = "Как мы учили микросервисы"
	RawBody  = "Сначала был каркас, потом инфраструктура, потом auth."
)

func GetTitle() postentity.Title {
	title, err := postentity.TitleFrom(RawTitle)
	if err != nil {
		panic(err)
	}

	return title
}

func GetBody() postentity.Body {
	body, err := postentity.BodyFrom(RawBody)
	if err != nil {
		panic(err)
	}

	return body
}

func GetPostBuilder() *postentity.Builder {
	return postentity.NewBuilder().
		ID(postentity.NewPostID()).
		AuthorID(commonuserentity.NewUserID()).
		Title(GetTitle()).
		Body(GetBody())
}

func GetValidPost() *postentity.Post {
	post, err := GetPostBuilder().Build()
	if err != nil {
		panic(err)
	}

	return post
}

func GetPostAuthoredBy(authorID commonuserentity.UserID) *postentity.Post {
	post, err := postentity.NewBuilder().
		ID(postentity.NewPostID()).
		AuthorID(authorID).
		Title(GetTitle()).
		Body(GetBody()).
		Build()
	if err != nil {
		panic(err)
	}

	return post
}
