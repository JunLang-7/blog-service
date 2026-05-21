package service

import (
	"errors"
	"mime/multipart"
	"os"

	"github.com/JunLang-7/blog-service/global"
	"github.com/JunLang-7/blog-service/pkg/upload"
)

type FileInfo struct {
	Name      string
	AccessUrl string
}

func (s *Service) UploadFile(fileType upload.FileType, file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error) {
	filename := upload.GetFileName(fileHeader.Filename)
	if !upload.CheckContainExt(fileType, filename) {
		return nil, errors.New("file suffix is not supported")
	}
	if upload.CheckMaxSize(fileType, file) {
		return nil, errors.New("file size exceeds limit")
	}

	uploadSavePath := upload.GetSavePath()
	if !upload.CheckSavePath(uploadSavePath) {
		if err := upload.CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
			return nil, errors.New("failed to create save directory")
		}
	}
	if upload.CheckPermission(uploadSavePath) {
		return nil, errors.New("permission denied")
	}

	dst := uploadSavePath + "/" + filename
	if err := upload.SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}

	accessUrl := global.AppSetting.UploadServerUrl + "/" + filename
	return &FileInfo{
		Name:      filename,
		AccessUrl: accessUrl,
	}, nil
}
