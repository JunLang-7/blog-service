package model

import (
	"github.com/JunLang-7/blog-service/pkg/app"
	"gorm.io/gorm"
)

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

func (t *BlogTag) Count(db *gorm.DB) (int, error) {
	var count int64
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)
	if err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (t *BlogTag) List(db *gorm.DB, pageOffset, pageSize int) ([]*BlogTag, error) {
	var tags []*BlogTag
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	db = db.Where("state = ?", t.State)
	if err := db.Where("is_del = ?", 0).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (t *BlogTag) Create(db *gorm.DB) error {
	return db.Create(t).Error
}

func (t *BlogTag) Update(db *gorm.DB, values interface{}) error {
	return db.Model(t).Where("id = ? and is_del = ?", t.ID, 0).Updates(values).Error
}

func (t *BlogTag) Delete(db *gorm.DB) error {
	return db.Where("id = ? and is_del = ?", t.ID, 0).Delete(t).Error
}
