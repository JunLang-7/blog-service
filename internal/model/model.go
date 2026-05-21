package model

import (
	"fmt"
	"time"

	"github.com/JunLang-7/blog-service/global"
	"github.com/JunLang-7/blog-service/pkg/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Model struct {
	ID         uint32 `gorm:"primary_key" json:"id"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	CreatedOn  uint32 `json:"created_on"`
	ModifiedOn uint32 `json:"modified_on"`
	DeletedOn  uint32 `json:"deleted_on"`
	IsDel      uint8  `json:"is_del"`
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

	_ = db.Callback().Create().Before("gorm:create").Register("update_timestamps_for_create", updateTimeStampForCreateCallback)
	_ = db.Callback().Update().Before("gorm:update").Register("update_timestamps_for_update", updateTimeStampForUpdateCallback)
	_ = db.Callback().Delete().Before("gorm:delete").Register("soft_delete", deleteCallback)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	sqlDB.SetMaxOpenConns(databaseSetting.MaxOpenConns)

	return db, nil
}

func updateTimeStampForCreateCallback(db *gorm.DB) {
	if db.Error != nil || db.Statement.Schema == nil {
		return
	}
	now := time.Now().Unix()
	if field := db.Statement.Schema.LookUpField("CreatedOn"); field != nil {
		if _, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue); isZero {
			db.Statement.SetColumn("CreatedOn", now, true)
		}
	}
	if field := db.Statement.Schema.LookUpField("ModifiedOn"); field != nil {
		if _, isZero := field.ValueOf(db.Statement.Context, db.Statement.ReflectValue); isZero {
			db.Statement.SetColumn("ModifiedOn", now, true)
		}
	}
}

func updateTimeStampForUpdateCallback(db *gorm.DB) {
	if db.Error != nil || db.Statement.Schema == nil {
		return
	}
	if _, ok := db.Get("gorm:update_column"); ok {
		return
	}
	db.Statement.SetColumn("ModifiedOn", time.Now().Unix(), true)
}

func deleteCallback(db *gorm.DB) {
	if db.Error != nil || db.Statement.Schema == nil {
		return
	}

	if db.Statement.Unscoped {
		return
	}

	_, hasDeletedOn := db.Statement.Schema.FieldsByName["DeletedOn"]
	_, hasIsDel := db.Statement.Schema.FieldsByName["IsDel"]
	if !hasDeletedOn || !hasIsDel {
		return
	}

	now := time.Now().Unix()

	db.Statement.AddClause(clause.Update{})
	db.Statement.AddClause(clause.Set{{
		Column: clause.Column{Name: "deleted_on"},
		Value:  now,
	}, {
		Column: clause.Column{Name: "is_del"},
		Value:  1,
	}})

	if whereClause, ok := db.Statement.Clauses["WHERE"]; ok {
		w := whereClause.Expression.(clause.Where)
		w.Exprs = append(w.Exprs, clause.Eq{
			Column: clause.Column{Name: "is_del"},
			Value:  0,
		})
		whereClause.Expression = w
		db.Statement.Clauses["WHERE"] = whereClause
	} else {
		db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{
			clause.Eq{Column: clause.Column{Name: "is_del"}, Value: 0},
		}})
	}

	db.Statement.BuildClauses = []string{"UPDATE", "SET", "WHERE"}
}
