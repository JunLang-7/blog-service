package model

type BlogTagArticle struct {
	*Model
	TagID     uint32 `json:"tag_id"`
	ArticleID uint32 `json:"article_id"`
}

func (*BlogTagArticle) TableName() string {
	return "blog_tag_article"
}
