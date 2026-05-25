package dao

import "github.com/JunLang-7/blog-service/internal/model"

func (d *Dao) GetArticleTagListByAID(articleID uint32) ([]*model.BlogTagArticle, error) {
	articleTag := model.BlogTagArticle{ArticleID: articleID}
	return articleTag.ListByAID(d.engine)
}

func (d *Dao) GetArticleTagListByTID(tagID uint32) ([]*model.BlogTagArticle, error) {
	articleTag := model.BlogTagArticle{TagID: tagID}
	return articleTag.ListByTID(d.engine)
}

func (d *Dao) GetArticleTagListByAIDs(articleIDs []uint32) ([]*model.BlogTagArticle, error) {
	articleTag := model.BlogTagArticle{}
	return articleTag.ListByAIDs(d.engine, articleIDs)
}

func (d *Dao) GetArticleIDsByTagID(tagID uint32) ([]uint32, error) {
	list, err := d.GetArticleTagListByTID(tagID)
	if err != nil {
		return nil, err
	}
	ids := make([]uint32, 0, len(list))
	for _, at := range list {
		ids = append(ids, at.ArticleID)
	}
	return ids, nil
}

func (d *Dao) CreateArticleTags(articleID uint32, tagIDs []uint32, createdBy string) error {
	for _, tagID := range tagIDs {
		articleTag := model.BlogTagArticle{
			ArticleID: articleID,
			TagID:     tagID,
			Model:     &model.Model{CreatedBy: createdBy},
		}
		if err := articleTag.Create(d.engine); err != nil {
			return err
		}
	}
	return nil
}

func (d *Dao) UpdateArticleTags(articleID uint32, tagIDs []uint32, modifiedBy string) error {
	if err := d.DeleteArticleTagsByArticleID(articleID); err != nil {
		return err
	}
	return d.CreateArticleTags(articleID, tagIDs, modifiedBy)
}

func (d *Dao) DeleteArticleTagsByArticleID(articleID uint32) error {
	articleTag := model.BlogTagArticle{ArticleID: articleID}
	return articleTag.DeleteByArticleID(d.engine)
}
