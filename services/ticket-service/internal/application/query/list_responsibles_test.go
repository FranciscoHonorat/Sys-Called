package query_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/query"
)

func TestListResponsiblesUseCase(t *testing.T) {
	t.Run("should list the responsibles with their names in registration order", func(t *testing.T) {
		directory := &outtest.ResponsibleDirectory{}
		require.NoError(t, directory.Upsert(context.Background(), "agent-1", "Ana Souza"))
		require.NoError(t, directory.Upsert(context.Background(), "agent-2", "Bruno Lima"))

		output, err := query.NewListResponsiblesUseCase(directory).Execute(context.Background())

		require.NoError(t, err)
		assert.Equal(t, []query.ResponsibleOutput{
			{ID: "agent-1", Name: "Ana Souza"},
			{ID: "agent-2", Name: "Bruno Lima"},
		}, output)
	})

	t.Run("should propagate directory errors", func(t *testing.T) {
		directory := &outtest.ResponsibleDirectory{Err: assert.AnError}

		_, err := query.NewListResponsiblesUseCase(directory).Execute(context.Background())

		assert.ErrorIs(t, err, assert.AnError)
	})
}
