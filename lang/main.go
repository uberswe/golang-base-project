package lang

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Service struct {
	bundle    *i18n.Bundle
	ctx       *gin.Context
	localizer *i18n.Localizer
}

func New(ctx *gin.Context, bundle *i18n.Bundle) Service {
	localizer := i18n.NewLocalizer(bundle, ctx.Request.Header.Get("Accept-Language"), "en")
	return Service{
		bundle:    bundle,
		ctx:       ctx,
		localizer: localizer,
	}
}

func (s *Service) Trans(str string) string {
	// TODO, modify this to handle plural and more types of phrases
	for _, m := range translationMessages {
		if m.ID == str {
			localizedString, _ := s.localizer.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &m,
			})
			return localizedString
		} else if m.Other == str {
			localizedString, _ := s.localizer.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &m,
			})
			return localizedString
		}
	}
	return str
}
