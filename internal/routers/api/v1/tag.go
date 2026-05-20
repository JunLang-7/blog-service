package v1

import "github.com/gin-gonic/gin"

type Tag struct{}

func NewTag() *Tag {
	return &Tag{}
}

// Get 获取指定标签
func (t *Tag) Get(c *gin.Context) {}

// List 获取标签列表
func (t *Tag) List(c *gin.Context) {}

// Create 新增标签
func (t *Tag) Create(c *gin.Context) {}

// Update 更新指定标签
func (t *Tag) Update(c *gin.Context) {}

// Delete 删除指定标签
func (t *Tag) Delete(c *gin.Context) {}
