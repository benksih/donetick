package errorx

import (
	"errors"

	"donetick.com/core/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type LocalizedError struct {
	MessageID string
	Err       error
}

func (e *LocalizedError) Error() string {
	return e.Err.Error()
}

func (e *LocalizedError) Localize(c *gin.Context) string {
	localizer, ok := c.Get(middleware.I18nKey)
	if !ok {
		return e.Error()
	}
	i18nLocalizer, ok := localizer.(*i18n.Localizer)
	if !ok {
		return e.Error()
	}

	localizedString, err := i18nLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID: e.MessageID,
	})
	if err != nil {
		// fallback to default message
		return e.Error()
	}
	return localizedString
}

func New(messageID string, defaultText string) *LocalizedError {
	return &LocalizedError{
		MessageID: messageID,
		Err:       errors.New(defaultText),
	}
}

var ErrNotEnoughSpace = New("not_enough_storage_space", "not enough storage space")
var ErrNotAPlusMember = New("not_a_plus_member", "not a plus member: access to plus features denied")
var ErrFileSizeTooLarge = New("file_size_too_large", "file size too large")
