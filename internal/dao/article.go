package dao

import (
	"errors"

	"github.com/JunLang-7/blog-service/internal/model"
	"github.com/JunLang-7/blog-service/pkg/app"
)

var ErrArticleAlreadyExists = errors.New("article already exists")

type Article struct {
	ID            uint32 `json:"id"`
	TagID         uint32 `json:"tag_id"`
	Title         string `json:"title"`
	Desc          string `json:"desc"`
	Content       string `json:"content"`
	CoverImageUrl string `json:"cover_image_url"`
	CreatedBy     string `json:"created_by"`
	ModifiedBy    string `json:"modified_by"`
	State         uint8  `json:"state"`
}

func (d *Dao) CreateArticle(param *Article) (*model.BlogArticle, error) {
	article := &model.BlogArticle{
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		Model:         &model.Model{CreatedBy: param.CreatedBy},
	}
	count, err := article.Count(d.engine, false)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrArticleAlreadyExists
	}
	return article.Create(d.engine)
}

func (d *Dao) UpdateArticle(param *Article) error {
	article := &model.BlogArticle{Model: &model.Model{ID: param.ID}}
	values := map[string]interface{}{
		"modified_by": param.ModifiedBy,
		"state":       param.State,
	}
	if param.Title != "" {
		values["title"] = param.Title
	}
	if param.CoverImageUrl != "" {
		values["cover_image_url"] = param.CoverImageUrl
	}
	if param.Desc != "" {
		values["desc"] = param.Desc
	}
	if param.Content != "" {
		values["content"] = param.Content
	}
	return article.Update(d.engine, values)
}

func (d *Dao) GetArticle(id uint32, state uint8, filterState bool) (*model.BlogArticle, error) {
	article := &model.BlogArticle{Model: &model.Model{ID: id}, State: state}
	return article.Get(d.engine, filterState)
}

func (d *Dao) DeleteArticle(id uint32) error {
	article := &model.BlogArticle{Model: &model.Model{ID: id}}
	return article.Delete(d.engine)
}

func (d *Dao) GetArticleList(ids []uint32, state uint8, page, pageSize int, filterState bool) ([]*model.BlogArticle, error) {
	article := &model.BlogArticle{State: state}
	return article.ListByIDs(d.engine, ids, app.GetPageOffset(page, pageSize), pageSize, filterState)
}

func (d *Dao) CountArticleList(ids []uint32, state uint8, filterState bool) (int, error) {
	article := &model.BlogArticle{State: state}
	return article.CountByIDs(d.engine, ids, filterState)
}
