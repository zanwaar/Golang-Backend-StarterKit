package migrations

import (
	"log"

	"golang-backend/entity"

	"gorm.io/gorm"
)

func migrateUsers(db *gorm.DB) {
	err := db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatalf("Failed to migrate Users: %v", err)
	}

	// Manual Indexing for Smart Search
	// Adding GIN index for Full-Text Search performance
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_users_fulltext ON users USING GIN (to_tsvector('indonesian', name || ' ' || email));").Error
	if err != nil {
		log.Printf("Warning: Failed to create full-text index: %v", err)
	} else {
		log.Println("✓ Full-text search index ensured")
	}
}
