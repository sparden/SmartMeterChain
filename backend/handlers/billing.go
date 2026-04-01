package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/smartmeterchain/backend/database"
	"github.com/smartmeterchain/backend/models"
	"github.com/smartmeterchain/backend/services"
	"github.com/smartmeterchain/backend/utils"
)

type BillingHandler struct {
	billingSvc *services.BillingService
}

func NewBillingHandler(billingSvc *services.BillingService) *BillingHandler {
	return &BillingHandler{billingSvc: billingSvc}
}

// GenerateBill creates a new bill for a meter
func (h *BillingHandler) GenerateBill(c *fiber.Ctx) error {
	var input services.GenerateBillInput
	if err := c.BodyParser(&input); err != nil {
		return utils.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	bill, txID, err := h.billingSvc.GenerateBill(input)
	if err != nil {
		return utils.Error(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessWithTx(c, "Bill generated successfully", bill, txID)
}

// GetBills returns bills for a consumer
func (h *BillingHandler) GetBills(c *fiber.Ctx) error {
	consumerID := c.Query("consumer_id")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	if consumerID == "" {
		// Return all bills for admin
		var bills []models.BillCache
		var total int64
		database.DB.Model(&models.BillCache{}).Count(&total)
		database.DB.Order("period_end desc").
			Offset((page - 1) * perPage).
			Limit(perPage).
			Find(&bills)
		return utils.Paginated(c, bills, total, page, perPage)
	}

	bills, total, err := h.billingSvc.GetBillsByConsumer(consumerID, page, perPage)
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.Paginated(c, bills, total, page, perPage)
}

// GetBill returns a single bill
func (h *BillingHandler) GetBill(c *fiber.Ctx) error {
	billID := c.Params("id")
	var bill models.BillCache
	if err := database.DB.Where("bill_id = ?", billID).First(&bill).Error; err != nil {
		return utils.Error(c, fiber.StatusNotFound, "Bill not found")
	}
	return utils.Success(c, bill)
}

// VerifyBill checks bill integrity against blockchain
func (h *BillingHandler) VerifyBill(c *fiber.Ctx) error {
	billID := c.Params("id")
	verified, bill, err := h.billingSvc.VerifyBill(billID)
	if err != nil {
		return utils.Error(c, fiber.StatusNotFound, err.Error())
	}

	return utils.Success(c, fiber.Map{
		"bill_id":  billID,
		"verified": verified,
		"bill":     bill,
	})
}

// PayBill marks a bill as paid
func (h *BillingHandler) PayBill(c *fiber.Ctx) error {
	billID := c.Params("id")
	result := database.DB.Model(&models.BillCache{}).
		Where("bill_id = ? AND status = ?", billID, "pending").
		Update("status", "paid")

	if result.RowsAffected == 0 {
		return utils.Error(c, fiber.StatusNotFound, "Bill not found or already paid")
	}
	return utils.SuccessWithMessage(c, "Bill marked as paid", fiber.Map{"bill_id": billID})
}
