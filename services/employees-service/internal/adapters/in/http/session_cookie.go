package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const refreshCookieName = "refresh_token"

func setRefreshCookie(c *gin.Context, value string, maxAge int) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     refreshCookieName,
		Value:    value,
		Path:     "/auth",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func clearRefreshCookie(c *gin.Context) {
	setRefreshCookie(c, "", -1)
}
