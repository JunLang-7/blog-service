package model

import (
	"gorm.io/gorm"
)

type BlogAuth struct {
	*Model
	AppKey    string `json:"app_key"`
	AppSecret string `json:"app_secret"`
}

func (a *BlogAuth) TableName() string {
	return "blog_auth"
}

func (a *BlogAuth) Get(db *gorm.DB) (*BlogAuth, error) {
	var auth BlogAuth
	db = db.Where("app_key= ? and app_secret = ? and is_del = ?", a.AppKey, a.AppSecret, 0)
	err := db.First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}
