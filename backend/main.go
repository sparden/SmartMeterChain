package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/smartmeterchain/backend/config"
	"github.com/smartmeterchain/backend/database"
	"github.com/smartmeterchain/backend/handlers"
	"github.com/smartmeterchain/backend/middleware"
	"github.com/smartmeterchain/backend/models"
	"github.com/smartmeterchain/backend/services"
	"github.com/smartmeterchain/backend/utils"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg := config.Load()

	// Connect database
	database.Connect(cfg.DBPath)

	// Initialize services
	fabricSvc := services.NewFabricService(cfg)
	meterSvc := services.NewMeterService(fabricSvc)
	billingSvc := services.NewBillingService(fabricSvc)
	alertSvc := services.NewAlertService()

	// Initialize handlers
	meterHandler := handlers.NewMeterHandler(meterSvc)
	billingHandler := handlers.NewBillingHandler(billingSvc)
	consumerHandler := handlers.NewConsumerHandler()
	tariffHandler := handlers.NewTariffHandler()
	disputeHandler := handlers.NewDisputeHandler()
	dashboardHandler := handlers.NewDashboardHandler(alertSvc)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:   "SmartMeterChain API v1.0",
		BodyLimit: 10 * 1024 * 1024, // 10MB
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "SmartMeterChain API",
			"version": "1.0.0",
		})
	})

	// ─── Public Routes ─────────────────────────────────────
	api := app.Group("/api/v1")

	// Auth
	api.Post("/auth/login", func(c *fiber.Ctx) error {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&req); err != nil {
			return utils.Error(c, fiber.StatusBadRequest, "Invalid request body")
		}
		var user models.User
		if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
			return utils.Error(c, fiber.StatusUnauthorized, "Invalid credentials")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return utils.Error(c, fiber.StatusUnauthorized, "Invalid credentials")
		}
		token, err := middleware.GenerateToken(&user, cfg)
		if err != nil {
			return utils.Error(c, fiber.StatusInternalServerError, "Failed to generate token")
		}
		return utils.Success(c, fiber.Map{
			"token":    token,
			"user_id":  user.ID,
			"username": user.Username,
			"role":     user.Role,
			"name":     user.Name,
		})
	})

	api.Post("/auth/register", consumerHandler.Register)

	// Bill verification (public)
	api.Get("/verify/bill/:id", billingHandler.VerifyBill)

	// ─── Protected Routes ──────────────────────────────────
	auth := api.Group("", middleware.AuthRequired(cfg))

	// Profile
	auth.Get("/profile", consumerHandler.GetProfile)

	// Meters
	auth.Get("/meters", meterHandler.GetMeters)
	auth.Get("/meters/:id", meterHandler.GetMeter)
	auth.Get("/meters/:id/readings", meterHandler.GetReadings)
	auth.Post("/meters", middleware.RoleRequired("admin"), meterHandler.RegisterMeter)
	auth.Post("/meters/readings", meterHandler.IngestReading)
	auth.Post("/meters/readings/bulk", middleware.RoleRequired("admin"), meterHandler.BulkIngest)

	// Billing
	auth.Get("/bills", billingHandler.GetBills)
	auth.Get("/bills/:id", billingHandler.GetBill)
	auth.Post("/bills/generate", middleware.RoleRequired("admin"), billingHandler.GenerateBill)
	auth.Post("/bills/:id/pay", billingHandler.PayBill)

	// Tariffs
	auth.Get("/tariffs", tariffHandler.GetTariffs)
	auth.Get("/tariffs/:id", tariffHandler.GetTariff)
	auth.Post("/tariffs", middleware.RoleRequired("admin", "regulator"), tariffHandler.CreateTariff)
	auth.Put("/tariffs/:id", middleware.RoleRequired("admin", "regulator"), tariffHandler.UpdateTariff)
	auth.Delete("/tariffs/:id", middleware.RoleRequired("admin"), tariffHandler.DeleteTariff)

	// Disputes
	auth.Get("/disputes", disputeHandler.GetDisputes)
	auth.Get("/disputes/:id", disputeHandler.GetDispute)
	auth.Post("/disputes", disputeHandler.FileDispute)
	auth.Put("/disputes/:id/resolve", middleware.RoleRequired("admin", "regulator"), disputeHandler.ResolveDispute)

	// Dashboard & Analytics
	auth.Get("/dashboard/stats", middleware.RoleRequired("admin", "regulator"), dashboardHandler.GetDashboardStats)
	auth.Get("/dashboard/consumption", dashboardHandler.GetConsumptionTrend)
	auth.Get("/dashboard/alerts", middleware.RoleRequired("admin", "regulator"), dashboardHandler.GetAlerts)
	auth.Post("/dashboard/alerts/ack", middleware.RoleRequired("admin"), dashboardHandler.AcknowledgeAlert)

	// Consumers (admin only)
	auth.Get("/consumers", middleware.RoleRequired("admin", "regulator"), consumerHandler.GetConsumers)

	log.Printf("SmartMeterChain API starting on :%s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
