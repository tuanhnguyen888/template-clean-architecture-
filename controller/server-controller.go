package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tuanhnguyen888/server/entity"
	"github.com/tuanhnguyen888/server/helper"
	"github.com/tuanhnguyen888/server/service"
	"github.com/tuanhnguyen888/server/validators"
	"gorm.io/gorm"
)

type ServerController interface {
	FindAll(ctx *gin.Context)
	Create(ctx *gin.Context) error
	Update(ctx *gin.Context) error
	Delete(ctx *gin.Context) error
	ShowAll(ctx *gin.Context)
}

type controller struct {
	ServerService service.ServerService
	LogRorate     service.LogRorate
	DB            *gorm.DB
}

var validate *validator.Validate

func NewServerController(serverService service.ServerService, logRorate service.LogRorate, db *gorm.DB) ServerController {
	validate = validator.New()
	validate.RegisterValidation("ipv4", validators.ValidateIpv4)
	return &controller{
		ServerService: serverService,
		LogRorate:     logRorate,
		DB:            db,
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
	c.ServerService.Create(server)
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
	server.UpdatedAt = time.Now().UnixMilli()
	err = validate.Struct(server)
	if err != nil {
		return err
	}
	c.ServerService.Update(server)
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
	c.ServerService.Delete(server)
	return nil
}

// FindAll implements ServerController
func (c *controller) FindAll(ctx *gin.Context) {
	// ------ pagination ----------
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}
	perPage := 10
	offset := (page - 1) * perPage

	// ----search
	servers := []entity.Server{}
	v := ctx.Query("value")
	c.DB.Where("name LIKE ?", "%"+v+"%").Order("name").Offset(offset).Limit(perPage).Find(&servers)
	if len(servers) == 0 {
		response := helper.BuildErrorResponse("Failed to process request", "not found data", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": servers,
	})

}

// ShowAll implements ServerController
func (c *controller) ShowAll(ctx *gin.Context) {
	servers := []entity.Server{}
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}
	perPage := 10
	offset := (page - 1) * perPage

	//  ------ sort -------
	sort := ctx.Query("sort")
	k := ctx.Query("kind")
	if sort != "" {
		c.DB.Order(sort + " " + k).Offset(offset).Limit(perPage).Find(&servers)
	} else {
		c.DB.Offset(offset).Limit(perPage).Find(&servers)
	}
	data := gin.H{
		"title":  "List server",
		"Server": servers,
	}
	ctx.HTML(http.StatusOK, "server.html", data)
}
