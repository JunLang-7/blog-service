package model

import (
	"errors"

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

func (t *BlogTag) Count(db *gorm.DB, filterState bool) (int, error) {
	var count int64
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	if filterState {
		db = db.Where("state = ?", t.State)
	}
	if err := db.Model(&t).Where("is_del = ?", 0).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (t *BlogTag) List(db *gorm.DB, pageOffset, pageSize int, filterState bool) ([]*BlogTag, error) {
	var tags []*BlogTag
	if pageOffset >= 0 && pageSize > 0 {
		db = db.Offset(pageOffset).Limit(pageSize)
	}
	if t.Name != "" {
		db = db.Where("name = ?", t.Name)
	}
	if filterState {
		db = db.Where("state = ?", t.State)
	}
	if err := db.Where("is_del = ?", 0).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (t *BlogTag) ListByIDs(db *gorm.DB, ids []uint32) ([]*BlogTag, error) {
	var tags []*BlogTag
	db = db.Where("state = ? and is_del = ?", t.State, 0)
	err := db.Where("id in (?)", ids).Find(&tags).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return tags, nil
}

func (t *BlogTag) Get(db *gorm.DB) (*BlogTag, error) {
	var tag BlogTag
	err := db.Where("id = ? and is_del = ? and state = ?", t.ID, 0, t.State).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
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
