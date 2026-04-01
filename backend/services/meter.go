package services

import (
	"fmt"
	"math"
	"time"

	"github.com/smartmeterchain/backend/database"
	"github.com/smartmeterchain/backend/models"
)

type MeterService struct {
	fabric *FabricService
}

func NewMeterService(fabric *FabricService) *MeterService {
	return &MeterService{fabric: fabric}
}

type ReadingInput struct {
	MeterID     string  `json:"meter_id" validate:"required"`
	Reading     float64 `json:"reading" validate:"required"`
	Timestamp   int64   `json:"timestamp"`
	Temperature float64 `json:"temperature"`
	Voltage     float64 `json:"voltage"`
}

// IngestReading processes a new meter reading with tamper detection
func (ms *MeterService) IngestReading(input ReadingInput) (*models.ReadingCache, string, error) {
	// Validate meter exists
	var meter models.Meter
	if err := database.DB.Where("meter_id = ?", input.MeterID).First(&meter).Error; err != nil {
		return nil, "", fmt.Errorf("meter %s not found", input.MeterID)
	}

	ts := time.Now()
	if input.Timestamp > 0 {
		ts = time.Unix(input.Timestamp, 0)
	}

	// Tamper detection
	anomaly := ms.detectAnomaly(&meter, input.Reading)

	// Submit to blockchain
	_, txID, err := ms.fabric.SubmitTransaction("SubmitReading",
		input.MeterID,
		fmt.Sprintf("%.4f", input.Reading),
		ts.Format(time.RFC3339),
	)
	if err != nil {
		txID = "pending-" + fmt.Sprintf("%d", time.Now().UnixNano())
	}

	// Cache reading off-chain
	reading := &models.ReadingCache{
		MeterID:     input.MeterID,
		Reading:     input.Reading,
		Timestamp:   ts,
		TxID:        txID,
		IsAnomaly:   anomaly,
		Temperature: input.Temperature,
		Voltage:     input.Voltage,
	}
	database.DB.Create(reading)

	// Update meter's last reading
	database.DB.Model(&meter).Updates(map[string]interface{}{
		"last_reading":  input.Reading,
		"last_synced_at": ts,
	})

	// Create tamper alert if anomaly detected
	if anomaly {
		ms.createTamperAlert(&meter, input.Reading)
	}

	return reading, txID, nil
}

func (ms *MeterService) detectAnomaly(meter *models.Meter, newReading float64) bool {
	if meter.LastReading == 0 {
		return false
	}

	// Reading should not decrease (reverse flow)
	if newReading < meter.LastReading {
		return true
	}

	// Spike detection: >500% increase is suspicious
	if meter.LastReading > 0 {
		ratio := newReading / meter.LastReading
		if ratio > 5.0 {
			return true
		}
	}

	// Sudden large jump (>1000 kWh in single reading)
	if math.Abs(newReading-meter.LastReading) > 1000 {
		return true
	}

	return false
}

func (ms *MeterService) createTamperAlert(meter *models.Meter, reading float64) {
	alertType := "spike"
	if reading < meter.LastReading {
		alertType = "reverse_flow"
	}

	alert := &models.TamperAlert{
		MeterID:      meter.MeterID,
		AlertType:    alertType,
		Description:  fmt.Sprintf("Anomaly detected: previous=%.2f, current=%.2f", meter.LastReading, reading),
		ReadingValue: reading,
		ExpectedMin:  meter.LastReading * 0.8,
		ExpectedMax:  meter.LastReading * 2.0,
		Severity:     "high",
		DetectedAt:   time.Now(),
	}
	database.DB.Create(alert)
}

// GetReadings returns paginated readings for a meter
func (ms *MeterService) GetReadings(meterID string, page, perPage int) ([]models.ReadingCache, int64, error) {
	var readings []models.ReadingCache
	var total int64

	query := database.DB.Where("meter_id = ?", meterID)
	query.Model(&models.ReadingCache{}).Count(&total)

	err := query.Order("timestamp desc").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&readings).Error

	return readings, total, err
}
