package utils

import "context"

type BaseInterface[T any] interface {
	Create(ctx context.Context, entity *T) error
	GetById(ctx context.Context, id uint) (*T, error)
	Update(ctx context.Context, entity *T) error
}
