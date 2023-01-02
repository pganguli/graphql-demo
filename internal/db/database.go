package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func Open() {
	dbUrl := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		log.Panic(err)
	}

	DB = db
}
