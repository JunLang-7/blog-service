package upload

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/JunLang-7/blog-service/global"
	"github.com/JunLang-7/blog-service/pkg/util"
)

type FileType int

const TypeImage FileType = iota + 1

func GetFileName(name string) string {
	ext := GetFileExt(name)
	filename := strings.TrimSuffix(name, ext)
	filename = util.EncodeMD5(filename)

	return filename + ext
}

func GetFileExt(name string) string {
	return path.Ext(name)
}

func GetSavePath() string {
	return global.AppSetting.UploadSavePath
}

func GetServerUrl() string {
	return global.AppSetting.UploadServerUrl
}

// CheckSavePath 检查保存目录是否存在
func CheckSavePath(dst string) bool {
	_, err := os.Stat(dst)
	return !os.IsNotExist(err)
}

// CheckContainExt 检查文件后缀是否包含在约定的后缀配置项中
func CheckContainExt(t FileType, name string) bool {
	ext := GetFileExt(name)
	ext = strings.ToLower(ext)
	switch t {
	case TypeImage:
		for _, allowExt := range global.AppSetting.UploadImageAllowExts {
			if allowExt == ext {
				return true
			}
		}
	}
	return false
}

// CheckMaxSize 检查文件大小是否超出最大大小限制
func CheckMaxSize(t FileType, f multipart.File) bool {
	content, _ := io.ReadAll(f)
	size := len(content)
	switch t {
	case TypeImage:
		if size >= global.AppSetting.UploadImageMaxSize*1024*1024 {
			return true
		}
	}
	return false
}

// CheckPermission 检查文件权限是否足够
func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)
	return os.IsPermission(err)
}

func CreateSavePath(dst string, perm os.FileMode) error {
	return os.MkdirAll(dst, perm)
}

// SaveFile 保存所上传的文件
func SaveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
