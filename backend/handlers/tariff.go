package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/smartmeterchain/backend/database"
	"github.com/smartmeterchain/backend/models"
	"github.com/smartmeterchain/backend/utils"
)

type TariffHandler struct{}

func NewTariffHandler() *TariffHandler {
	return &TariffHandler{}
}

// GetTariffs returns all tariff slabs
func (h *TariffHandler) GetTariffs(c *fiber.Ctx) error {
	category := c.Query("category")
	var tariffs []models.TariffCache

	query := database.DB.Model(&models.TariffCache{})
	if category != "" {
		query = query.Where("category = ?", category)
	}
	query.Order("category asc, slab_start asc").Find(&tariffs)

	return utils.Success(c, tariffs)
}

// GetTariff returns a single tariff by ID
func (h *TariffHandler) GetTariff(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var tariff models.TariffCache
	if err := database.DB.First(&tariff, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Tariff not found")
	}
	return utils.Success(c, tariff)
}

// CreateTariff adds a new tariff slab
func (h *TariffHandler) CreateTariff(c *fiber.Ctx) error {
	var tariff models.TariffCache
	if err := c.BodyParser(&tariff); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if tariff.TariffID == "" || tariff.Category == "" {
		return utils.Error(c, fiber.StatusBadRequest, "tariff_id and category are required")
	}

	if err := database.DB.Create(&tariff).Error; err != nil {
		return utils.Error(c, fiber.StatusConflict, "Tariff ID already exists")
	}
	return utils.Created(c, tariff)
}

// UpdateTariff modifies an existing tariff
func (h *TariffHandler) UpdateTariff(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var existing models.TariffCache
	if err := database.DB.First(&existing, id).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Tariff not found")
	}

	var updates models.TariffCache
	if err := c.BodyParser(&updates); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	database.DB.Model(&existing).Updates(updates)
	return utils.Success(c, existing)
}

// DeleteTariff removes a tariff slab
func (h *TariffHandler) DeleteTariff(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	result := database.DB.Delete(&models.TariffCache{}, id)
	if result.RowsAffected == 0 {
		return utils.Error(c, fiber.StatusNotFound, "Tariff not found")
	}
	return utils.SuccessWithMessage(c, "Tariff deleted", nil)
}
