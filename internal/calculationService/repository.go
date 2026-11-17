package calculationService

import (
	"gorm.io/gorm"
)

type CalculationRepository interface {
	CreateCalculation(calc Calculation) error
	GetAllCalculation() ([]Calculation, error)
	GetCalculationByID(id string) (Calculation, error)
	UpdateCalculation(calc Calculation) error
	DeleteCalculation(id string) error
}

type calcRepository struct {
	db *gorm.DB
}

func NewCalculationRepository(db *gorm.DB) CalculationRepository {
	return &calcRepository{db: db}
}

func (r *calcRepository) CreateCalculation(calc Calculation) error {
	return r.db.Create(&calc).Error
}

func (r *calcRepository) GetAllCalculation() ([]Calculation, error) {
	var calculations []Calculation
	err := r.db.Find(&calculations).Error
	return calculations, err
}

func (r *calcRepository) GetCalculationByID(id string) (Calculation, error) {
	var calc Calculation
	err := r.db.First(&calc, "id = ?", id).Error
	return calc, err
}

func (r *calcRepository) UpdateCalculation(calc Calculation) error {
	// Используем Updates вместо Save, чтобы проверить RowsAffected
	result := r.db.Model(&Calculation{}).Where("id = ?", calc.ID).Updates(&calc)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCalculationNotFound
	}
	return nil
}

func (r *calcRepository) DeleteCalculation(id string) error {
	result := r.db.Delete(&Calculation{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCalculationNotFound
	}
	return nil
}
