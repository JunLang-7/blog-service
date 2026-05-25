package service

import (
	"errors"

	"github.com/JunLang-7/blog-service/internal/dao"
	"github.com/JunLang-7/blog-service/internal/model"
	"github.com/JunLang-7/blog-service/pkg/app"
	"github.com/JunLang-7/blog-service/pkg/errcode"
)

type ArticleRequest struct {
	ID    uint32 `form:"id" json:"id" uri:"id" binding:"required,gte=1"`
	State uint8  `form:"state" json:"state" binding:"oneof=0 1"`
}

type ArticleListRequest struct {
	TagID uint32 `form:"tag_id" json:"tag_id" binding:"omitempty,gte=1"`
	State uint8  `form:"state" json:"state" binding:"omitempty,oneof=0 1"`
}

type CreateArticleRequest struct {
	TagIDs        []uint32 `form:"tag_ids" json:"tag_ids" binding:"required,min=1"`
	Title         string   `form:"title" json:"title" binding:"required,min=2,max=100"`
	Desc          string   `form:"desc" json:"desc" binding:"required,min=2,max=255"`
	Content       string   `form:"content" json:"content" binding:"required,min=2,max=4294967295"`
	CoverImageUrl string   `form:"cover_image_url" json:"cover_image_url" binding:"required,url"`
	CreatedBy     string   `form:"created_by" json:"created_by" binding:"required,min=2,max=100"`
	State         uint8    `form:"state" json:"state" binding:"oneof=0 1"`
}

type UpdateArticleRequest struct {
	ID            uint32   `form:"id" json:"id" uri:"id" binding:"required,gte=1"`
	TagIDs        []uint32 `form:"tag_ids" json:"tag_ids" binding:"required,min=1"`
	Title         string   `form:"title" json:"title" binding:"min=2,max=100"`
	Desc          string   `form:"desc" json:"desc" binding:"min=2,max=255"`
	Content       string   `form:"content" json:"content" binding:"min=2,max=4294967295"`
	CoverImageUrl string   `form:"cover_image_url" json:"cover_image_url" binding:"url"`
	ModifiedBy    string   `form:"modified_by" json:"modified_by" binding:"required,min=2,max=100"`
	State         uint8    `form:"state" json:"state" binding:"oneof=0 1"`
}

type DeleteArticleRequest struct {
	ID uint32 `form:"id" json:"id" uri:"id" binding:"required,gte=1"`
}

type Article struct {
	ID            uint32           `json:"id"`
	Title         string           `json:"title"`
	Desc          string           `json:"desc"`
	Content       string           `json:"content"`
	CoverImageUrl string           `json:"cover_image_url"`
	State         uint8            `json:"state"`
	Tags          []*model.BlogTag `json:"tags"`
}

func (s *Service) GetArticle(param *ArticleRequest, filterState bool) (*Article, error) {
	article, err := s.dao.GetArticle(param.ID, param.State, filterState)
	if err != nil {
		return nil, err
	}

	tags, err := s.getTagsByArticleID(article.ID)
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
		Tags:          tags,
	}, nil
}

func (s *Service) GetArticleList(param *ArticleListRequest, pager *app.Pager, filterTag bool, filterState bool) ([]*Article, int, error) {
	var articleIDs []uint32
	if filterTag {
		ids, err := s.dao.GetArticleIDsByTagID(param.TagID)
		if err != nil {
			return nil, 0, err
		}
		articleIDs = ids
		if len(articleIDs) == 0 {
			return []*Article{}, 0, nil
		}
	}

	totalRows, err := s.dao.CountArticleList(articleIDs, param.State, filterState)
	if err != nil {
		return nil, 0, err
	}

	articles, err := s.dao.GetArticleList(articleIDs, param.State, pager.Page, pager.PageSize, filterState)
	if err != nil {
		return nil, 0, err
	}

	articleList, err := s.enrichArticlesWithTags(articles)
	if err != nil {
		return nil, 0, err
	}

	return articleList, totalRows, nil
}

func (s *Service) CreateArticle(param *CreateArticleRequest) error {
	return s.dao.Transaction(func(txDao *dao.Dao) error {
		article, err := txDao.CreateArticle(&dao.Article{
			Title:         param.Title,
			Desc:          param.Desc,
			Content:       param.Content,
			CoverImageUrl: param.CoverImageUrl,
			State:         param.State,
			CreatedBy:     param.CreatedBy,
		})
		if err != nil {
			if errors.Is(err, dao.ErrArticleAlreadyExists) {
				return errcode.ErrorDuplicateArticle
			}
			return err
		}
		return txDao.CreateArticleTags(article.ID, param.TagIDs, param.CreatedBy)
	})
}

func (s *Service) UpdateArticle(param *UpdateArticleRequest) error {
	return s.dao.Transaction(func(txDao *dao.Dao) error {
		err := txDao.UpdateArticle(&dao.Article{
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
		return txDao.UpdateArticleTags(param.ID, param.TagIDs, param.ModifiedBy)
	})
}

func (s *Service) DeleteArticle(param *DeleteArticleRequest) error {
	return s.dao.Transaction(func(txDao *dao.Dao) error {
		err := txDao.DeleteArticle(param.ID)
		if err != nil {
			return err
		}
		return txDao.DeleteArticleTagsByArticleID(param.ID)
	})
}

func (s *Service) getTagsByArticleID(articleID uint32) ([]*model.BlogTag, error) {
	articleTags, err := s.dao.GetArticleTagListByAID(articleID)
	if err != nil {
		return nil, err
	}
	if len(articleTags) == 0 {
		return []*model.BlogTag{}, nil
	}

	tagIDs := make([]uint32, 0, len(articleTags))
	for _, at := range articleTags {
		tagIDs = append(tagIDs, at.TagID)
	}

	return s.dao.GetTagListByIDs(tagIDs, model.STATE_OPEN)
}

func (s *Service) enrichArticlesWithTags(articles []*model.BlogArticle) ([]*Article, error) {
	if len(articles) == 0 {
		return []*Article{}, nil
	}

	articleIDs := make([]uint32, 0, len(articles))
	for _, a := range articles {
		articleIDs = append(articleIDs, a.ID)
	}

	allArticleTags, err := s.dao.GetArticleTagListByAIDs(articleIDs)
	if err != nil {
		return nil, err
	}

	tagIDs := make([]uint32, 0)
	tagMap := make(map[uint32][]uint32)
	for _, at := range allArticleTags {
		tagMap[at.ArticleID] = append(tagMap[at.ArticleID], at.TagID)
		tagIDs = append(tagIDs, at.TagID)
	}

	allTags, err := s.dao.GetTagListByIDs(tagIDs, model.STATE_OPEN)
	if err != nil {
		return nil, err
	}
	tagByID := make(map[uint32]*model.BlogTag)
	for _, t := range allTags {
		tagByID[t.ID] = t
	}

	var result []*Article
	for _, a := range articles {
		var tags []*model.BlogTag
		for _, tid := range tagMap[a.ID] {
			if t, ok := tagByID[tid]; ok {
				tags = append(tags, t)
			}
		}
		if tags == nil {
			tags = []*model.BlogTag{}
		}
		result = append(result, &Article{
			ID:            a.ID,
			Title:         a.Title,
			Desc:          a.Desc,
			Content:       a.Content,
			CoverImageUrl: a.CoverImageUrl,
			State:         a.State,
			Tags:          tags,
		})
	}
	return result, nil
}
