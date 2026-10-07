package postentity

import (
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	timeprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/time"
)

type Post struct {
	id        PostID
	authorID  commonuserentity.UserID
	title     Title
	body      Body
	createdAt timeprimitive.Timestamp
	updatedAt timeprimitive.Timestamp
}

func (p *Post) ID() PostID {
	return p.id
}

func (p *Post) AuthorID() commonuserentity.UserID {
	return p.authorID
}

func (p *Post) Title() Title {
	return p.title
}

func (p *Post) Body() Body {
	return p.body
}

func (p *Post) CreatedAt() timeprimitive.Timestamp {
	return p.createdAt
}

func (p *Post) UpdatedAt() timeprimitive.Timestamp {
	return p.updatedAt
}

func (p *Post) IsAuthoredBy(userID commonuserentity.UserID) bool {
	return p.authorID == userID
}
