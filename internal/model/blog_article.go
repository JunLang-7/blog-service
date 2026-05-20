package model

import "github.com/JunLang-7/blog-service/pkg/app"

type BlogArticle struct {
	*Model
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	State         uint8  `json:"state"`
}

type BlogArticleSwagger struct {
	List  []*BlogArticle
	Pager *app.Pager
}

func (BlogArticle) TableName() string {
	return "blog_article"
}
