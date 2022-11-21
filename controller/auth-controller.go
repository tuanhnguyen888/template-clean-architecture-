package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuanhnguyen888/server/entity"
	"github.com/tuanhnguyen888/server/helper"
	"github.com/tuanhnguyen888/server/service"
)

type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	LogoutUser(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
	jWtService  service.JWTService
}

func NewAuthController(authService service.AuthService,
	jWtService service.JWTService) AuthController {
	return &authController{
		authService: authService,
		jWtService:  jWtService,
	}
}

func (c *authController) Register(ctx *gin.Context) {
	registerUser := entity.User{}
	err := ctx.ShouldBind(&registerUser)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if !c.authService.IsDuplicateEmail(registerUser.Email) {
		response := helper.BuildErrorResponse("Failed to process request", "Duplicate email", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
		return
	}

	hashedPassword, err := entity.HashPassword(registerUser.Password)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", "Error hashed password", helper.EmptyObj{})
		ctx.JSON(http.StatusConflict, response)
		return
	}

	registerUser.Password = hashedPassword

	err = c.authService.CreateUser(registerUser)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process create user", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	token, _ := c.jWtService.GenerateToken(registerUser.Email, true)
	ctx.JSON(http.StatusOK, gin.H{
		"data":  registerUser,
		"token": token,
	})

}

func (c *authController) Login(ctx *gin.Context) {
	loginUser := entity.User{}
	err := ctx.ShouldBind(&loginUser)
	if err != nil {
		response := helper.BuildErrorResponse("Failed to process request", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	authResult := c.authService.VerifyCredential(loginUser.Email, loginUser.Password)

	if authResult == nil {
		token, _ := c.jWtService.GenerateToken(loginUser.Email, true)
		response := helper.BuildResponse(true, "OK!", token)
		ctx.JSON(http.StatusOK, response)
		return
	}
	response := helper.BuildErrorResponse("Please check again your credential", "Invalid Credential", helper.EmptyObj{})
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func (c *authController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("token", "", 150, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
