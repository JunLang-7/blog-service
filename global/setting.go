package global

import (
	"github.com/JunLang-7/blog-service/pkg/logger"
	"github.com/JunLang-7/blog-service/pkg/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	JWTSetting      *setting.JWTSettingS
	DatabaseSetting *setting.DatabaseSettingS
	RedisSetting    *setting.RedisSettingS
	EmailSetting    *setting.EmailSettingS
	Logger          *logger.Logger
)
