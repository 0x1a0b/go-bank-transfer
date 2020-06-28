package validator

import (
	"github.com/gsabadini/go-bank-transfer/infrastructure/logger"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type GoPlayground struct {
	validator *validator.Validate
	translate ut.Translator
	log       logger.Logger
	err       error
	msgs      []string
}

func NewGoPlayground(log logger.Logger) Validator {
	language := en.New()
	uni := ut.New(language, language)
	translate, found := uni.GetTranslator("en")
	if !found {
		log.Fatalln("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, translate); err != nil {
		log.Fatalln("translator not found")
	}

	return &GoPlayground{validator: v, translate: translate, log: log}
}

func (g *GoPlayground) Validate(i interface{}) error {
	errs := g.validator.Struct(i)
	if errs != nil {
		g.err = errs
		return g.err
	}

	return nil
}

func (g *GoPlayground) Messages() []string {
	if g.err != nil {
		for _, err := range g.err.(validator.ValidationErrors) {
			g.msgs = append(g.msgs, err.Translate(g.translate))
		}
	}

	return g.msgs
}
