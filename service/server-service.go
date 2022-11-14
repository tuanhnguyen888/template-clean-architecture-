package service

import (
	"github.com/tuanhnguyen888/server/entity"
	"github.com/tuanhnguyen888/server/repository"
)

type ServerService interface {
	Create(server entity.Server) error
	Update(server entity.Server) error
	Delete(server entity.Server) error
	FindAll() []entity.Server
}

type serverService struct {
	serverRepo repository.ServerRepository
}

func New(serverRepository repository.ServerRepository) ServerService {
	return &serverService{
		serverRepo: serverRepository,
	}
}

func (service *serverService) Create(server entity.Server) error {
	service.serverRepo.Create(server)
	return nil
}

func (service *serverService) Update(server entity.Server) error {
	service.serverRepo.Update(server)
	return nil
}

func (service *serverService) Delete(server entity.Server) error {
	service.serverRepo.Delete(server)
	return nil
}

func (service *serverService) FindAll() []entity.Server {
	return service.serverRepo.FindAll()
}
