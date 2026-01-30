package migrations

import (
	"log"

	"golang-backend/entity"

	"gorm.io/gorm"
)

func migrateRBAC(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.Role{},
		&entity.Permission{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate RBAC: %v", err)
	}
}
