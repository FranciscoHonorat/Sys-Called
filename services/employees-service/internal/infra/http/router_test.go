package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	httpapi "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/infra/http"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/infra/memory"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestListEmployeesEndpoint(t *testing.T) {
	handler := httpapi.NewHandler(memory.NewEmployeeRepository())
	router := httpapi.NewRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/employees", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var employees []application.EmployeeOutput
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &employees))
	assert.GreaterOrEqual(t, len(employees), 3)
}
