package v1

import (
	"github.com/JunLang-7/blog-service/pkg/app"
	"github.com/JunLang-7/blog-service/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type Article struct{}

func NewArticle() *Article {
	return &Article{}
}

// Get 获取指定文章
func (a *Article) Get(c *gin.Context) {
	app.NewResponse(c).ToErrorResponse(errcode.ServerError)
	return
}

// List 获取文章列表
func (a *Article) List(c *gin.Context) {}

// Create 新增文章
func (a *Article) Create(c *gin.Context) {}

// Update 更新指定文章
func (a *Article) Update(c *gin.Context) {}

// Delete 删除指定文章
func (a *Article) Delete(c *gin.Context) {}
