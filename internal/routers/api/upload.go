package api

import (
	"github.com/JunLang-7/blog-service/global"
	"github.com/JunLang-7/blog-service/internal/service"
	"github.com/JunLang-7/blog-service/pkg/app"
	"github.com/JunLang-7/blog-service/pkg/convert"
	"github.com/JunLang-7/blog-service/pkg/errcode"
	"github.com/JunLang-7/blog-service/pkg/upload"
	"github.com/gin-gonic/gin"
)

type Upload struct{}

func NewUpload() *Upload {
	return new(Upload)
}

func (u *Upload) UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	fileType := convert.StrTo(c.PostForm("type")).MustInt()
	if header == nil || fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, header)
	if err != nil {
		global.Logger.Errorf(c, "svc.UploadFile err: %v", err)
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}
	response.ToResponse(gin.H{
		"file_access_url": fileInfo.AccessUrl,
	})
}

func (u *Upload) UploadFiles(c *gin.Context) {
	response := app.NewResponse(c)
	form, err := c.MultipartForm()
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}
	headers := form.File["files"]

	fileType := convert.StrTo(c.PostForm("type")).MustInt()
	if fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	svc := service.New(c.Request.Context())
	var resp []string
	for _, header := range headers {
		file, err := header.Open()
		if err != nil {
			global.Logger.Errorf(c, "svc.UploadFiles err: %v", err)
			response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
			return
		}
		fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, header)
		_ = file.Close()
		if err != nil {
			global.Logger.Errorf(c, "svc.UploadFile err: %v", err)
			response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
			return
		}
		resp = append(resp, fileInfo.AccessUrl)
	}
	response.ToResponse(gin.H{
		"file_access_urls": resp,
	})
}
