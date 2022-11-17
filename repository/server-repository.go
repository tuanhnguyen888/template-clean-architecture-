package repository

import (
	"github.com/tuanhnguyen888/server/entity"
	"gorm.io/gorm"
)

type ServerRepository interface {
	Create(server entity.Server)
	Update(server entity.Server)
	Delete(server entity.Server)
	FindAll() []entity.Server
	CloseDB()
}

type database struct {
	DB *gorm.DB
}

func NewServerRepository(db *gorm.DB) ServerRepository {
	return &database{
		DB: db,
	}
}

// CloseDB implements ServerRepository
func (db *database) CloseDB() {
	dbSQl, err := db.DB.DB()
	if err != nil {
		panic("Failed to close connection from database")
	}
	dbSQl.Close()

}

// Create implements ServerRepository
func (db *database) Create(server entity.Server) {
	db.DB.Create(&server)
}

// Update implements ServerRepository
func (db *database) Update(server entity.Server) {
	db.DB.Save(&server)
}

// Delete implements ServerRepository
func (db *database) Delete(server entity.Server) {
	db.DB.Delete(&server)
}

// FindAll implements ServerRepository
func (db *database) FindAll() []entity.Server {
	servers := []entity.Server{}
	db.DB.Set("gorm:auto_preload", true).Find(&servers)
	return servers
}
