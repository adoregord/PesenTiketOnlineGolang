package handler

import (
	"strings"
	"time"

	"github.com/go-playground/validator"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("noblank", func(fl validator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})
	validate.RegisterValidation("Datetime", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		_, err := time.Parse("02-Jan-2006 15:04:05", dateStr)
		return err == nil
	})
}