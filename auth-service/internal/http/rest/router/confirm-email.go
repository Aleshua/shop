package router

import (
	"net/http"

	req "auth/internal/http/rest/dto/requests"
	res "auth/internal/http/rest/dto/responses"
	l "auth/internal/logger"
	"auth/internal/utils"

	"github.com/gin-gonic/gin"
)

func (r *GinRouter) confirmEmail(c *gin.Context) {
	var in req.ConfirmEmailRequest

	parser := utils.NewParamParser(c)
	if parser.DecodeAndValidateJSONBody(&in).HasErrors() {
		c.JSON(http.StatusBadRequest, parser.GetErrors())
		return
	}

	if err := r.uc.ConfirmEmail(c.Request.Context(), in.UserId, in.Code); err != nil {
		r.logger.With(l.NewField("body", in)).Errorf("не удалось подтвердить почту: %s", err.Error())
		msg, status := utils.TranslateErrorToHTTP(err)
		c.JSON(status, res.NewResponse(msg, nil, nil))
		return
	}

	c.JSON(http.StatusOK, res.NewResponse("почта подтверждена", nil, nil))
}
