package command

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type SyncResponsibleInput struct {
	ID   string
	Name string
}

type SyncResponsibleUseCase struct {
	responsibles out.ResponsibleDirectory
}

func NewSyncResponsibleUseCase(responsibles out.ResponsibleDirectory) *SyncResponsibleUseCase {
	return &SyncResponsibleUseCase{responsibles: responsibles}
}

func (uc *SyncResponsibleUseCase) Execute(ctx context.Context, input SyncResponsibleInput) error {
	id, err := valueobjects.NewAssigneeID(input.ID)
	if err != nil {
		return err
	}

	return uc.responsibles.Upsert(ctx, id.GetAssigneeID(), input.Name)
}
