package database

import (
	"fmt"
	"log"
	"time"

	"github.com/Qubitopia/quantum-scholar-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPgsql() {
	var err error

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		PGSQL_HOST, PGSQL_PORT, PGSQL_USER, PGSQL_PASSWORD, PGSQL_NAME, PGSQL_SSLMODE)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
}

func dropConflictingTables() {
	err := DB.Migrator().DropTable(
		&models.User{},
		&models.MagicLink{},
		&models.Test{},
		&models.TestAssignedToUser{},
		&models.PaymentTable{},
		&models.Answer{},
	)
	if err != nil {
		log.Fatal("Failed to drop tables:", err)
	}
	log.Println("Conflicting tables dropped successfully")
}

func MigratePgsql() {
	// Migrate tables in dependency order to avoid foreign key errors
	err := DB.AutoMigrate(
		&models.User{},
		&models.MagicLink{},
		&models.Test{},
		&models.TestAssignedToUser{},
		&models.PaymentTable{},
		&models.Answer{},
	)
	if err != nil {
		if GIN_MODE == "release" {
			log.Fatal("Failed to migrate database:", err)
			log.Println("In production mode, not attempting to drop tables. Exiting.")
			return
		}
		
		log.Println("Failed to migrate database:", err)

		log.Println("Attempting to drop conflicting tables and retry migration in 10 seconds as the GIN_MODE is not 'release'")
		time.Sleep(10 * time.Second)
		log.Println("10 seconds elapsed. Dropping tables now.")

		dropConflictingTables()
		err := DB.AutoMigrate(
			&models.User{},
			&models.MagicLink{},
			&models.Test{},
			&models.TestAssignedToUser{},
			&models.PaymentTable{},
			&models.Answer{},
		)
		if err != nil {
			log.Fatal("Failed to migrate database even after dropping tables:", err)
		}
	}
	log.Println("Database migration completed")
}
