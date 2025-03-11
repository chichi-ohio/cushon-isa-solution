package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Investment struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CustomerID uint      `json:"customer_id"`
	FundID     string    `json:"fund_id"`
	Amount     float64   `json:"amount" binding:"required,gt=0"`
	Status     string    `json:"status" gorm:"default:pending"`
	Customer   Customer  `json:"customer" gorm:"foreignKey:CustomerID"`
	Fund       Fund      `json:"fund" gorm:"foreignKey:FundID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type InvestmentRepository struct {
	db *gorm.DB
}

func NewInvestmentRepository(db *gorm.DB) *InvestmentRepository {
	return &InvestmentRepository{db: db}
}

func (r *InvestmentRepository) Create(inv *Investment) error {
	// Validate investment amount against fund limits
	var fund Fund
	if err := r.db.First(&fund, "id = ?", inv.FundID).Error; err != nil {
		return err
	}

	if inv.Amount < fund.MinInvestment || (fund.MaxInvestment > 0 && inv.Amount > fund.MaxInvestment) {
		return fmt.Errorf("investment amount outside fund limits (min: %.2f, max: %.2f)", fund.MinInvestment, fund.MaxInvestment)
	}

	return r.db.Create(inv).Error
}

func (r *InvestmentRepository) GetByID(id uint) (*Investment, error) {
	var investment Investment
	err := r.db.Preload("Customer").Preload("Fund").First(&investment, id).Error
	if err != nil {
		return nil, err
	}
	return &investment, nil
}

func (r *InvestmentRepository) List() ([]Investment, error) {
	var investments []Investment
	err := r.db.Preload("Customer").Preload("Fund").Order("created_at desc").Find(&investments).Error
	return investments, err
}

func (r *InvestmentRepository) ListByCustomer(customerID uint) ([]Investment, error) {
	var investments []Investment
	err := r.db.Preload("Fund").Where("customer_id = ?", customerID).Order("created_at desc").Find(&investments).Error
	return investments, err
}

func (r *InvestmentRepository) UpdateStatus(id uint, status string) error {
	result := r.db.Model(&Investment{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *InvestmentRepository) WithTx(tx *gorm.DB) *InvestmentRepository {
	return &InvestmentRepository{db: tx}
}
