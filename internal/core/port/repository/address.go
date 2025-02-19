package repository

import (
	"context"

	"project-api/internal/core/entity"
)

type IAddressRepository interface {
	Create(ctx context.Context, address *entity.Address) error
}
