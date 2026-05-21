package dao

import "github.com/JunLang-7/blog-service/internal/model"

func (d *Dao) GetAuth(appKey, appSecret string) (*model.BlogAuth, error) {
	auth := &model.BlogAuth{AppKey: appKey, AppSecret: appSecret}
	return auth.Get(d.engine)
}
