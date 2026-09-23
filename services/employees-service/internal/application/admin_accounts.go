package application

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

func adminTarget(ctx context.Context, repo out.EmployeeRepository, caller out.Caller, id string) (employee.Employee, error) {
	if caller.Role != employee.RoleAdmin {
		return employee.Employee{}, domainErr.ErrForbidden
	}
	e, found, err := repo.FindByID(ctx, id)
	if err != nil {
		return employee.Employee{}, err
	}
	if !found {
		return employee.Employee{}, domainErr.ErrEmployeeNotFound
	}
	return e, nil
}

type ApproveEmployeeUseCase struct {
	repo out.EmployeeRepository
}

func NewApproveEmployeeUseCase(repo out.EmployeeRepository) *ApproveEmployeeUseCase {
	return &ApproveEmployeeUseCase{repo: repo}
}

func (uc *ApproveEmployeeUseCase) Execute(ctx context.Context, caller out.Caller, id string) error {
	e, err := adminTarget(ctx, uc.repo, caller, id)
	if err != nil {
		return err
	}
	e.Approve()
	return uc.repo.Save(ctx, e)
}

type IssueTemporaryPasswordUseCase struct {
	repo      out.EmployeeRepository
	hasher    out.PasswordHasher
	generator out.PasswordGenerator
	sessions  *Sessions
}

func NewIssueTemporaryPasswordUseCase(repo out.EmployeeRepository, hasher out.PasswordHasher, generator out.PasswordGenerator, sessions *Sessions) *IssueTemporaryPasswordUseCase {
	return &IssueTemporaryPasswordUseCase{repo: repo, hasher: hasher, generator: generator, sessions: sessions}
}

func (uc *IssueTemporaryPasswordUseCase) Execute(ctx context.Context, caller out.Caller, id string) (string, error) {
	e, err := adminTarget(ctx, uc.repo, caller, id)
	if err != nil {
		return "", err
	}
	password, err := uc.generator.Generate()
	if err != nil {
		return "", err
	}
	hash, err := uc.hasher.Hash(password)
	if err != nil {
		return "", err
	}
	e.SetTemporaryPassword(hash)
	if err := uc.repo.Save(ctx, e); err != nil {
		return "", err
	}
	return password, uc.sessions.EndAll(ctx, e.GetID())
}

type ChangePasswordInput struct {
	Current string
	New     string
}

type ChangePasswordUseCase struct {
	repo   out.EmployeeRepository
	hasher out.PasswordHasher
}

func NewChangePasswordUseCase(repo out.EmployeeRepository, hasher out.PasswordHasher) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{repo: repo, hasher: hasher}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, caller out.Caller, input ChangePasswordInput) error {
	e, found, err := uc.repo.FindByID(ctx, caller.ID)
	if err != nil {
		return err
	}
	if !found || !uc.hasher.Matches(e.GetPasswordHash(), input.Current) {
		return domainErr.ErrInvalidCredentials
	}
	if err := checkPasswordStrength(input.New); err != nil {
		return err
	}
	hash, err := uc.hasher.Hash(input.New)
	if err != nil {
		return err
	}
	e.ChangePassword(hash)
	return uc.repo.Save(ctx, e)
}
