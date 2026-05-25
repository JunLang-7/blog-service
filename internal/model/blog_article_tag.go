package model

import (
	"errors"

	"gorm.io/gorm"
)

type BlogTagArticle struct {
	*Model
	TagID     uint32 `json:"tag_id"`
	ArticleID uint32 `json:"article_id"`
}

func (a *BlogTagArticle) TableName() string {
	return "blog_article_tag"
}

func (a *BlogTagArticle) GetByAID(db *gorm.DB) (*BlogTagArticle, error) {
	var articleTag BlogTagArticle
	err := db.Where("article_id = ? and is_del = ?", a.ArticleID, 0).First(&articleTag).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &articleTag, nil
}

func (a *BlogTagArticle) ListByAID(db *gorm.DB) ([]*BlogTagArticle, error) {
	var articleTags []*BlogTagArticle
	err := db.Where("article_id = ? and is_del = ?", a.ArticleID, 0).Find(&articleTags).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return articleTags, nil
}

func (a *BlogTagArticle) ListByTID(db *gorm.DB) ([]*BlogTagArticle, error) {
	var articleTags []*BlogTagArticle
	if err := db.Where("tag_id = ? and is_del = ?", a.TagID, 0).Find(&articleTags).Error; err != nil {
		return nil, err
	}
	return articleTags, nil
}

func (a *BlogTagArticle) ListByAIDs(db *gorm.DB, articleIDs []uint32) ([]*BlogTagArticle, error) {
	var articleTags []*BlogTagArticle
	err := db.Where("article_id IN (?) and is_del = ?", articleIDs, 0).Find(&articleTags).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return articleTags, nil
}

func (a *BlogTagArticle) Create(db *gorm.DB) error {
	return db.Create(a).Error
}

func (a *BlogTagArticle) UpdateOne(db *gorm.DB, values interface{}) error {
	return db.Model(a).Where("article_id = ? and is_del = ?", a.ArticleID, 0).Limit(1).Updates(values).Error
}

func (a *BlogTagArticle) Delete(db *gorm.DB) error {
	return db.Where("id = ? and is_del = ?", a.ID, 0).Delete(a).Error
}

func (a *BlogTagArticle) DeleteOne(db *gorm.DB) error {
	return db.Where("id = ? and is_del = ?", a.ID, 0).Delete(a).Limit(1).Error
}

func (a *BlogTagArticle) DeleteByArticleID(db *gorm.DB) error {
	return db.Where("article_id = ? and is_del = ?", a.ArticleID, 0).Delete(a).Error
}
