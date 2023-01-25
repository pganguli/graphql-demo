package users

import (
	"errors"
	db "github.com/pganguli/hnews/internal/db"
	hash "github.com/pganguli/hnews/pkg/hash/pbkdf2"
	"gorm.io/gorm"
	"log"
)

type User struct {
	ID       int    `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

func (user *User) Create() error {
	hashedPassword, err := hash.HashPassword(user.Password)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = hashedPassword
	return db.DB.Create(user).Error
}

func (user *User) Authenticate() (bool, error) {
	var hashedPassword string

	err := db.DB.Model(&User{}).Select("password").Where("username = ?", user.Username).First(&hashedPassword).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Print(err)
		}
		return false, err
	}

	match, err := hash.CheckPasswordHash(user.Password, hashedPassword)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	return match, nil
}

// GetUserIdByUsername check if a user exists in database by given username
func GetUserIdByUsername(username string) (int, error) {
	var Id int

	err := db.DB.Model(&User{}).Select("id").Where("username = ?", username).First(&Id).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Print(err)
		}
		return -1, err
	}

	return Id, nil
}
