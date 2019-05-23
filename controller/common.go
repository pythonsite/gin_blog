package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	SESSION_KEY          = "UserID"       // session key
	CONTEXT_USER_KEY     = "User"         // context user key
	SESSION_GITHUB_STATE = "GITHUB_STATE" // github state session key
	SESSION_CAPTCHA      = "GIN_CAPTCHA"  // captcha session key
)

func writeJson(ctx *gin.Context, h gin.H) {
	if _, ok := h["succeed"]; !ok {
		h["succeed"] =false
	}
	ctx.JSON(http.StatusOK, h)
}