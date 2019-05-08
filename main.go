package main

import (
	"fmt"
	"log"

	"github.com/go-playground/locales/en_US"
	"github.com/go-playground/locales/ja_JP"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

func main() {
	japanese := ja_JP.New()
	english := en_US.New()
	trans := ut.New(japanese, japanese, english)

	ja, found := trans.GetTranslator("ja_JP")
	if !found {
		log.Fatal("not found translator")
	}
	err := ja.Add("Hello", "こんにちは", false)
	if err != nil {
		log.Fatal(err)
	}

	en, found := trans.GetTranslator("en_US")
	if !found {
		log.Fatal("not found translator")
	}
	err = en.Add("Hello", "Hello", false)
	if err != nil {
		log.Fatal(err)
	}

	validate := validator.New()
	validate.RegisterTranslation("required", ja,
		func(ut ut.Translator) error {
			return ut.Add("required", "{0}は必須項目です", false)
		},
		TransFn,
	)
	validate.RegisterTranslation("required", en,
		func(ut ut.Translator) error {
			return ut.Add("required", "{0} is required", false)
		},
		TransFn,
	)

	s := struct {
		Hello string `validate:"required"`
	}{}
	errs := validate.Struct(s)

	fmt.Println(errs.(validator.ValidationErrors)[0].Translate(ja))
	fmt.Println(errs.(validator.ValidationErrors)[0].Translate(en))
}

func TransFn(ut ut.Translator, fe validator.FieldError) string {
	fld, _ := ut.T(fe.Field())
	t, err := ut.T(fe.Tag(), fld)
	if err != nil {
		return fe.(error).Error()
	}
	return t
}
