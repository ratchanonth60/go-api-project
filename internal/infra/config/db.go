package config

import (
	"fmt"

	"project-api/internal/core/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	DB     *gorm.DB
	Config *gorm.Config
}

func (g *GormDB) Connect() error {
	c := Config.Database

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		c.Host, c.Port, c.User, c.Password, c.DBName)
	db, err := gorm.Open(postgres.Open(dsn), g.Config)
	if err != nil {
		return err
	}

	g.DB = db
	models := []interface{}{
		&entity.User{},
		&entity.Address{},
		&entity.File{},
	}
	if err := db.AutoMigrate(models...); err != nil {
		return nil
	}
	return nil
}
