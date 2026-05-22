package dao

import "github.com/JunLang-7/blog-service/internal/model"

func (d *Dao) GetArticleTagByAID(articleID uint32) (*model.BlogTagArticle, error) {
	articleTag := model.BlogTagArticle{ArticleID: articleID}
	return articleTag.GetByAID(d.engine)
}

func (d *Dao) GetArticleTagListByTID(tagID uint32) ([]*model.BlogTagArticle, error) {
	articleTag := model.BlogTagArticle{TagID: tagID}
	return articleTag.ListByTID(d.engine)
}

func (d *Dao) GetArticleTAgListByAIDs(articleIDs []uint32) ([]*model.BlogTagArticle, error) {
	articleTag := model.BlogTagArticle{}
	return articleTag.ListByAIDs(d.engine, articleIDs)
}

func (d *Dao) CreateArticleTag(articleID uint32, tagID uint32, createdBy string) error {
	articleTag := model.BlogTagArticle{
		ArticleID: articleID,
		TagID:     tagID,
		Model: &model.Model{
			CreatedBy: createdBy,
		},
	}
	return articleTag.Create(d.engine)
}

func (d *Dao) UpdateArticleTag(articleID uint32, tagID uint32, updatedBy string) error {
	articleTag := model.BlogTagArticle{ArticleID: articleID}
	values := map[string]interface{}{
		"article_id":  articleID,
		"tag_id":      tagID,
		"modified_by": updatedBy,
	}
	return articleTag.UpdateOne(d.engine, values)
}

func (d *Dao) DeleteArticleTag(articleID uint32) error {
	articleTag := model.BlogTagArticle{ArticleID: articleID}
	return articleTag.DeleteOne(d.engine)
}
