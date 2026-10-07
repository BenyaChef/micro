package postrepomodel

import (
	"time"

	postentity "github.com/BenyaChef/micro/blog/domain/entity/post"
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	timeprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/time"
)

type Post struct {
	ID        string    `db:"id"`
	AuthorID  string    `db:"author_id"`
	Title     string    `db:"title"`
	Body      string    `db:"body"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func ToModel(post *postentity.Post) *Post {
	return &Post{
		ID:        post.ID().String(),
		AuthorID:  post.AuthorID().String(),
		Title:     post.Title().String(),
		Body:      post.Body().String(),
		CreatedAt: post.CreatedAt().Time(),
		UpdatedAt: post.UpdatedAt().Time(),
	}
}

func ToEntity(model *Post) (*postentity.Post, error) {
	postID, err := postentity.PostIDFrom(model.ID)
	if err != nil {
		return nil, err
	}

	authorID, err := commonuserentity.UserIDFrom(model.AuthorID)
	if err != nil {
		return nil, err
	}

	title, err := postentity.TitleFrom(model.Title)
	if err != nil {
		return nil, err
	}

	body, err := postentity.BodyFrom(model.Body)
	if err != nil {
		return nil, err
	}

	return postentity.NewBuilder().
		ID(postID).
		AuthorID(authorID).
		Title(title).
		Body(body).
		CreatedAt(timeprimitive.TimestampFrom(model.CreatedAt)).
		UpdatedAt(timeprimitive.TimestampFrom(model.UpdatedAt)).
		Build()
}

func ToEntities(models []*Post) ([]*postentity.Post, error) {
	posts := make([]*postentity.Post, 0, len(models))

	for _, model := range models {
		post, err := ToEntity(model)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
