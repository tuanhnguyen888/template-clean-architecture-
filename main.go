package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tuanhnguyen888/server/controller"
	"github.com/tuanhnguyen888/server/database"
	"github.com/tuanhnguyen888/server/middlewares"
	"github.com/tuanhnguyen888/server/repository"
	"github.com/tuanhnguyen888/server/service"
)

var (
	db                 = database.NewInitPG()
	redis, errRedis    = database.NewInitResdis()
	elasic, errElastic = database.NewInitEsClient()

	serverRepo repository.ServerRepository = repository.NewServerRepository(db)
	userRepo   repository.UserRepository   = repository.NewUserRepository(db)

	serverService    service.ServerService = service.NewServerService(serverRepo)
	authService      service.AuthService   = service.NewAuthService(userRepo)
	jwtService       service.JWTService    = service.NewJWTService()
	logRorateService service.LogRorate     = service.NewLogrorate()

	serverController   controller.ServerController   = controller.NewServerController(serverService, logRorateService, db)
	excelController    controller.ExcelController    = controller.NewExcelController(db, serverService)
	authController     controller.AuthController     = controller.NewAuthController(authService, jwtService)
	periodicController controller.PeriodicController = controller.NewPeriodicController(serverService, logRorateService, db, redis, elasic)
)

func main() {
	defer serverRepo.CloseDB()
	if errElastic != nil {
		log.Panic(errElastic)
	}

	if errRedis != nil {
		log.Panic(errRedis)
	}

	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(gin.Logger())

	app.Static("/css", "./templates/css")

	app.LoadHTMLGlob("templates/*.html")

	// auth
	apiAuth := app.Group("/api/auth")
	{
		apiAuth.POST("/register", authController.Register)
		apiAuth.POST("/login", authController.Login)
		apiAuth.GET("/logout", authController.LogoutUser)
	}

	apiRoutes := app.Group("/api")
	apiRoutes.Use(middlewares.AuthorizeJWT())
	{
		apiRoutes.GET("servers", serverController.FindAll)

		apiRoutes.POST("server", func(ctx *gin.Context) {
			err := serverController.Create(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Add Success!",
			})
		})

		apiRoutes.PUT("/server/:id", func(ctx *gin.Context) {
			err := serverController.Update(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Update Success!"})
			}

		})

		apiRoutes.DELETE("/server/:id", func(ctx *gin.Context) {
			err := serverController.Delete(ctx)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				ctx.JSON(http.StatusOK, gin.H{"message": "Deleted!"})
			}

		})

		apiRoutes.POST("/exportServer", excelController.ExportServer)
		apiRoutes.POST("/importServer", excelController.ImportServer)
	}

	viewRoutes := app.Group("/view")
	{
		viewRoutes.GET("/servers", serverController.ShowAll)
	}

	go periodicController.Cron()

	// We can setup this env variable from the EB console
	port := os.Getenv("PORT")
	// Elastic Beanstalk forwards requests to port 5000
	if port == "" {
		port = "3000"
	}
	app.Run(":" + port)
}

// JWT
