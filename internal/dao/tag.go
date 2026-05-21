package dao

import (
	"github.com/JunLang-7/blog-service/internal/model"
	"github.com/JunLang-7/blog-service/pkg/app"
)

func (d *Dao) CountTag(name string, state uint8) (int, error) {
	tag := model.BlogTag{Name: name, State: state}
	return tag.Count(d.engine)
}

func (d *Dao) GetTagList(name string, state uint8, page, pageSize int) ([]*model.BlogTag, error) {
	tag := model.BlogTag{Name: name, State: state}
	pageOffset := app.GetPageOffset(page, pageSize)
	return tag.List(d.engine, pageOffset, pageSize)
}

func (d *Dao) CreateTag(name string, state uint8, createBy string) error {
	tag := model.BlogTag{
		Name:  name,
		State: state,
		Model: &model.Model{
			CreatedBy: createBy,
		},
	}
	return tag.Create(d.engine)
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
