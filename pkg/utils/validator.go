package utils

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/yorukot/go-template/pkg/logger"
)

var LangList = []string{"en", "es", "zh-tw", "zh-cn"}

type LangLocal string

const (
	English   LangLocal = "en"
	Spanish   LangLocal = "es"
	ChineseTW LangLocal = "zh-tw"
	ChineseCN LangLocal = "zh-cn"
)

var langLocalValidator validator.Func = func(fl validator.FieldLevel) bool {
	lang := fl.Field().String()
	for _, validLang := range LangList {
		if lang == validLang {
			return true
		}
	}
	return false
}

var usernameValidator validator.Func = func(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	matched, _ := regexp.MatchString(`^[a-z0-9._]+$`, username)
	return matched
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("lang", langLocalValidator)
		v.RegisterValidation("username", usernameValidator)
	} else {
		logger.Log.Fatal("error registering validator")
	}
}
