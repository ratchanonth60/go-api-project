package repository

import "context"

type BaseInterface interface {
	Create(ctx context.Context, entity interface{}) error
}
