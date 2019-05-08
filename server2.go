package main

import (
	"log"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
)

func test3() {
	r := gin.Default()

	r.GET("/hello", func(gc *gin.Context) {
		req := struct {
			Hello int `binding:"required"`
		}{}
		if err := gc.BindQuery(&req); err != nil {
			log.Println(reflect.TypeOf(err))
			switch err.(type) {
			case validator.ValidationErrors:
				verrs := err.(validator.ValidationErrors)
				for _, v := range verrs {
					log.Printf("%+v\n", v)
				}
			default:
				gc.String(http.StatusBadRequest, err.Error())
			}
			return
		}

		gc.String(http.StatusOK, "hello")
	})

	r.Run(":8888")
}
