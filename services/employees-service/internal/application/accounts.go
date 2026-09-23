package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

const (
	minimumPasswordLength = 8
	maximumPasswordBytes  = 72
)

func checkPasswordStrength(password string) error {
	switch {
	case len([]rune(password)) < minimumPasswordLength:
		return domainErr.ErrWeakPassword
	case len(password) > maximumPasswordBytes:
		return domainErr.ErrPasswordTooLong
	}
	return nil
}

type SignUpInput struct {
	Name     string
	Username string
	Password string
}

type SignUpUseCase struct {
	repo   out.EmployeeRepository
	hasher out.PasswordHasher
}

func NewSignUpUseCase(repo out.EmployeeRepository, hasher out.PasswordHasher) *SignUpUseCase {
	return &SignUpUseCase{repo: repo, hasher: hasher}
}

func (uc *SignUpUseCase) Execute(ctx context.Context, input SignUpInput) error {
	if err := checkPasswordStrength(input.Password); err != nil {
		return err
	}
	hash, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return err
	}
	e, err := employee.SignUp("user-"+uuid.NewString(), input.Name, input.Username, hash)
	if err != nil {
		return err
	}
	return uc.repo.Register(ctx, e)
}

type RequestPasswordResetUseCase struct {
	repo out.EmployeeRepository
}

func NewRequestPasswordResetUseCase(repo out.EmployeeRepository) *RequestPasswordResetUseCase {
	return &RequestPasswordResetUseCase{repo: repo}
}

func (uc *RequestPasswordResetUseCase) Execute(ctx context.Context, username string) error {
	e, found, err := uc.repo.FindByUsername(ctx, username)
	if err != nil || !found || !e.CanLogIn() {
		return err
	}
	e.RequestPasswordReset()
	return uc.repo.Save(ctx, e)
}
