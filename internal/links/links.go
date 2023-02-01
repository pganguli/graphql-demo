package links

import (
	"errors"
	db "github.com/pganguli/hnews/internal/db"
	"github.com/pganguli/hnews/internal/users"
	"gorm.io/gorm"
	"log"
)

type Link struct {
	ID      int    `gorm:"primaryKey"`
	Title   string `gorm:"not null"`
	Address string `gorm:"not null"`
	UserID  int
	User    users.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (link Link) Save() (int, error) {
	err := db.DB.Create(&link).Error
	return link.ID, err
}

func GetAll() []Link {
	var links []Link

	result := db.DB.Preload("User").Find(&links)
	if err := result.Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Print(err)
		}
	}

	return links
}
