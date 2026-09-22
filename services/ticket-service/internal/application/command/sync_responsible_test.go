package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

func TestSyncResponsibleUseCase(t *testing.T) {
	t.Run("should register a new responsible", func(t *testing.T) {
		directory := &outtest.ResponsibleDirectory{}
		uc := command.NewSyncResponsibleUseCase(directory)

		err := uc.Execute(context.Background(), command.SyncResponsibleInput{ID: "agent-1", Name: "Ana Souza"})

		require.NoError(t, err)
		assert.Equal(t, []string{"agent-1"}, directory.Responsibles)
		assert.Equal(t, "Ana Souza", directory.Names["agent-1"])
	})

	t.Run("should be idempotent for a duplicated event", func(t *testing.T) {
		directory := &outtest.ResponsibleDirectory{}
		uc := command.NewSyncResponsibleUseCase(directory)
		input := command.SyncResponsibleInput{ID: "agent-1", Name: "Ana Souza"}

		require.NoError(t, uc.Execute(context.Background(), input))
		require.NoError(t, uc.Execute(context.Background(), input))

		assert.Equal(t, []string{"agent-1"}, directory.Responsibles)
	})

	t.Run("should reject an empty ID", func(t *testing.T) {
		directory := &outtest.ResponsibleDirectory{}
		uc := command.NewSyncResponsibleUseCase(directory)

		err := uc.Execute(context.Background(), command.SyncResponsibleInput{Name: "Ana Souza"})

		assert.ErrorIs(t, err, domainErr.ErrInvalidAssignee)
		assert.Empty(t, directory.Responsibles)
	})

	t.Run("should propagate directory errors", func(t *testing.T) {
		directory := &outtest.ResponsibleDirectory{Err: assert.AnError}
		uc := command.NewSyncResponsibleUseCase(directory)

		err := uc.Execute(context.Background(), command.SyncResponsibleInput{ID: "agent-1", Name: "Ana Souza"})

		assert.ErrorIs(t, err, assert.AnError)
	})
}
