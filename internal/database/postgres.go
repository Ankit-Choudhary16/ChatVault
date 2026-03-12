package database

import (
	"embed"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Connect establishes a connection to the database using the given URL
// and runs pending migrations. Retries up to 10 times for K8s startup timing.
func Connect(databaseURL string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Database connection attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect after retries: %w", err)
	}

	upBytes, err := migrationsFS.ReadFile("migrations/000001_create_tables.up.sql")
	if err != nil {
		return nil, fmt.Errorf("failed to read UP migration file: %w", err)
	}
	upSQL := string(upBytes)

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Exec(upSQL).Error; err != nil {
		log.Printf("Migration failed: %v. Rolling back...", err)

		downBytes, downErr := migrationsFS.ReadFile("migrations/000001_create_tables.down.sql")
		if downErr != nil {
			log.Printf("Failed to read DOWN migration: %v", downErr)
			tx.Rollback()
			return nil, fmt.Errorf("migration failed and rollback failed: %w", downErr)
		}
		downSQL := string(downBytes)
		if execErr := tx.Exec(downSQL).Error; execErr != nil {
			tx.Rollback()
			return nil, fmt.Errorf("migration failed and rollback failed: %w", execErr)
		}
		tx.Rollback()
		return nil, fmt.Errorf("migration failed, rollback executed")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit migration transaction: %w", err)
	}

	log.Println("Database migrated successfully!")
	return db, nil
}
