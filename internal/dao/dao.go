package dao

import "gorm.io/gorm"

type Dao struct {
	engine *gorm.DB
}

func New(engine *gorm.DB) *Dao {
	return &Dao{engine: engine}
}

func (d *Dao) Transaction(fn func(txDao *Dao) error) error {
	return d.engine.Transaction(func(tx *gorm.DB) error {
		return fn(&Dao{engine: tx})
	})
}
