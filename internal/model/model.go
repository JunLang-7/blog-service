package model

import (
	"fmt"

	"github.com/JunLang-7/blog-service/global"
	"github.com/JunLang-7/blog-service/pkg/app"
	"github.com/JunLang-7/blog-service/pkg/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type BlogTagSwagger struct {
	List  []*BlogTag
	Pager *app.Pager
}

type BlogArticleSwagger struct {
	List  []*BlogArticle
	Pager *app.Pager
}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*gorm.DB, error) {
	s := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	db, err := gorm.Open(mysql.Open(fmt.Sprintf(s,
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	if global.ServerSetting.RunMode == "debug" {
		db.Logger.LogMode(logger.Info)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)

	return db, nil
}
