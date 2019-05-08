package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en_US"
	"github.com/go-playground/locales/ja_JP"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

func main() {
	r := gin.Default()

	trans := ut.New(ja_JP.New(), ja_JP.New(), en_US.New())
	ja, _ := trans.GetTranslator("ja_JP")
	_ = ja.Add("Hello", "こんにちは", false)

	en, _ := trans.GetTranslator("en_US")
	_ = en.Add("Hello", "Hello", false)

	validate := validator.New()
	validate.RegisterTranslation("required", ja, func(ut ut.Translator) error {
		return ut.Add("required", "{0}は必須項目です", false)
	}, TransFunc)
	validate.RegisterTranslation("required", en, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is required", false)
	}, TransFunc)

	r.GET("/hello", func(gc *gin.Context) {
		req := struct {
			Hello string `validate:"required"`
		}{}
		if err := gc.BindQuery(&req); err != nil {
			gc.Status(http.StatusBadRequest)
			return
		}
		errs := validate.Struct(req)

		verrs := errs.(validator.ValidationErrors)
		if len(verrs) > 0 {
			gc.String(http.StatusBadRequest, verrs[0].Translate(en))
			return
		}

		gc.String(http.StatusOK, "hello")
	})

	r.Run(":8888")
}

func TransFunc(ut ut.Translator, fe validator.FieldError) string {
	fld, _ := ut.T(fe.Field())
	t, err := ut.T(fe.Tag(), fld)
	if err != nil {
		return fe.(error).Error()
	}
	return t
}
