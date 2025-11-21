package router

import (
	"net/http"

	req "auth/internal/http/rest/dto/requests"
	res "auth/internal/http/rest/dto/responses"
	locUtils "auth/internal/utils"
	shUtils "shared/utils"

	"github.com/gin-gonic/gin"
)

func (r *GinRouter) confirmEmail(c *gin.Context) {
	var in req.ConfirmEmailRequest

	parser := shUtils.NewParamParser(c)
	if parser.DecodeAndValidateJSONBody(&in).HasErrors() {
		c.JSON(http.StatusBadRequest, parser.GetErrors())
		return
	}

	if err := r.uc.ConfirmEmail(c.Request.Context(), in.UserId, in.Code); err != nil {
		msg, status := locUtils.TranslateErrorToHTTP(err)
		c.JSON(status, res.NewResponse(msg, nil, nil))
		return
	}

	c.JSON(http.StatusOK, res.NewResponse("почта подтверждена", nil, nil))
}
