package service

import (
	"context"

	"github.com/JunLang-7/blog-service/global"
	"github.com/JunLang-7/blog-service/internal/dao"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context) *Service {
	return &Service{
		ctx: ctx,
		dao: dao.New(global.DBEngine),
	}
}
