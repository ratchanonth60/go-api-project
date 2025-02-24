package repository

import (
	"context"
	"project-api/internal/core/entity"
	"project-api/internal/core/port/utils"

	"gorm.io/gorm"
)

type IFileRepository interface {
	utils.BaseInterface[entity.File]
	BeginTransaction(ctx context.Context) *gorm.DB
	FindByKey(ctx context.Context, key string, file *entity.File) error
	FindByKeyForUpdate(ctx context.Context, key string, file *entity.File) error // New: with lock
	Update(ctx context.Context, file *entity.File) error
}
