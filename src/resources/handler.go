package resources

import (
	"encoding/json"

	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
)

func Localize(rootPath string) gin.HandlerFunc {
	config := &i18n.BundleCfg{
		RootPath:         rootPath,
		AcceptLanguage:   []language.Tag{language.Russian, language.English},
		DefaultLanguage:  language.English,
		UnmarshalFunc:    json.Unmarshal,
		FormatBundleFile: "json",
	}
	opt := i18n.WithBundle(config)
	return i18n.Localize(opt)
}
