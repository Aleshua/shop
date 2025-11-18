package router

import (
	"net/http"

	req "auth/internal/http/rest/dto/requests"
	res "auth/internal/http/rest/dto/responses"
	l "auth/internal/logger"
	"auth/internal/utils"

	"github.com/gin-gonic/gin"
)

func (r *GinRouter) logout(c *gin.Context) {
	var in req.LogoutRequest

	parser := utils.NewParamParser(c)
	if parser.DecodeAndValidateJSONBody(&in).HasErrors() {
		c.JSON(http.StatusOK, parser.GetErrors())
		return
	}

	if err := r.uc.Logout(c.Request.Context(), in.RefreshToken); err != nil {
		r.logger.With(l.NewField("body", in)).Errorf("не удалось выйти из аккаунта: %s", err.Error())
		msg, status := utils.TranslateErrorToHTTP(err)
		c.JSON(status, res.NewResponse(msg, nil, nil))
		return
	}

	c.JSON(http.StatusOK, res.NewResponse("пользователь вышел из аккаунта", nil, nil))
}
