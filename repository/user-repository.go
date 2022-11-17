package repository

import (
	"fmt"

	"github.com/tuanhnguyen888/server/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	IsDuplicateEmail(email string) (tx *gorm.DB)
	CreateUser(user entity.User)
	VerifyCredential(email string, password string) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &database{
		DB: db,
	}
}

func (db *database) IsDuplicateEmail(email string) (tx *gorm.DB) {
	user := entity.User{}
	return db.DB.Where("email = ?", email).Take(&user)
}

// CreateUser implements ServerRepository
func (db *database) CreateUser(user entity.User) {
	db.DB.Create(&user)
}

func (db *database) VerifyCredential(email string, password string) error {
	user := entity.User{}
	result := db.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return fmt.Errorf("invalid email or Password")
	}
	if err := entity.VerifyPassword(user.Password, password); err != nil {
		return fmt.Errorf("invalid email or Password")
	}

	return nil
}
