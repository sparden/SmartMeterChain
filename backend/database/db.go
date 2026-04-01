package database

import (
	"log"
	"time"

	"github.com/smartmeterchain/backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(dbPath string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate all models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Meter{},
		&models.ReadingCache{},
		&models.BillCache{},
		&models.TariffCache{},
		&models.DisputeCache{},
		&models.TamperAlert{},
		&models.AuditLog{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	seed()
	log.Println("Database connected and migrated successfully")
}

func seed() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	log.Println("Seeding database with initial data...")

	// Create default admin
	hashedPw, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		Username: "admin",
		Email:    "admin@discom.gov.in",
		Password: string(hashedPw),
		Role:     "admin",
		Name:     "DISCOM Administrator",
	}
	DB.Create(&admin)

	// Create sample consumer
	consumerPw, _ := bcrypt.GenerateFromPassword([]byte("consumer123"), bcrypt.DefaultCost)
	consumer := models.User{
		Username: "consumer1",
		Email:    "consumer1@example.com",
		Password: string(consumerPw),
		Role:     "consumer",
		Name:     "Rajesh Kumar",
		Phone:    "+919876543210",
	}
	DB.Create(&consumer)

	// Create regulator
	regPw, _ := bcrypt.GenerateFromPassword([]byte("regulator123"), bcrypt.DefaultCost)
	regulator := models.User{
		Username: "regulator",
		Email:    "regulator@serc.gov.in",
		Password: string(regPw),
		Role:     "regulator",
		Name:     "State Electricity Regulatory Commission",
	}
	DB.Create(&regulator)

	// Create sample meters
	meters := []models.Meter{
		{MeterID: "SM-DEL-001", ConsumerID: "consumer1", Location: "New Delhi, Sector 15", MeterType: "domestic", Status: "active", InstallDate: time.Now().AddDate(0, -6, 0)},
		{MeterID: "SM-DEL-002", ConsumerID: "consumer1", Location: "New Delhi, Sector 22", MeterType: "domestic", Status: "active", InstallDate: time.Now().AddDate(0, -3, 0)},
		{MeterID: "SM-MUM-001", ConsumerID: "commercial1", Location: "Mumbai, BKC", MeterType: "commercial", Status: "active", InstallDate: time.Now().AddDate(-1, 0, 0)},
		{MeterID: "SM-BLR-001", ConsumerID: "industrial1", Location: "Bangalore, Electronic City", MeterType: "industrial", Status: "active", InstallDate: time.Now().AddDate(-2, 0, 0)},
	}
	DB.Create(&meters)

	// Create default tariff slabs (domestic)
	tariffs := []models.TariffCache{
		{TariffID: "TAR-DOM-001", Category: "domestic", SlabStart: 0, SlabEnd: 100, RatePerUnit: 3.0, FixedCharge: 50, EffectiveAt: time.Now()},
		{TariffID: "TAR-DOM-002", Category: "domestic", SlabStart: 100, SlabEnd: 300, RatePerUnit: 5.5, FixedCharge: 100, EffectiveAt: time.Now()},
		{TariffID: "TAR-DOM-003", Category: "domestic", SlabStart: 300, SlabEnd: 99999, RatePerUnit: 8.0, FixedCharge: 150, EffectiveAt: time.Now()},
		{TariffID: "TAR-COM-001", Category: "commercial", SlabStart: 0, SlabEnd: 99999, RatePerUnit: 9.5, FixedCharge: 500, EffectiveAt: time.Now()},
		{TariffID: "TAR-IND-001", Category: "industrial", SlabStart: 0, SlabEnd: 99999, RatePerUnit: 7.0, FixedCharge: 1000, EffectiveAt: time.Now()},
	}
	DB.Create(&tariffs)

	log.Println("Seeding complete")
}
