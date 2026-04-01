package main

// MeterReading represents a single smart meter consumption reading
type MeterReading struct {
	ID           string  `json:"id"`
	MeterID      string  `json:"meterId"`
	ConsumerID   string  `json:"consumerId"`
	Reading      float64 `json:"reading"`       // Cumulative kWh
	Timestamp    string  `json:"timestamp"`      // ISO 8601
	Hash         string  `json:"hash"`           // SHA-256 of meterID+reading+timestamp
	PreviousHash string  `json:"previousHash"`
	Status       string  `json:"status"`         // valid, suspect, tampered
	DocType      string  `json:"docType"`
}

// Meter represents a registered smart meter device
type Meter struct {
	ID              string `json:"id"`
	MeterID         string `json:"meterId"`
	ConsumerID      string `json:"consumerId"`
	Location        string `json:"location"`
	Type            string `json:"type"`            // domestic, commercial, industrial
	Status          string `json:"status"`          // active, inactive, tampered
	InstalledDate   string `json:"installedDate"`
	LastReadingTime string `json:"lastReadingTime"`
	LastReadingHash string `json:"lastReadingHash"`
	DocType         string `json:"docType"`
}

// Consumer represents an electricity consumer
type Consumer struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ConsumerNumber string `json:"consumerNumber"`
	MeterID        string `json:"meterId"`
	Category       string `json:"category"` // domestic, commercial, industrial
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	DocType        string `json:"docType"`
}

// Bill represents a generated electricity bill
type Bill struct {
	ID                 string  `json:"id"`
	BillNumber         string  `json:"billNumber"`
	ConsumerID         string  `json:"consumerId"`
	MeterID            string  `json:"meterId"`
	BillingPeriodStart string  `json:"billingPeriodStart"`
	BillingPeriodEnd   string  `json:"billingPeriodEnd"`
	UnitsConsumed      float64 `json:"unitsConsumed"`
	SlabBreakdown      []SlabCharge `json:"slabBreakdown"`
	FixedCharge        float64 `json:"fixedCharge"`
	TotalAmount        float64 `json:"totalAmount"`
	Hash               string  `json:"hash"`
	Status             string  `json:"status"` // generated, paid, disputed, overdue
	GeneratedAt        string  `json:"generatedAt"`
	PaidAt             string  `json:"paidAt,omitempty"`
	DocType            string  `json:"docType"`
}

// SlabCharge represents the charge for a specific tariff slab
type SlabCharge struct {
	SlabName    string  `json:"slabName"`
	Units       float64 `json:"units"`
	RatePerUnit float64 `json:"ratePerUnit"`
	Amount      float64 `json:"amount"`
}

// Tariff represents a tariff structure for a consumer category
type Tariff struct {
	ID            string       `json:"id"`
	Category      string       `json:"category"` // domestic, commercial, industrial
	Slabs         []TariffSlab `json:"slabs"`
	FixedCharge   float64      `json:"fixedCharge"`
	EffectiveFrom string       `json:"effectiveFrom"`
	EffectiveTo   string       `json:"effectiveTo,omitempty"`
	DocType       string       `json:"docType"`
}

// TariffSlab represents a single slab in a tariff
type TariffSlab struct {
	MinUnits    float64 `json:"minUnits"`
	MaxUnits    float64 `json:"maxUnits"`    // -1 for unlimited
	RatePerUnit float64 `json:"ratePerUnit"` // INR per kWh
}

// Dispute represents a billing dispute
type Dispute struct {
	ID         string `json:"id"`
	BillID     string `json:"billId"`
	ConsumerID string `json:"consumerId"`
	Reason     string `json:"reason"`
	Status     string `json:"status"` // open, investigating, resolved, rejected
	FiledAt    string `json:"filedAt"`
	ResolvedAt string `json:"resolvedAt,omitempty"`
	Resolution string `json:"resolution,omitempty"`
	DocType    string `json:"docType"`
}

// TamperAlert represents a detected tamper event
type TamperAlert struct {
	ID           string `json:"id"`
	MeterID      string `json:"meterId"`
	ExpectedHash string `json:"expectedHash"`
	ActualHash   string `json:"actualHash"`
	DetectedAt   string `json:"detectedAt"`
	Severity     string `json:"severity"` // low, medium, high, critical
	DocType      string `json:"docType"`
}
