package main

import (
	v1 "centralService/internal/api/v1"
	"centralService/internal/controllers/authcontroller"
	"centralService/internal/controllers/equipmentcontroller"
	"centralService/internal/controllers/usercontroller"
	"centralService/internal/http/middlewares"
	"centralService/internal/service/authservice"
	"centralService/internal/service/equipmentservice"
	"centralService/internal/storage"
	"centralService/internal/storage/authstorage"
	"centralService/internal/storage/db"
	"centralService/internal/storage/equipmentstorage"
	"time"

	"log/slog"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	project_state := os.Getenv("PROJECT_STATE")

	logger_storage := storage.NewLoggerStorage(project_state)
	logger_storage.SetLogger()

	dbUri := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")
	collUsers := os.Getenv("DB_USERS")
	secretkey := os.Getenv("SECRET_KEY")

	//webhook_url := os.Getenv("ALERT_WEBHOOK")
	//dbCollaborators := os.Getenv("DB_COLLABORATORS")
	//dbShipments := os.Getenv("DB_SHIPMENTS")
	collEquipments := os.Getenv("DB_EQUIPMENTS")
	google_client_id := os.Getenv("GOOGLE_CLIENT_ID")

	slog.Info(".env varibles loaded", "db-name", dbName, "db-url", dbUri)

	client, err := db.NewClientDB(dbUri, dbName)

	if err != nil {
		slog.Error("ERROR: occurred at db", "error", err.Error())
		return
	}

	router := gin.Default()

	userStorage := authstorage.NewUserStorage(client, collUsers)

	equipmentStorage := equipmentstorage.NewStorage(client, collEquipments)

	err = equipmentStorage.InitIndexes()

	if err != nil {
		slog.Error("ERROR: occurred at init index db", "error", err.Error())
		return
	}
	err = userStorage.InitIndexes()

	if err != nil {
		slog.Error("ERROR: occurred at init index db", "error", err.Error())
		return
	}

	tokenStorage := authstorage.NewTokenStorage(secretkey)
	authStorage := authstorage.NewAuthStorage(google_client_id)

	userService := authservice.NewUserService(userStorage)
	tokenService := authservice.NewTokenService(tokenStorage)
	authService := authservice.NewAuthService(userService, authStorage, tokenService)

	equipmentService := equipmentservice.NewService(equipmentStorage)

	authController := authcontroller.NewAuthController(authService)
	equipmentController := equipmentcontroller.NewController(equipmentService)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "GET", "PATCH", "PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	api := router.Group("/api/v1")

	v1.RegisterUserRoutes(api, usercontroller.NewUserController(userService))

	v1.RegisterAuthRoutes(api, authController)

	v1.RegisterEquipmentsRoutes(api, equipmentController, middlewares.NewAuthenticateRole(tokenService, userService))

	router.Run("localhost:5000")

}
