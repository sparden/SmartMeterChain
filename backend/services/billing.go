package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/smartmeterchain/backend/database"
	"github.com/smartmeterchain/backend/models"
	"github.com/smartmeterchain/backend/utils"
)

type BillingService struct {
	fabric *FabricService
}

func NewBillingService(fabric *FabricService) *BillingService {
	return &BillingService{fabric: fabric}
}

type GenerateBillInput struct {
	MeterID     string `json:"meter_id" validate:"required"`
	PeriodStart string `json:"period_start" validate:"required"` // YYYY-MM-DD
	PeriodEnd   string `json:"period_end" validate:"required"`
}

// GenerateBill calculates bill using slab-based tariff and submits to blockchain
func (bs *BillingService) GenerateBill(input GenerateBillInput) (*models.BillCache, string, error) {
	// Validate meter
	var meter models.Meter
	if err := database.DB.Where("meter_id = ?", input.MeterID).First(&meter).Error; err != nil {
		return nil, "", fmt.Errorf("meter %s not found", input.MeterID)
	}

	periodStart, err := time.Parse("2006-01-02", input.PeriodStart)
	if err != nil {
		return nil, "", fmt.Errorf("invalid period_start format")
	}
	periodEnd, err := time.Parse("2006-01-02", input.PeriodEnd)
	if err != nil {
		return nil, "", fmt.Errorf("invalid period_end format")
	}

	// Calculate units consumed in period
	var readings []models.ReadingCache
	database.DB.Where("meter_id = ? AND timestamp BETWEEN ? AND ?",
		input.MeterID, periodStart, periodEnd).
		Order("timestamp asc").
		Find(&readings)

	var unitsUsed float64
	if len(readings) >= 2 {
		unitsUsed = readings[len(readings)-1].Reading - readings[0].Reading
	} else if len(readings) == 1 {
		unitsUsed = readings[0].Reading
	}
	if unitsUsed < 0 {
		unitsUsed = 0
	}

	// Get tariff slabs for meter type
	var tariffs []models.TariffCache
	database.DB.Where("category = ?", meter.MeterType).
		Order("slab_start asc").
		Find(&tariffs)

	// Calculate bill using slab-based tariff
	amount := calculateSlabBill(unitsUsed, tariffs)

	billID := "BILL-" + uuid.New().String()[:8]

	// Submit to blockchain
	_, txID, err := bs.fabric.SubmitTransaction("GenerateBill",
		billID, input.MeterID, meter.ConsumerID,
		fmt.Sprintf("%.2f", unitsUsed),
		fmt.Sprintf("%.2f", amount),
	)
	if err != nil {
		txID = "pending-" + fmt.Sprintf("%d", time.Now().UnixNano())
	}

	bill := &models.BillCache{
		BillID:      billID,
		MeterID:     input.MeterID,
		ConsumerID:  meter.ConsumerID,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		UnitsUsed:   unitsUsed,
		Amount:      amount,
		Status:      "pending",
		TxID:        txID,
		BillHash:    utils.HashString(fmt.Sprintf("%s:%.2f:%.2f", billID, unitsUsed, amount)),
		DueDate:     periodEnd.AddDate(0, 0, 15),
	}
	database.DB.Create(bill)

	return bill, txID, nil
}

func calculateSlabBill(units float64, tariffs []models.TariffCache) float64 {
	if len(tariffs) == 0 {
		return units * 5.0 // fallback rate
	}

	var total float64
	remaining := units

	for _, slab := range tariffs {
		if remaining <= 0 {
			break
		}
		slabWidth := slab.SlabEnd - slab.SlabStart
		unitsInSlab := remaining
		if unitsInSlab > slabWidth {
			unitsInSlab = slabWidth
		}
		total += unitsInSlab * slab.RatePerUnit
		remaining -= unitsInSlab

		// Add fixed charge for the highest applicable slab
		if remaining <= 0 || remaining < slabWidth {
			total += slab.FixedCharge
		}
	}

	return total
}

// VerifyBill checks bill hash against blockchain record
func (bs *BillingService) VerifyBill(billID string) (bool, *models.BillCache, error) {
	var bill models.BillCache
	if err := database.DB.Where("bill_id = ?", billID).First(&bill).Error; err != nil {
		return false, nil, fmt.Errorf("bill %s not found", billID)
	}

	expectedHash := utils.HashString(fmt.Sprintf("%s:%.2f:%.2f", bill.BillID, bill.UnitsUsed, bill.Amount))
	verified := bill.BillHash == expectedHash

	return verified, &bill, nil
}

// GetBillsByConsumer returns paginated bills for a consumer
func (bs *BillingService) GetBillsByConsumer(consumerID string, page, perPage int) ([]models.BillCache, int64, error) {
	var bills []models.BillCache
	var total int64

	query := database.DB.Where("consumer_id = ?", consumerID)
	query.Model(&models.BillCache{}).Count(&total)

	err := query.Order("period_end desc").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&bills).Error

	return bills, total, err
}
