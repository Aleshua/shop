package router

import (
	"net/http"

	req "auth/internal/http/rest/dto/requests"
	res "auth/internal/http/rest/dto/responses"
	l "auth/internal/logger"
	"auth/internal/utils"

	"github.com/gin-gonic/gin"
)

func (r *GinRouter) refreshToken(c *gin.Context) {
	var in req.RefreshTokenRequest

	parser := utils.NewParamParser(c)
	if parser.DecodeAndValidateJSONBody(&in).HasErrors() {
		c.JSON(http.StatusBadRequest, parser.GetErrors())
		return
	}

	accessToken, err := r.uc.RefreshToken(c.Request.Context(), in.RefreshToken)
	if err != nil {
		r.logger.With(l.NewField("body", in)).Errorf("не удалось обновить токен: %s", err.Error())
		msg, status := utils.TranslateErrorToHTTP(err)
		c.JSON(status, res.NewResponse(msg, nil, nil))
		return
	}

	c.JSON(http.StatusOK, res.NewResponse("", gin.H{"access_token": accessToken}, nil))
}
