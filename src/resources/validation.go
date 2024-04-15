package resources

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	ru_translations "github.com/go-playground/validator/v10/translations/ru"
)

var (
	Translator *ut.UniversalTranslator
)

func init() {
	en := en.New()
	ru := ru.New()
	Translator = ut.New(en, en, ru)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en_trans, _ := Translator.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(v, en_trans)

		ru_trans, _ := Translator.GetTranslator("ru")
		ru_translations.RegisterDefaultTranslations(v, ru_trans)
	}
}
