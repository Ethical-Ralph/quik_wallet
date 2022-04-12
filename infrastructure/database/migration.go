package migration

import (
	"github.com/Ethical-Ralph/quik_wallet/infrastructure/database/models"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	db.Migrator().CreateTable(&models.Wallet{})
}
