package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tuanhnguyen888/server/service"
)

type LoginController interface {
	Login(ctx *gin.Context) string
}

type loginController struct {
	loginService service.LoginService
	jWtService   service.JWTService
}

func NewLoginController(loginService service.LoginService,
	jWtService service.JWTService) LoginController {
	return &loginController{
		loginService: loginService,
		jWtService:   jWtService,
	}
}

func (controller *loginController) Login(ctx *gin.Context) string {
	account := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := ctx.ShouldBind(&account)
	if err != nil {
		return ""
	}
	isAuthenticated := controller.loginService.Login(account.Username, account.Password)
	if isAuthenticated {
		x, _ := controller.jWtService.GenerateToken(account.Username, true)
		return x
	}
	return ""
}
