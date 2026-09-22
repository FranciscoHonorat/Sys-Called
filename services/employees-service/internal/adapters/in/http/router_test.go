package httpapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	httpapi "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/in/http"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type fakeEmployeeRepository struct {
	employees []employee.Employee
}

func (r *fakeEmployeeRepository) List(_ context.Context) ([]employee.Employee, error) {
	return r.employees, nil
}

func (r *fakeEmployeeRepository) Register(_ context.Context, e employee.Employee) error {
	r.employees = append(r.employees, e)
	return nil
}

func TestListEmployeesEndpoint(t *testing.T) {
	handler := httpapi.NewHandler(application.NewListEmployeesUseCase(&fakeEmployeeRepository{employees: []employee.Employee{
		employee.NewEmployee("agent-1", "Ana Souza"),
		employee.NewEmployee("agent-2", "Bruno Lima"),
		employee.NewEmployee("agent-3", "Carla Melo"),
	}}))
	router := httpapi.NewRouter(handler)

	req := httptest.NewRequest(http.MethodGet, "/employees", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var employees []application.EmployeeOutput
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &employees))
	assert.GreaterOrEqual(t, len(employees), 3)
}
