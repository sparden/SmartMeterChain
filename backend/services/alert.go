package services

import (
	"github.com/smartmeterchain/backend/database"
	"github.com/smartmeterchain/backend/models"
)

type AlertService struct{}

func NewAlertService() *AlertService {
	return &AlertService{}
}

// GetAlerts returns paginated tamper alerts
func (as *AlertService) GetAlerts(page, perPage int, severity string) ([]models.TamperAlert, int64, error) {
	var alerts []models.TamperAlert
	var total int64

	query := database.DB.Model(&models.TamperAlert{})
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	query.Count(&total)

	err := query.Order("detected_at desc").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&alerts).Error

	return alerts, total, err
}

// GetAlertsByMeter returns alerts for a specific meter
func (as *AlertService) GetAlertsByMeter(meterID string) ([]models.TamperAlert, error) {
	var alerts []models.TamperAlert
	err := database.DB.Where("meter_id = ?", meterID).
		Order("detected_at desc").
		Limit(50).
		Find(&alerts).Error
	return alerts, err
}

// AcknowledgeAlert marks an alert as acknowledged
func (as *AlertService) AcknowledgeAlert(alertID uint) error {
	return database.DB.Model(&models.TamperAlert{}).
		Where("id = ?", alertID).
		Update("acknowledged", true).Error
}

// GetAlertStats returns summary statistics
func (as *AlertService) GetAlertStats() map[string]interface{} {
	var totalAlerts, unacknowledged, critical int64

	database.DB.Model(&models.TamperAlert{}).Count(&totalAlerts)
	database.DB.Model(&models.TamperAlert{}).Where("acknowledged = ?", false).Count(&unacknowledged)
	database.DB.Model(&models.TamperAlert{}).Where("severity = ?", "critical").Count(&critical)

	return map[string]interface{}{
		"total_alerts":   totalAlerts,
		"unacknowledged": unacknowledged,
		"critical":       critical,
	}
}
