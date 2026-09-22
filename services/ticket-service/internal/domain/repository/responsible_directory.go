package repository

import "context"

type ResponsibleDirectory interface {
	List(ctx context.Context) ([]string, error)
}
