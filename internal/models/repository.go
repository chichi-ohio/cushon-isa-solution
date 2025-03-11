package models

// Repository defines the interface for investment data operations
type Repository interface {
	Create(inv *Investment) error
	GetByID(id uint) (*Investment, error)
	UpdateStatus(id uint, status string) error
	List(limit, offset int) ([]Investment, error)
}
