package db

import (
	"cushion-isa/internal/config"
	"cushion-isa/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func NewClient(cfg config.DatabaseConfig) (*Client, error) {
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (c *Client) Migrate() error {
	return c.db.AutoMigrate(&models.Investment{})
}

func (c *Client) CreateInvestment(inv *models.Investment) error {
	return c.db.Create(inv).Error
}

func (c *Client) GetInvestment(id uint) (*models.Investment, error) {
	var inv models.Investment
	if err := c.db.First(&inv, id).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (c *Client) ListInvestments() ([]models.Investment, error) {
	var investments []models.Investment
	if err := c.db.Find(&investments).Error; err != nil {
		return nil, err
	}
	return investments, nil
}

func (c *Client) UpdateInvestmentStatus(id uint, status string) error {
	return c.db.Model(&models.Investment{}).Where("id = ?", id).Update("status", status).Error
}
