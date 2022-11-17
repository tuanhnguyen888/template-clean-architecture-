package service

import (
	"github.com/tuanhnguyen888/server/entity"
	"github.com/tuanhnguyen888/server/repository"
)

type AuthService interface {
	// Login(username string, password string) bool
	IsDuplicateEmail(email string) bool
	CreateUser(user entity.User) error
	VerifyCredential(email string, password string) error
}

type authService struct {
	authRepo repository.UserRepository
}

func NewAuthService(authRepository repository.UserRepository) AuthService {
	return &authService{
		authRepo: authRepository,
	}
}

// func (service *loginService) Login(username string, password string) bool {
// 	return service.authorizedUsername == username && service.authorizedPassword == password
// }

func (service *authService) IsDuplicateEmail(email string) bool {
	res := service.authRepo.IsDuplicateEmail(email)
	return !(res.Error == nil)
}

func (service *authService) CreateUser(user entity.User) error {
	service.authRepo.CreateUser(user)
	return nil
}

func (service *authService) VerifyCredential(email string, password string) error {
	return service.authRepo.VerifyCredential(email, password)
}
