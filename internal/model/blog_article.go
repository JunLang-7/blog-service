package model

import (
	"errors"

	"github.com/JunLang-7/blog-service/pkg/app"
	"gorm.io/gorm"
)

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

type ArticleRow struct {
	ArticleID     uint32
	TagID         uint32
	TagName       string
	ArticleTitle  string
	ArticleDesc   string
	CoverImageUrl string
	Content       string
	State         uint8
}

func (a *BlogArticle) TableName() string {
	return "blog_article"
}

func (a *BlogArticle) Get(db *gorm.DB, filterState bool) (*BlogArticle, error) {
	var article BlogArticle
	db = db.Where("id = ? and is_del = ?", a.ID, 0)
	if filterState {
		db = db.Where("state = ?", a.State)
	}
	err := db.First(&article).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &article, nil
}

func (a *BlogArticle) Create(db *gorm.DB) (*BlogArticle, error) {
	if err := db.Create(a).Error; err != nil {
		return nil, err
	}
	return a, nil
}

func (a *BlogArticle) Update(db *gorm.DB, values interface{}) error {
	return db.Model(a).Where("id = ? and is_del = ?", a.ID, 0).Updates(values).Error
}

func (a *BlogArticle) Delete(db *gorm.DB) error {
	return db.Where("id = ? and is_del = ?", a.ID, 0).Delete(a).Error
}

func (a *BlogArticle) ListByTagID(db *gorm.DB, tagID uint32, pageOffset, pageSize int, filterTag bool, filterState bool) ([]*ArticleRow, error) {
	fields := []string{"ar.id AS article_id", "ar.title AS article_title", "ar.desc AS article_desc", "ar.cover_image_url", "ar.content"}
	fields = append(fields, []string{"t.id AS tag_id", "t.name AS tag_name"}...)

	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	db = db.Select(fields).Table(new(BlogTagArticle).TableName()+" AS at").
		Joins("LEFT JOIN `"+new(BlogTag).TableName()+"` AS t ON at.tag_id = t.id").
		Joins("LEFT JOIN `"+new(BlogArticle).TableName()+"` AS ar ON at.article_id = ar.id").
		Where("ar.is_del = ?", 0)
	if filterState {
		db = db.Where("ar.state = ?", a.State)
	}
	if filterTag {
		db = db.Where("at.`tag_id` = ?", tagID)
	}
	rows, err := db.Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*ArticleRow
	for rows.Next() {
		var row ArticleRow
		if err := db.ScanRows(rows, &row); err != nil {
			return nil, err
		}
		list = append(list, &row)
	}
	return list, nil
}

func (a *BlogArticle) CountByTagID(db *gorm.DB, tagID uint32, filterTag bool, filterState bool) (int, error) {
	var count int64
	db = db.Table(new(BlogTagArticle).TableName()+" AS at").
		Joins("LEFT JOIN `"+new(BlogArticle).TableName()+"` AS ar ON at.article_id = ar.id").
		Where("ar.is_del = ?", 0)
	if filterState {
		db = db.Where("ar.state = ?", a.State)
	}
	if filterTag {
		db = db.Where("at.`tag_id` = ?", tagID)
	}
	err := db.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
