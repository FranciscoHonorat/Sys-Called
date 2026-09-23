package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httpapi "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/in/http"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out/outtest"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestRouter() *gin.Engine {
	repo := &outtest.EmployeeRepository{Employees: []employee.Employee{
		employee.NewEmployee("agent-1", "Ana Souza", "ana", employee.RoleSupport, "hashed:senha123"),
		employee.NewEmployee("agent-2", "Bruno Lima", "bruno", employee.RoleSupport, "hashed:senha123"),
		employee.NewEmployee("agent-3", "Carla Melo", "carla", employee.RoleSupport, "hashed:senha123"),
	}}
	sessions := application.NewSessions(
		outtest.TokenIssuer{}, &outtest.RefreshTokenGenerator{}, outtest.NewRefreshTokenStore(), 7*24*time.Hour, time.Now,
	)
	verifier := outtest.TokenVerifier{Callers: map[string]out.Caller{
		"admin-token": {ID: "admin-1", Role: employee.RoleAdmin},
		"agent-token": {ID: "agent-1", Role: employee.RoleSupport},
	}}
	hasher := &outtest.PasswordHasher{}
	return httpapi.NewRouter(httpapi.NewHandler(httpapi.UseCases{
		Authenticate:           application.NewAuthenticateUseCase(verifier),
		SignUp:                 application.NewSignUpUseCase(repo, hasher),
		RequestPasswordReset:   application.NewRequestPasswordResetUseCase(repo),
		ApproveEmployee:        application.NewApproveEmployeeUseCase(repo),
		IssueTemporaryPassword: application.NewIssueTemporaryPasswordUseCase(repo, hasher, outtest.PasswordGenerator{Password: "Temp-1234"}, sessions),
		ChangePassword:         application.NewChangePasswordUseCase(repo, hasher),
		ListEmployees:          application.NewListEmployeesUseCase(repo),
		Login:                  application.NewLoginUseCase(repo, hasher, sessions),
		RefreshSession:         application.NewRefreshSessionUseCase(repo, sessions),
		Logout:                 application.NewLogoutUseCase(sessions),
		PublicKeys:             application.NewGetPublicKeysUseCase(outtest.SigningKeys{}),
	}))
}

func TestJWKSEndpoint(t *testing.T) {
	w := doRequest(newTestRouter(), http.MethodGet, "/.well-known/jwks.json", "")

	require.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"keys":[{"kid":"key-1","kty":"OKP","crv":"Ed25519","alg":"EdDSA","use":"sig","x":"public-x"}]}`, w.Body.String())
}

func doRequest(router *gin.Engine, method, path, body string, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		req.AddCookie(c)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestListEmployeesEndpoint(t *testing.T) {
	get := func(authorization string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/employees", nil)
		if authorization != "" {
			req.Header.Set("Authorization", authorization)
		}
		w := httptest.NewRecorder()
		newTestRouter().ServeHTTP(w, req)
		return w
	}

	t.Run("lists the employees with username and role for admins", func(t *testing.T) {
		w := get("Bearer admin-token")

		require.Equal(t, http.StatusOK, w.Code)
		var employees []application.EmployeeOutput
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &employees))
		assert.Contains(t, employees, application.EmployeeOutput{ID: "agent-1", Name: "Ana Souza", Username: "ana", Role: "support", Status: "active"})
	})

	t.Run("is forbidden to anyone else", func(t *testing.T) {
		assert.Equal(t, http.StatusForbidden, get("Bearer agent-token").Code)
	})

	t.Run("requires authentication", func(t *testing.T) {
		assert.Equal(t, http.StatusUnauthorized, get("").Code)
		assert.Equal(t, http.StatusUnauthorized, get("Bearer forged").Code)
	})
}

func TestLoginEndpoint(t *testing.T) {
	t.Run("should return a bearer token for valid credentials", func(t *testing.T) {
		w := doRequest(newTestRouter(), http.MethodPost, "/auth/login", `{"username":"ana","password":"senha123"}`)

		require.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"access_token":"token-for-agent-1","token_type":"Bearer","expires_in":900}`, w.Body.String())
	})

	t.Run("should return 401 with a generic message for invalid credentials", func(t *testing.T) {
		w := doRequest(newTestRouter(), http.MethodPost, "/auth/login", `{"username":"ana","password":"wrong"}`)

		require.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"error":"invalid credentials"}`, w.Body.String())
	})

	t.Run("should return 400 for a malformed body", func(t *testing.T) {
		w := doRequest(newTestRouter(), http.MethodPost, "/auth/login", `not json`)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}
