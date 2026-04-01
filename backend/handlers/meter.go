package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/smartmeterchain/backend/database"
	"github.com/smartmeterchain/backend/models"
	"github.com/smartmeterchain/backend/services"
	"github.com/smartmeterchain/backend/utils"
)

type MeterHandler struct {
	meterSvc *services.MeterService
}

func NewMeterHandler(meterSvc *services.MeterService) *MeterHandler {
	return &MeterHandler{meterSvc: meterSvc}
}

// RegisterMeter creates a new meter in the system
func (h *MeterHandler) RegisterMeter(c *fiber.Ctx) error {
	var meter models.Meter
	if err := c.BodyParser(&meter); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if meter.MeterID == "" || meter.ConsumerID == "" {
		return utils.Error(c, fiber.StatusBadRequest, "meter_id and consumer_id are required")
	}

	if err := database.DB.Create(&meter).Error; err != nil {
		return utils.Error(c, fiber.StatusConflict, "Meter already exists or DB error")
	}
	return utils.Created(c, meter)
}

// GetMeters returns all meters with optional filters
func (h *MeterHandler) GetMeters(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))
	status := c.Query("status")
	meterType := c.Query("type")

	var meters []models.Meter
	var total int64

	query := database.DB.Model(&models.Meter{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if meterType != "" {
		query = query.Where("meter_type = ?", meterType)
	}
	query.Count(&total)

	query.Order("created_at desc").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&meters)

	return utils.Paginated(c, meters, total, page, perPage)
}

// GetMeter returns a single meter by ID
func (h *MeterHandler) GetMeter(c *fiber.Ctx) error {
	meterID := c.Params("id")
	var meter models.Meter
	if err := database.DB.Where("meter_id = ?", meterID).First(&meter).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Meter not found")
	}
	return utils.Success(c, meter)
}

// IngestReading submits a new meter reading
func (h *MeterHandler) IngestReading(c *fiber.Ctx) error {
	var input services.ReadingInput
	if err := c.BodyParser(&input); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if input.MeterID == "" {
		return utils.Error(c, fiber.StatusBadRequest, "meter_id is required")
	}

	reading, txID, err := h.meterSvc.IngestReading(input)
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessWithTx(c, "Reading ingested successfully", reading, txID)
}

// GetReadings returns readings for a meter
func (h *MeterHandler) GetReadings(c *fiber.Ctx) error {
	meterID := c.Params("id")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "50"))

	readings, total, err := h.meterSvc.GetReadings(meterID, page, perPage)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.Paginated(c, readings, total, page, perPage)
}

// BulkIngest handles batch meter reading submissions
func (h *MeterHandler) BulkIngest(c *fiber.Ctx) error {
	var inputs []services.ReadingInput
	if err := c.BodyParser(&inputs); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body — expected array of readings")
	}

	var results []map[string]interface{}
	var errors []string

	for _, input := range inputs {
		reading, txID, err := h.meterSvc.IngestReading(input)
		if err != nil {
			errors = append(errors, input.MeterID+": "+err.Error())
			continue
		}
		results = append(results, map[string]interface{}{
			"meter_id":   input.MeterID,
			"reading_id": reading.ID,
			"tx_id":      txID,
		})
	}

	return utils.Success(c, fiber.Map{
		"ingested": len(results),
		"failed":   len(errors),
		"results":  results,
		"errors":   errors,
	})
}
