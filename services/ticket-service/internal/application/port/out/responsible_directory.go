package out

import "context"

type ResponsibleDirectory interface {
	List(ctx context.Context) ([]string, error)
	Upsert(ctx context.Context, id, name string) error
}
