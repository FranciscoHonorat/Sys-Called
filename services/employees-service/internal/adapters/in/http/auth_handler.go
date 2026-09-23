package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Login(c *gin.Context) {
	var body loginRequest
	if !bindJSON(c, &body) {
		return
	}

	output, err := h.useCases.Login.Execute(c.Request.Context(), application.LoginInput{
		Username: body.Username,
		Password: body.Password,
	})
	if err == nil {
		setRefreshCookie(c, output.RefreshToken, output.RefreshExpiresIn)
	}
	respondJSON(c, output, err)
}

func (h *Handler) RefreshSession(c *gin.Context) {
	refreshToken, _ := c.Cookie(refreshCookieName)
	if refreshToken == "" {
		respondError(c, domainErr.ErrInvalidRefreshToken)
		return
	}

	output, err := h.useCases.RefreshSession.Execute(c.Request.Context(), refreshToken)
	if err != nil {
		clearRefreshCookie(c)
		respondError(c, err)
		return
	}
	setRefreshCookie(c, output.RefreshToken, output.RefreshExpiresIn)
	c.JSON(http.StatusOK, output)
}

func (h *Handler) Logout(c *gin.Context) {
	if refreshToken, _ := c.Cookie(refreshCookieName); refreshToken != "" {
		if err := h.useCases.Logout.Execute(c.Request.Context(), refreshToken); err != nil {
			respondError(c, err)
			return
		}
	}
	clearRefreshCookie(c)
	c.Status(http.StatusNoContent)
}

type jwk struct {
	ID        string `json:"kid"`
	KeyType   string `json:"kty"`
	Curve     string `json:"crv"`
	Algorithm string `json:"alg"`
	Use       string `json:"use"`
	X         string `json:"x"`
}

func (h *Handler) JWKS(c *gin.Context) {
	keys := []jwk{}
	for _, k := range h.useCases.PublicKeys.Execute() {
		keys = append(keys, toJWK(k))
	}
	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

func toJWK(k out.PublicKey) jwk {
	return jwk{ID: k.ID, KeyType: k.KeyType, Curve: k.Curve, Algorithm: k.Algorithm, Use: "sig", X: k.X}
}
