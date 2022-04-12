package database

import (
	"github.com/Ethical-Ralph/quik_wallet/infrastructure/database/models"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {
	err := db.Migrator().CreateTable(&models.Wallet{})
	if err != nil {
		return err
	}

	return nil
}
