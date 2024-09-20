package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"mxshop-api/userop-web/global"
	"reflect"
	"strings"
)

func InitTrans(locale string) (err error) {
	// 修改gin框架中的Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器

		uni := ut.New(enT, zhT, enT)
		var ok bool
		// 也可以使用 uni.FindTranslator(...) 传入多个locale进行查找
		global.Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}

		// 注册翻译器
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(v, global.Trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, global.Trans)
		default:
			en_translations.RegisterDefaultTranslations(v, global.Trans)
		}
		return
	}
	return
}
