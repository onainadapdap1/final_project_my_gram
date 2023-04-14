package seeds

import (
	"log"
	"math/rand"

	"github.com/icrowley/fake"
	"github.com/jinzhu/gorm"
	"github.com/onainadapdap1/dev/kode/my_gram/driver"
	"github.com/onainadapdap1/dev/kode/my_gram/helpers"
	"github.com/onainadapdap1/dev/kode/my_gram/models"
)

func randomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func SeedUsers(db *gorm.DB) {
	users := []models.User{}
	usersToSeed := 0
	tx := db.Begin()
	tx.Model(&models.User{}).Find(&users).Count(&usersToSeed)
	usersToSeed = 5 - usersToSeed
	if usersToSeed > 0 {
		for i := 0; i < usersToSeed; i++ {
			password, _ := helpers.HassPass("password")
			user := models.User {
				Username: fake.UserName(),
				Email: fake.EmailAddress(),
				Password: password,
				Age: randomInt(8, 80),
			}

			err := tx.Debug().Create(&user).Error

			if err != nil {
				tx.Rollback()
				log.Println(err)
			}
		}		
	}

	tx.Commit()
}

func Seed() {
	db := driver.ConnectDB()
	SeedUsers(db)
}
