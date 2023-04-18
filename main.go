package main

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/dev/kode/my_gram/driver"
	"github.com/onainadapdap1/dev/kode/my_gram/models"
	"github.com/onainadapdap1/dev/kode/my_gram/seeds"
	"github.com/onainadapdap1/dev/kode/my_gram/server"
)
func drop(db *gorm.DB) {
	if err := db.DropTableIfExists(
		&models.User{}, 
		&models.Photo{},
		&models.Comment{},
		&models.SocialMedia{},).Error; err != nil {
		log.Fatalf("Error dropping tables: %v", err)
	}
}

func migrate(db *gorm.DB) {
	if err := db.Debug().AutoMigrate(
		&models.User{}, 
		&models.Photo{},
		&models.Comment{},
		&models.SocialMedia{},).Error; err != nil {
		log.Fatalf("Error migrating tables: %v", err)
	}
}

func create(database *gorm.DB) {
	drop(database)
	migrate(database)
}

func main() {
	database := driver.ConnectDB()
	defer database.Close()

	args := os.Args
	if len(args) > 1 {
		first := args[1]

		if first == "create" {
			create(database)
			log.Println("Database created successfully")
			os.Exit(0)
		} else if first == "seed" {
			seeds.Seed()
			log.Println("Database seeded successfully")
			os.Exit(0)
		} else if first == "migrate" {
			migrate(database)
			log.Println("Database migrated successfully")
			os.Exit(0)
		} else {
			log.Fatalf("Unknown command %s", first)
		}
	}

	server.StartServer()
}