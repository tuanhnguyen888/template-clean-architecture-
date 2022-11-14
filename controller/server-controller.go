package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tuanhnguyen888/server/entity"
	"github.com/tuanhnguyen888/server/service"
	"github.com/tuanhnguyen888/server/validators"
)

type ServerController interface {
	FindAll() []entity.Server
	Create(ctx *gin.Context) error
	Update(ctx *gin.Context) error
	Delete(ctx *gin.Context) error
	ShowAll(ctx *gin.Context)
}

type controller struct {
	service service.ServerService
}

var validate *validator.Validate

func New(service service.ServerService) ServerController {
	validate = validator.New()
	validate.RegisterValidation("ipv4", validators.ValidateIpv4)
	return &controller{
		service: service,
	}
}

// Create implements ServerController
func (c *controller) Create(ctx *gin.Context) error {
	server := entity.Server{}
	err := ctx.ShouldBindJSON(&server)
	if err != nil {
		return err
	}
	server.CreatedAt = time.Now().UnixMilli()
	server.UpdatedAt = time.Now().UnixMilli()
	err = validate.Struct(server)
	if err != nil {
		return err
	}
	c.service.Create(server)
	return nil

}

// Update implements ServerController
func (c *controller) Update(ctx *gin.Context) error {
	server := entity.Server{}
	err := ctx.ShouldBindJSON(&server)
	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		return err
	}

	server.ID = id
	err = validate.Struct(server)
	if err != nil {
		return err
	}
	c.service.Update(server)
	return nil
}

// Delete implements ServerController
func (c *controller) Delete(ctx *gin.Context) error {
	server := entity.Server{}
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		return err
	}
	server.ID = id
	c.service.Delete(server)
	return nil
}

// FindAll implements ServerController
func (c *controller) FindAll() []entity.Server {
	return c.service.FindAll()
}

// ShowAll implements ServerController
func (c *controller) ShowAll(ctx *gin.Context) {
	servers := c.service.FindAll()
	data := gin.H{
		"title":  "List server",
		"Server": servers,
	}
	ctx.HTML(http.StatusOK, "server.html", data)
}
