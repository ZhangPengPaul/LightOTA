package main

import (
	"log"

	"github.com/ZhangPengPaul/LightOTA/internal/config"
	"github.com/ZhangPengPaul/LightOTA/internal/handler"
	"github.com/ZhangPengPaul/LightOTA/internal/mqtt"
	"github.com/ZhangPengPaul/LightOTA/internal/repository"
	"github.com/ZhangPengPaul/LightOTA/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(sqlite.Open(cfg.Database.Path), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repo := repository.New(db)
	err = repo.AutoMigrate()
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	var mqttClient *mqtt.Client
	if cfg.MQTT.Enabled {
		mqttClient, err = mqtt.NewClient(cfg.MQTT)
		if err != nil {
			log.Printf("Warning: Failed to connect to MQTT: %v", err)
		} else {
			defer mqttClient.Disconnect()
		}
	}

	tenanthandlerSvc := service.NewTenantService(repo)
	tenantHandler := handler.NewTenantHandler(tenanthandlerSvc)

	productSvc := service.NewProductService(repo)
	productHandler := handler.NewProductHandler(productSvc)

	firmwareSvc := service.NewFirmwareService(repo, cfg.Firmware.StoragePath)
	firmwareHandler := handler.NewFirmwareHandler(firmwareSvc)

	upgradeSvc := service.NewUpgradeService(repo, repo, repo, repo, mqttClient, cfg)
	upgradeHandler := handler.NewUpgradeHandler(upgradeSvc)

	deviceHandler := handler.NewDeviceHandler(upgradeSvc)

	r := gin.Default()

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(corsConfig))

	api := r.Group("/api/v1")
	{
		api.Use(handler.AuthMiddleware(repo))
		{
			tenantHandler.Register(api)
			productHandler.Register(api)
			firmwareHandler.Register(api)
			upgradeHandler.Register(api)
		}
	}

	deviceApi := r.Group("/api/v1/ota")
	{
		deviceApi.Use(handler.DeviceAuthMiddleware())
		deviceHandler.Register(deviceApi)
	}

	log.Printf("Server starting on port %s...\n", cfg.Server.Port)
	log.Fatal(r.Run(":" + cfg.Server.Port))
}
