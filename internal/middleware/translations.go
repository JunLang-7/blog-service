package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	uni        *ut.UniversalTranslator
	registerOnce sync.Once
)

func Translations() gin.HandlerFunc {
	// 共享的 uni 实例，保证 GetTranslator 返回同一 translator
	uni = ut.New(en.New(), zh.New(), zh_Hant_TW.New())

	// RegisterDefaultTranslations 修改的是全局 validator 的内部 map，
	// 只需调用一次，并发请求重复调用会导致 concurrent map write panic
	registerOnce.Do(func() {
		v, _ := binding.Validator.Engine().(*validator.Validate)
		zhTrans, _ := uni.GetTranslator("zh")
		_ = zh_translations.RegisterDefaultTranslations(v, zhTrans)
		enTrans, _ := uni.GetTranslator("en")
		_ = en_translations.RegisterDefaultTranslations(v, enTrans)
	})

	return func(c *gin.Context) {
		locale := c.GetHeader("locale")
		trans, _ := uni.GetTranslator(locale)
		c.Set("trans", trans)
		c.Next()
	}
}
