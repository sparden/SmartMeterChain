package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/smartmeterchain/backend/database"
	"github.com/smartmeterchain/backend/models"
	"github.com/smartmeterchain/backend/services"
	"github.com/smartmeterchain/backend/utils"
)

type DashboardHandler struct {
	alertSvc *services.AlertService
}

func NewDashboardHandler(alertSvc *services.AlertService) *DashboardHandler {
	return &DashboardHandler{alertSvc: alertSvc}
}

// GetDashboardStats returns overview statistics
func (h *DashboardHandler) GetDashboardStats(c *fiber.Ctx) error {
	var totalMeters, activeMeters int64
	var totalBills int64
	var totalRevenue float64
	var openDisputes int64
	var todayReadings int64

	database.DB.Model(&models.Meter{}).Count(&totalMeters)
	database.DB.Model(&models.Meter{}).Where("status = ?", "active").Count(&activeMeters)
	database.DB.Model(&models.BillCache{}).Count(&totalBills)
	database.DB.Model(&models.BillCache{}).Where("status = ?", "paid").Select("COALESCE(SUM(amount), 0)").Scan(&totalRevenue)
	database.DB.Model(&models.DisputeCache{}).Where("status = ?", "open").Count(&openDisputes)

	today := time.Now().Truncate(24 * time.Hour)
	database.DB.Model(&models.ReadingCache{}).Where("timestamp >= ?", today).Count(&todayReadings)

	alertStats := h.alertSvc.GetAlertStats()

	return utils.Success(c, fiber.Map{
		"meters": fiber.Map{
			"total":  totalMeters,
			"active": activeMeters,
		},
		"billing": fiber.Map{
			"total_bills":   totalBills,
			"total_revenue": totalRevenue,
		},
		"disputes": fiber.Map{
			"open": openDisputes,
		},
		"readings": fiber.Map{
			"today": todayReadings,
		},
		"alerts": alertStats,
	})
}

// GetConsumptionTrend returns daily consumption over the last N days
func (h *DashboardHandler) GetConsumptionTrend(c *fiber.Ctx) error {
	days := c.QueryInt("days", 30)

	type DailyConsumption struct {
		Date     string  `json:"date"`
		Readings int64   `json:"readings"`
		AvgValue float64 `json:"avg_value"`
	}

	var trend []DailyConsumption
	startDate := time.Now().AddDate(0, 0, -days)

	database.DB.Model(&models.ReadingCache{}).
		Select("DATE(timestamp) as date, COUNT(*) as readings, AVG(reading) as avg_value").
		Where("timestamp >= ?", startDate).
		Group("DATE(timestamp)").
		Order("date asc").
		Scan(&trend)

	return utils.Success(c, trend)
}

// GetAlerts returns tamper alerts
func (h *DashboardHandler) GetAlerts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 20)
	severity := c.Query("severity")

	alerts, total, err := h.alertSvc.GetAlerts(page, perPage, severity)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.Paginated(c, alerts, total, page, perPage)
}

// AcknowledgeAlert marks an alert as acknowledged
func (h *DashboardHandler) AcknowledgeAlert(c *fiber.Ctx) error {
	alertID := c.QueryInt("id", 0)
	if alertID == 0 {
		return utils.Error(c, fiber.StatusBadRequest, "Alert ID required")
	}
	if err := h.alertSvc.AcknowledgeAlert(uint(alertID)); err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Alert not found")
	}
	return utils.SuccessWithMessage(c, "Alert acknowledged", nil)
}
