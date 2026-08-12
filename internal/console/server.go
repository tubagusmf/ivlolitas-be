package console

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/tubagusmf/ivlolitas-be/db"
	"github.com/tubagusmf/ivlolitas-be/internal/config"
	"github.com/tubagusmf/ivlolitas-be/internal/repository"
	"github.com/tubagusmf/ivlolitas-be/internal/storage"
	"github.com/tubagusmf/ivlolitas-be/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	handlerHttp "github.com/tubagusmf/ivlolitas-be/internal/delivery/http"
	appValidator "github.com/tubagusmf/ivlolitas-be/internal/helper"
	jwtService "github.com/tubagusmf/ivlolitas-be/internal/jwt"
)

func init() {
	rootCmd.AddCommand(serverCMD)
}

var serverCMD = &cobra.Command{
	Use:   "httpsrv",
	Short: "Start HTTP server",
	Long:  "Start the HTTP server to handle incoming requests for the to-do list application.",
	Run:   httpServer,
}

func httpServer(cmd *cobra.Command, args []string) {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	config.LoadWithViper()

	postgresDB := db.NewPostgres()
	sqlDB, err := postgresDB.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB from Gorm: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	cloudStorage, err := storage.NewCloudinaryStorage(
		config.CloudinaryCloudName(),
		config.CloudinaryAPIKey(),
		config.CloudinaryAPISecret(),
	)
	if err != nil {
		log.Fatalf("failed to initialize Cloudinary: %v", err)
	}

	userRepo := repository.NewUserRepository(postgresDB)
	roleRepo := repository.NewRoleRepository(postgresDB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(postgresDB)
	categoryRepo := repository.NewCategoryRepository(postgresDB)
	productRepo := repository.NewProductRepository(postgresDB)
	productImageRepo := repository.NewProductImageRepository(postgresDB)
	productVariantRepo := repository.NewProductVariantRepository(postgresDB)
	inventoryRepo := repository.NewInventoryRepository(postgresDB)
	inventoryTransactionRepo := repository.NewInventoryTransactionRepository(postgresDB)

	jwt := jwtService.New(os.Getenv("JWT_SECRET"))

	userUsecase := usecase.NewUserUsecase(userRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo, refreshTokenRepo, jwt)
	roleUsecase := usecase.NewRoleUsecase(roleRepo)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	productUsecase := usecase.NewProductUsecase(productRepo, categoryRepo)
	productImageUsecase := usecase.NewProductImageUsecase(productRepo, productImageRepo, cloudStorage)
	productVariantUsecase := usecase.NewProductVariantUsecase(productVariantRepo, productRepo, cloudStorage)
	inventoryUsecase := usecase.NewInventoryUsecase(inventoryRepo, inventoryTransactionRepo, postgresDB)

	authMiddleware := handlerHttp.NewAuthMiddleware(jwt)

	e := echo.New()

	e.Validator = appValidator.New()

	handlerHttp.NewUserHandler(e, userUsecase, authMiddleware)
	handlerHttp.NewAuthHandler(e, authUsecase, authMiddleware)
	handlerHttp.NewroleHandler(e, roleUsecase, authMiddleware)
	handlerHttp.NewCategoryHandler(e, categoryUsecase, authMiddleware)
	handlerHttp.NewProductHandler(e, productUsecase, authMiddleware)
	handlerHttp.NewProductImageHandler(e, productImageUsecase, authMiddleware)
	handlerHttp.NewProductVariantHandler(e, productVariantUsecase, authMiddleware)
	handlerHttp.NewInventoryHandler(e, inventoryUsecase, authMiddleware)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:3001"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	wg.Add(2)

	go func() {
		defer wg.Done()
		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}

		errCh <- e.Start(":" + port)
	}()

	go func() {
		defer wg.Done()
		<-errCh
	}()

	wg.Wait()

	if err := <-errCh; err != nil {
		if err != http.ErrServerClosed {
			logrus.Errorf("HTTP server error: %v", err)
		}
	}
}
