package router

import (
	"net/http"

	req "auth/internal/http/rest/dto/requests"
	res "auth/internal/http/rest/dto/responses"
	locUtils "auth/internal/utils"
	shUtils "shared/utils"

	"github.com/gin-gonic/gin"
)

func (r *GinRouter) logout(c *gin.Context) {
	var in req.LogoutRequest

	parser := shUtils.NewParamParser(c)
	if parser.DecodeAndValidateJSONBody(&in).HasErrors() {
		c.JSON(http.StatusOK, parser.GetErrors())
		return
	}

	if err := r.uc.Logout(c.Request.Context(), in.RefreshToken); err != nil {
		msg, status := locUtils.TranslateErrorToHTTP(err)
		c.JSON(status, res.NewResponse(msg, nil, nil))
		return
	}

	c.JSON(http.StatusOK, res.NewResponse("пользователь вышел из аккаунта", nil, nil))
}
