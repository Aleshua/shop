package router

import (
	"net/http"

	req "auth/internal/http/rest/dto/requests"
	res "auth/internal/http/rest/dto/responses"
	locUtils "auth/internal/utils"
	shUtils "shared/utils"

	"github.com/gin-gonic/gin"
)

func (r *GinRouter) refreshToken(c *gin.Context) {
	var in req.RefreshTokenRequest

	parser := shUtils.NewParamParser(c)
	if parser.DecodeAndValidateJSONBody(&in).HasErrors() {
		c.JSON(http.StatusBadRequest, parser.GetErrors())
		return
	}

	accessToken, err := r.uc.RefreshToken(c.Request.Context(), in.RefreshToken)
	if err != nil {
		msg, status := locUtils.TranslateErrorToHTTP(err)
		c.JSON(status, res.NewResponse(msg, nil, nil))
		return
	}

	c.JSON(http.StatusOK, res.NewResponse("", gin.H{"access_token": accessToken}, nil))
}
