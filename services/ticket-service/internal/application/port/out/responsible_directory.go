package out

import "context"

type Responsible struct {
	ID   string
	Name string
}

type ResponsibleDirectory interface {
	List(ctx context.Context) ([]string, error)
	ListWithNames(ctx context.Context) ([]Responsible, error)
	Upsert(ctx context.Context, id, name string) error
}
