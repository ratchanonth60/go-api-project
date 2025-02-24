package repository

import (
	"context"
	"project-api/internal/core/entity"
	"project-api/internal/core/port/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) repository.IFileRepository {
	return &FileRepository{
		db: db,
	}
}
func (f *FileRepository) BeginTransaction(ctx context.Context) *gorm.DB {
	return f.db.WithContext(ctx).Begin()
}
func (f *FileRepository) Create(ctx context.Context, entity *entity.File) error {
	return f.db.WithContext(ctx).Create(entity).Error
}

func (f *FileRepository) GetById(ctx context.Context, id uint) (*entity.File, error) {
	file := &entity.File{}
	if err := f.db.WithContext(ctx).Where("id = ?", id).First(file).Error; err != nil {
		return file, err
	}
	return file, nil
}

func (f *FileRepository) FindByKey(ctx context.Context, key string, file *entity.File) error {
	return f.db.WithContext(ctx).Where("file_path = ? AND is_deleted = ?", key, false).First(file).Error
}

func (f *FileRepository) FindByKeyForUpdate(ctx context.Context, key string, file *entity.File) error {
	return f.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", key, false).Clauses(clause.Locking{Strength: "UPDATE"}).First(file).Error
}

func (f *FileRepository) Update(ctx context.Context, file *entity.File) error {
	return f.db.WithContext(ctx).Save(file).Error
}
