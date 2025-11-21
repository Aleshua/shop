package router

import (
	"net/http"

	req "auth/internal/http/rest/dto/requests"
	res "auth/internal/http/rest/dto/responses"
	locUtils "auth/internal/utils"
	shUtils "shared/utils"

	"github.com/gin-gonic/gin"
)

func (r *GinRouter) register(c *gin.Context) {
	var in req.RegisterRequest

	parser := shUtils.NewParamParser(c)

	if parser.DecodeAndValidateJSONBody(&in).HasErrors() {
		c.JSON(http.StatusBadRequest, parser.GetErrors())
		return
	}

	userId, err := r.uc.Register(c.Request.Context(), in.Password, in.Email)
	if err != nil {
		msg, status := locUtils.TranslateErrorToHTTP(err)
		c.JSON(status, res.NewResponse(msg, nil, nil))
		return
	}

	c.JSON(
		http.StatusOK,
		res.NewResponse("пользователь зарегистрирован. код подтверждения отправлен на почту", gin.H{"id": userId}, nil),
	)
}
