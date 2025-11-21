package utils

import (
	"context"

	"github.com/gin-gonic/gin"
)

type ParamParser struct {
	errors  map[string]string
	request *gin.Context
}

func NewParamParser(c *gin.Context) *ParamParser {
	return &ParamParser{
		request: c,
		errors:  make(map[string]string),
	}
}

func (p *ParamParser) HasErrors() bool {
	return len(p.errors) > 0
}

func (p *ParamParser) GetErrors() map[string]string {
	return p.errors
}

type StructWithValidate interface {
	Validate(ctx context.Context) error
}

func (p *ParamParser) DecodeAndValidateJSONBody(dest StructWithValidate) *ParamParser {
	if err := p.request.ShouldBindJSON(dest); err != nil {
		p.errors["body"] = "не удалось прочитать body"
		return p
	}

	err := dest.Validate(p.request)
	if err != nil {
		p.errors = MergeMaps(p.errors, TranslateValidatorError(err))
	}

	return p
}
