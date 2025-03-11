package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	Email        string       `json:"email" gorm:"uniqueIndex" binding:"required,email"`
	Name         string       `json:"name" binding:"required"`
	DateOfBirth  time.Time    `json:"date_of_birth" binding:"required"`
	Address      string       `json:"address" binding:"required"`
	Investments  []Investment `json:"investments" gorm:"foreignKey:CustomerID"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	PasswordHash string       `json:"-"` // Password hash not exposed in JSON
}

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) Create(customer *Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) GetByID(id uint) (*Customer, error) {
	var customer Customer
	err := r.db.Preload("Investments").First(&customer, id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) GetByEmail(email string) (*Customer, error) {
	var customer Customer
	err := r.db.Where("email = ?", email).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) List() ([]Customer, error) {
	var customers []Customer
	err := r.db.Find(&customers).Error
	return customers, err
}

func (r *CustomerRepository) Update(customer *Customer) error {
	return r.db.Save(customer).Error
}
