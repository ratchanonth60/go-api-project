package repository

import (
	"context"

	"project-api/internal/core/entity"
	"project-api/internal/core/port/repository"

	"gorm.io/gorm"
)

type AddressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) repository.IAddressRepository {
	return &AddressRepository{
		db: db,
	}
}

func (a *AddressRepository) Create(ctx context.Context, address *entity.Address) error {
	return a.db.WithContext(ctx).Create(&address).Error
}
