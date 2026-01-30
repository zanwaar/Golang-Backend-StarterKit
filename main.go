package main

import (
	"flag"
	"io"
	"log"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"

	"golang-backend/config"
	"golang-backend/controller"
	"golang-backend/middleware"
	"golang-backend/migrations"
	"golang-backend/repository"
	"golang-backend/routes"
	"golang-backend/service"
	"golang-backend/utils"

	"github.com/gin-gonic/gin"
)

// @title           Golang Backend API
// @version         1.0
// @description     This is a sample server for Golang Backend.
// @host            localhost:8080
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Load environment variables
	config.LoadEnv()

	// Setup logging with rotation (lumberjack)
	logFile := &lumberjack.Logger{
		Filename:   "logs/server.log",
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	// Set gin writer to both file and console
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(gin.DefaultWriter)

	// Setup slog for structured logging (JSON)
	logger := slog.New(slog.NewJSONHandler(gin.DefaultWriter, nil))
	slog.SetDefault(logger)

	// 2. Initialize database
	db := config.InitDB()

	// 3. Run Migrations (Optional via flag)
	migrate := flag.Bool("migrate", false, "Run database migrations")
	seed := flag.Bool("seed", false, "Run database seeder")
	flag.Parse()

	if *migrate {
		migrations.RunMigrations(db)
	}

	// 4. Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// 4. Initialize services
	userService := service.NewUserService(userRepo)

	// 5. Initialize controllers
	userCtrl := controller.NewUserController(userService)
	roleCtrl := controller.NewRoleController(db)

	// Run Seeder
	if *seed {
		utils.SeedRolesAndPermissions(db)

		utils.SeedUsers(db)
	}

	// 6. Initialize Gin router
	app := gin.Default()

	// 7. Apply global middleware
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.ErrorHandlerMiddleware())
	app.Use(middleware.RateLimiterMiddleware())

	// 8. Setup routes
	routes.SetupRoutes(app, userCtrl, roleCtrl)

	// 9. Health check endpoint
	app.GET("/health", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			utils.ErrorResponse(c, "Database connection failed", 500, nil)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			utils.ErrorResponse(c, "Database ping failed", 500, nil)
			return
		}

		utils.SuccessResponse(c, "API is running and DB is connected", nil)
	})

	// 10. Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := app.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
