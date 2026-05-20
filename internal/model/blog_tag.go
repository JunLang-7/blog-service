package model

import "github.com/JunLang-7/blog-service/pkg/app"

type BlogTag struct {
	*Model
	Name  string `json:"name"`
	State uint8  `json:"state"`
}

type BlogTagSwagger struct {
	List  []*BlogTag
	Pager *app.Pager
}

func (t *BlogTag) TableName() string {
	return "blog_tag"
}
