package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const I18nKey = "i18n"

func Localizer(bundle *i18n.Bundle) gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.GetHeader("Accept-Language")
		localizer := i18n.NewLocalizer(bundle, lang)
		c.Set(I18nKey, localizer)
		c.Next()
	}
}
