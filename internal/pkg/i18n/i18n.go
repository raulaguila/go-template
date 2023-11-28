package i18n

import (
	"errors"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var I18nTranslations map[string]*Translation = map[string]*Translation{}

type Translation struct {
	language string

	ErrGeneric            error
	ErrInvalidId          error
	ErrInvalidDatas       error
	ErrManyRequest        error
	ErrorNonexistentRoute error
	ErrUndefinedColumn    error
	ErrExpiredToken       error
	ErrDisabledUser       error
	ErrIncorrectPassword  error
	ErrPassUnmatch        error
	ErrUserHasPass        error

	ErrProductUsed       error
	ErrProductNotFound   error
	ErrProductRegistered error

	ErrProfileUsed       error
	ErrProfileNotFound   error
	ErrProfileRegistered error

	ErrUserUsed       error
	ErrUserNotFound   error
	ErrUserRegistered error
}

func (s *Translation) SetLanguage(lang string) {
	s.language = lang
}

func (s *Translation) SetTranslations(localizer *i18n.Localizer) {
	s.ErrGeneric = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrGeneric"}, PluralCount: 1}))
	s.ErrInvalidId = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrInvalidId"}, PluralCount: 1}))
	s.ErrInvalidDatas = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrInvalidDatas"}, PluralCount: 1}))
	s.ErrManyRequest = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrManyRequest"}, PluralCount: 1}))
	s.ErrorNonexistentRoute = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrorNonexistentRoute"}, PluralCount: 1}))
	s.ErrUndefinedColumn = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrUndefinedColumn"}, PluralCount: 1}))
	s.ErrExpiredToken = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrExpiredToken"}, PluralCount: 1}))
	s.ErrDisabledUser = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrDisabledUser"}, PluralCount: 1}))
	s.ErrIncorrectPassword = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrIncorrectPassword"}, PluralCount: 1}))
	s.ErrPassUnmatch = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrPassUnmatch"}, PluralCount: 1}))
	s.ErrUserHasPass = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrUserHasPass"}, PluralCount: 1}))

	s.ErrProductUsed = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrProductUsed"}, PluralCount: 1}))
	s.ErrProductNotFound = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrProductNotFound"}, PluralCount: 1}))
	s.ErrProductRegistered = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrProductRegistered"}, PluralCount: 1}))

	s.ErrProfileUsed = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrProfileUsed"}, PluralCount: 1}))
	s.ErrProfileNotFound = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrProfileNotFound"}, PluralCount: 1}))
	s.ErrProfileRegistered = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrProfileRegistered"}, PluralCount: 1}))

	s.ErrUserUsed = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrUserUsed"}, PluralCount: 1}))
	s.ErrUserNotFound = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrUserNotFound"}, PluralCount: 1}))
	s.ErrUserRegistered = errors.New(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: "ErrUserRegistered"}, PluralCount: 1}))
}
