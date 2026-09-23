package query

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
)

type ResponsibleOutput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListResponsiblesUseCase struct {
	directory out.ResponsibleDirectory
}

func NewListResponsiblesUseCase(directory out.ResponsibleDirectory) *ListResponsiblesUseCase {
	return &ListResponsiblesUseCase{directory: directory}
}

func (uc *ListResponsiblesUseCase) Execute(ctx context.Context) ([]ResponsibleOutput, error) {
	responsibles, err := uc.directory.ListWithNames(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make([]ResponsibleOutput, 0, len(responsibles))
	for _, r := range responsibles {
		outputs = append(outputs, ResponsibleOutput{ID: r.ID, Name: r.Name})
	}
	return outputs, nil
}
