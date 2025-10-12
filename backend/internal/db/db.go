package db

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"example.com/defect-control-system/internal/models"
)

func Connect() (*gorm.DB, error) {
	dsn := viper.GetString("database.url")
	if dsn == "" {
		return nil, fmt.Errorf("database.url is empty")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	// AutoMigrate models (add more models as they are implemented)
	if err := db.AutoMigrate(&models.User{}, &models.Project{}, &models.Defect{}, &models.Attachment{}); err != nil {
		return nil, err
	}
	return db, nil
}
