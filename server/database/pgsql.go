package database

import (
	"fmt"
	"log"

	"github.com/Qubitopia/QuantumScholar/server/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPgsql() {
	var err error

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Kolkata",
		PGSQL_HOST, PGSQL_PORT, PGSQL_USER, PGSQL_PASSWORD, PGSQL_NAME, PGSQL_SSLMODE)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
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
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed")
}
