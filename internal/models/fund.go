package models

import (
	"time"

	"gorm.io/gorm"
)

type Fund struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" binding:"required"`
	Description   string    `json:"description"`
	RiskLevel     string    `json:"risk_level"`
	MinInvestment float64   `json:"min_investment"`
	MaxInvestment float64   `json:"max_investment"`
	Performance   []Performance `json:"performance" gorm:"foreignKey:FundID"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Performance struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FundID    string    `json:"fund_id"`
	Date      time.Time `json:"date"`
	Value     float64   `json:"value"`
	Change    float64   `json:"change"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FundRepository struct {
	db *gorm.DB
}

func NewFundRepository(db *gorm.DB) *FundRepository {
	return &FundRepository{db: db}
}

func (r *FundRepository) List() ([]Fund, error) {
	var funds []Fund
	err := r.db.Preload("Performance").Find(&funds).Error
	return funds, err
}

func (r *FundRepository) GetByID(id string) (*Fund, error) {
	var fund Fund
	err := r.db.Preload("Performance").First(&fund, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &fund, nil
}

func (r *FundRepository) Create(fund *Fund) error {
	return r.db.Create(fund).Error
}

func (r *FundRepository) AddPerformance(performance *Performance) error {
	return r.db.Create(performance).Error
} 