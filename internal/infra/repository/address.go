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
func (a *AddressRepository) GetById(ctx context.Context, id uint) (*entity.Address, error) {
	address := &entity.Address{}
	if err := a.db.WithContext(ctx).Where("id = ?", id).First(&address).Error; err != nil {
		return address, err
	}
	return address, nil

}
func (a *AddressRepository) Create(ctx context.Context, entity *entity.Address) error {
	return a.db.WithContext(ctx).Create(entity).Error
}
