package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/smartmeterchain/backend/database"
	"github.com/smartmeterchain/backend/models"
	"github.com/smartmeterchain/backend/utils"
)

type DisputeHandler struct{}

func NewDisputeHandler() *DisputeHandler {
	return &DisputeHandler{}
}

type FileDisputeRequest struct {
	BillID string `json:"bill_id" validate:"required"`
	Reason string `json:"reason" validate:"required"`
}

// FileDispute creates a new billing dispute
func (h *DisputeHandler) FileDispute(c *fiber.Ctx) error {
	var req FileDisputeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if req.BillID == "" || req.Reason == "" {
		return utils.Error(c, fiber.StatusBadRequest, "bill_id and reason are required")
	}

	// Verify bill exists
	var bill models.BillCache
	if err := database.DB.Where("bill_id = ?", req.BillID).First(&bill).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Bill not found")
	}

	dispute := models.DisputeCache{
		DisputeID:  "DSP-" + uuid.New().String()[:8],
		BillID:     req.BillID,
		ConsumerID: bill.ConsumerID,
		Reason:     req.Reason,
		Status:     "open",
		FiledAt:    time.Now(),
	}

	database.DB.Create(&dispute)

	// Update bill status
	database.DB.Model(&bill).Update("status", "disputed")

	return utils.Created(c, dispute)
}

// GetDisputes returns all disputes with optional filters
func (h *DisputeHandler) GetDisputes(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))
	status := c.Query("status")
	consumerID := c.Query("consumer_id")

	var disputes []models.DisputeCache
	var total int64

	query := database.DB.Model(&models.DisputeCache{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if consumerID != "" {
		query = query.Where("consumer_id = ?", consumerID)
	}
	query.Count(&total)

	query.Order("filed_at desc").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&disputes)

	return utils.Paginated(c, disputes, total, page, perPage)
}

// GetDispute returns a single dispute
func (h *DisputeHandler) GetDispute(c *fiber.Ctx) error {
	disputeID := c.Params("id")
	var dispute models.DisputeCache
	if err := database.DB.Where("dispute_id = ?", disputeID).First(&dispute).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Dispute not found")
	}
	return utils.Success(c, dispute)
}

type ResolveDisputeRequest struct {
	Resolution string `json:"resolution" validate:"required"`
	Status     string `json:"status"` // resolved or rejected
}

// ResolveDispute updates a dispute with resolution
func (h *DisputeHandler) ResolveDispute(c *fiber.Ctx) error {
	disputeID := c.Params("id")
	var req ResolveDisputeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	var dispute models.DisputeCache
	if err := database.DB.Where("dispute_id = ?", disputeID).First(&dispute).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Dispute not found")
	}

	if dispute.Status != "open" && dispute.Status != "investigating" {
		return utils.Error(c, fiber.StatusBadRequest, "Dispute is already resolved")
	}

	status := "resolved"
	if req.Status == "rejected" {
		status = "rejected"
	}

	now := time.Now()
	database.DB.Model(&dispute).Updates(map[string]interface{}{
		"status":      status,
		"resolution":  req.Resolution,
		"resolved_at": &now,
	})

	return utils.SuccessWithMessage(c, "Dispute "+status, dispute)
}
