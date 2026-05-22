package service

import (
	"github.com/JunLang-7/blog-service/internal/dao"
	"github.com/JunLang-7/blog-service/internal/model"
	"github.com/JunLang-7/blog-service/pkg/app"
)

type ArticleRequest struct {
	ID    uint32 `form:"id" binding:"required,gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type ArticleListRequest struct {
	TagID uint32 `form:"tag_id" binding:"gte=1"`
	State uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type CreateArticleRequest struct {
	TagID         uint32 `form:"tag_id" binding:"required,gte=1"`
	Title         string `form:"title" binding:"required,min=2,max=100"`
	Desc          string `form:"desc" binding:"required,min=2,max=255"`
	Content       string `form:"content" binding:"required,min=2,max=4294967295"`
	CoverImageUrl string `form:"cover_image_url" binding:"required,url"`
	CreatedBy     string `form:"created_by" binding:"required,min=2,max=100"`
	State         uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type UpdateArticleRequest struct {
	ID            uint32 `form:"id" binding:"required,gte=1"`
	TagID         uint32 `form:"tag_id" binding:"required,gte=1"`
	Title         string `form:"title" binding:"min=2,max=100"`
	Desc          string `form:"desc" binding:"min=2,max=255"`
	Content       string `form:"content" binding:"min=2,max=4294967295"`
	CoverImageUrl string `form:"cover_image_url" binding:"url"`
	ModifiedBy    string `form:"modified_by" binding:"required,min=2,max=100"`
	State         uint8  `form:"state,default=1" binding:"oneof=0 1"`
}

type DeleteArticleRequest struct {
	ID uint32 `form:"id" binding:"required,gte=1"`
}

type Article struct {
	ID            uint32         `json:"id"`
	Title         string         `json:"title"`
	Desc          string         `json:"desc"`
	Content       string         `json:"content"`
	CoverImageUrl string         `json:"cover_image_url"`
	State         uint8          `json:"state"`
	Tag           *model.BlogTag `json:"tag"`
}

func (s *Service) GetArticle(param *ArticleRequest) (*Article, error) {
	article, err := s.dao.GetArticle(param.ID, param.State)
	if err != nil {
		return nil, err
	}

	articleTag, err := s.dao.GetArticleTagByAID(article.ID)
	if err != nil {
		return nil, err
	}

	tag, err := s.dao.GetTag(articleTag.TagID, model.STATE_OPEN)
	if err != nil {
		return nil, err
	}

	return &Article{
		ID:            article.ID,
		Title:         article.Title,
		Desc:          article.Desc,
		Content:       article.Content,
		CoverImageUrl: article.CoverImageUrl,
		State:         article.State,
		Tag:           tag,
	}, nil
}

func (s *Service) GetArticleList(param *ArticleListRequest, pager *app.Pager) ([]*Article, int, error) {
	articleCount, err := s.dao.CountArticleListByTagID(param.TagID, param.State)
	if err != nil {
		return nil, 0, err
	}

	articles, err := s.dao.GetArticleListByTagID(param.TagID, param.State, pager.Page, pager.PageSize)
	if err != nil {
		return nil, 0, err
	}

	var articleList []*Article
	for _, article := range articles {
		articleList = append(articleList, &Article{
			ID:            article.ArticleID,
			Title:         article.ArticleTitle,
			Desc:          article.ArticleDesc,
			Content:       article.Content,
			CoverImageUrl: article.CoverImageUrl,
			Tag:           &model.BlogTag{Model: &model.Model{ID: article.TagID}, Name: article.TagName},
		})
	}
	return articleList, articleCount, nil
}

func (s *Service) CreateArticle(param *CreateArticleRequest) error {
	article, err := s.dao.CreateArticle(&dao.Article{
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		CreatedBy:     param.CreatedBy,
	})
	if err != nil {
		return err
	}

	return s.dao.CreateArticleTag(article.ID, param.TagID, param.CreatedBy)
}

func (s *Service) UpdateArticle(param *UpdateArticleRequest) error {
	err := s.dao.UpdateArticle(&dao.Article{
		ID:            param.ID,
		Title:         param.Title,
		Desc:          param.Desc,
		Content:       param.Content,
		CoverImageUrl: param.CoverImageUrl,
		State:         param.State,
		ModifiedBy:    param.ModifiedBy,
	})
	if err != nil {
		return err
	}

	return s.dao.UpdateArticleTag(param.ID, param.TagID, param.ModifiedBy)
}

func (s *Service) DeleteArticle(param *DeleteArticleRequest) error {
	err := s.dao.DeleteArticle(param.ID)
	if err != nil {
		return err
	}
	return s.dao.DeleteArticleTag(param.ID)
}
