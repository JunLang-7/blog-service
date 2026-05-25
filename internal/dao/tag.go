package dao

import (
	"errors"

	"github.com/JunLang-7/blog-service/internal/model"
	"github.com/JunLang-7/blog-service/pkg/app"
)

var ErrTagAlreadyExists = errors.New("tag already exists")

func (d *Dao) CountTag(name string, state uint8, filterState bool) (int, error) {
	tag := model.BlogTag{Name: name, State: state}
	return tag.Count(d.engine, filterState)
}

func (d *Dao) GetTag(id uint32, state uint8) (*model.BlogTag, error) {
	tag := model.BlogTag{Model: &model.Model{ID: id}, State: state}
	return tag.Get(d.engine)
}

func (d *Dao) GetTagList(name string, state uint8, page, pageSize int, filterState bool) ([]*model.BlogTag, error) {
	tag := model.BlogTag{Name: name, State: state}
	pageOffset := app.GetPageOffset(page, pageSize)
	return tag.List(d.engine, pageOffset, pageSize, filterState)
}

func (d *Dao) GetTagListByIDs(ids []uint32, state uint8) ([]*model.BlogTag, error) {
	tag := model.BlogTag{State: state}
	return tag.ListByIDs(d.engine, ids)
}

func (d *Dao) CreateTag(name string, state uint8, createBy string) error {
	tag := model.BlogTag{Name: name}
	count, err := tag.Count(d.engine, false)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrTagAlreadyExists
	}

	newTag := model.BlogTag{
		Name:  name,
		State: state,
		Model: &model.Model{
			CreatedBy: createBy,
		},
	}
	return newTag.Create(d.engine)
}

func (d *Dao) UpdateTag(id uint32, name string, state uint8, updateBy string) error {
	tag := model.BlogTag{Model: &model.Model{ID: id}}
	values := map[string]interface{}{
		"state":       state,
		"modified_by": updateBy,
	}
	if name != "" {
		values["name"] = name
	}
	return tag.Update(d.engine, values)
}

func (d *Dao) DeleteTag(id uint32) error {
	tag := model.BlogTag{Model: &model.Model{ID: id}}
	return tag.Delete(d.engine)
}
