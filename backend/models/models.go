package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents an API user (DISCOM admin, consumer, regulator)
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null" json:"username"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"not null;default:consumer" json:"role"` // admin, consumer, regulator
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

// Meter represents a smart meter registered in the system
type Meter struct {
	gorm.Model
	MeterID      string    `gorm:"uniqueIndex;not null" json:"meter_id"`
	ConsumerID   string    `gorm:"index;not null" json:"consumer_id"`
	Location     string    `json:"location"`
	MeterType    string    `json:"meter_type"` // domestic, commercial, industrial
	Status       string    `gorm:"default:active" json:"status"`
	InstallDate  time.Time `json:"install_date"`
	LastReading  float64   `json:"last_reading"`
	LastSyncedAt time.Time `json:"last_synced_at"`
	OnChain      bool      `gorm:"default:false" json:"on_chain"`
}

// ReadingCache caches meter readings off-chain for fast queries
type ReadingCache struct {
	gorm.Model
	MeterID     string    `gorm:"index;not null" json:"meter_id"`
	Reading     float64   `gorm:"not null" json:"reading"`
	Unit        string    `gorm:"default:kWh" json:"unit"`
	Timestamp   time.Time `gorm:"not null" json:"timestamp"`
	TxID        string    `json:"tx_id"`
	DataHash    string    `json:"data_hash"`
	IsAnomaly   bool      `gorm:"default:false" json:"is_anomaly"`
	Temperature float64   `json:"temperature"`
	Voltage     float64   `json:"voltage"`
}

// BillCache caches generated bills off-chain
type BillCache struct {
	gorm.Model
	BillID      string    `gorm:"uniqueIndex;not null" json:"bill_id"`
	MeterID     string    `gorm:"index;not null" json:"meter_id"`
	ConsumerID  string    `gorm:"index;not null" json:"consumer_id"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	UnitsUsed   float64   `json:"units_used"`
	Amount      float64   `json:"amount"`
	TariffID    string    `json:"tariff_id"`
	Status      string    `gorm:"default:pending" json:"status"` // pending, paid, disputed, overdue
	TxID        string    `json:"tx_id"`
	BillHash    string    `json:"bill_hash"`
	DueDate     time.Time `json:"due_date"`
}

// TariffCache caches tariff slabs off-chain
type TariffCache struct {
	gorm.Model
	TariffID    string  `gorm:"uniqueIndex;not null" json:"tariff_id"`
	Category    string  `gorm:"not null" json:"category"` // domestic, commercial, industrial
	SlabStart   float64 `json:"slab_start"`
	SlabEnd     float64 `json:"slab_end"`
	RatePerUnit float64 `gorm:"not null" json:"rate_per_unit"`
	FixedCharge float64 `json:"fixed_charge"`
	EffectiveAt time.Time `json:"effective_at"`
	OnChain     bool    `gorm:"default:false" json:"on_chain"`
}

// DisputeCache caches disputes off-chain
type DisputeCache struct {
	gorm.Model
	DisputeID   string    `gorm:"uniqueIndex;not null" json:"dispute_id"`
	BillID      string    `gorm:"index" json:"bill_id"`
	ConsumerID  string    `gorm:"index;not null" json:"consumer_id"`
	Reason      string    `gorm:"not null" json:"reason"`
	Status      string    `gorm:"default:open" json:"status"` // open, investigating, resolved, rejected
	Resolution  string    `json:"resolution"`
	FiledAt     time.Time `json:"filed_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
	TxID        string    `json:"tx_id"`
}

// TamperAlert stores tamper detection events
type TamperAlert struct {
	gorm.Model
	MeterID      string    `gorm:"index;not null" json:"meter_id"`
	AlertType    string    `gorm:"not null" json:"alert_type"` // spike, drop, voltage_anomaly, reverse_flow
	Description  string    `json:"description"`
	ReadingValue float64   `json:"reading_value"`
	ExpectedMin  float64   `json:"expected_min"`
	ExpectedMax  float64   `json:"expected_max"`
	Severity     string    `gorm:"default:medium" json:"severity"` // low, medium, high, critical
	Acknowledged bool      `gorm:"default:false" json:"acknowledged"`
	DetectedAt   time.Time `json:"detected_at"`
}

// AuditLog tracks all blockchain transactions
type AuditLog struct {
	gorm.Model
	Action    string `gorm:"not null" json:"action"`
	EntityID  string `gorm:"index" json:"entity_id"`
	UserID    uint   `gorm:"index" json:"user_id"`
	TxID      string `json:"tx_id"`
	Details   string `json:"details"`
	IPAddress string `json:"ip_address"`
}
