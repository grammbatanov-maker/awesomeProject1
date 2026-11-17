package calculationService

import (
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
)

type CalculationService interface {
	CreateCalculation(expression string) (Calculation, error)
	GetAllCalculation() ([]Calculation, error)
	GetCalculationById(id string) (Calculation, error)
	UpdateCalculation(id, expression string) (Calculation, error)
	DeleteCalculation(id string) error
}

type calcService struct {
	repo CalculationRepository
}

func NewCalculationService(r CalculationRepository) CalculationService {
	return &calcService{repo: r}
}

func (s *calcService) calculateExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return "", err
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), nil
}

// Вспомогательная функция для проверки UUID
func validateID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}
	return nil
}

func (s *calcService) CreateCalculation(expression string) (Calculation, error) {
	if strings.TrimSpace(expression) == "" {
		return Calculation{}, ErrInvalidExpression
	}

	result, err := s.calculateExpression(expression)
	if err != nil {
		return Calculation{}, err
	}

	calc := Calculation{
		ID:         uuid.NewString(),
		Expression: expression,
		Result:     result,
	}
	if err := s.repo.CreateCalculation(calc); err != nil {
		return Calculation{}, err
	}

	return calc, nil
}

func (s *calcService) GetAllCalculation() ([]Calculation, error) {
	return s.repo.GetAllCalculation()
}

func (s *calcService) GetCalculationById(id string) (Calculation, error) {
	if err := validateID(id); err != nil {
		return Calculation{}, err
	}
	return s.repo.GetCalculationByID(id)
}

func (s *calcService) UpdateCalculation(id, expression string) (Calculation, error) {
	if err := validateID(id); err != nil {
		return Calculation{}, err
	}
	if strings.TrimSpace(expression) == "" {
		return Calculation{}, ErrInvalidExpression
	}

	calc, err := s.repo.GetCalculationByID(id)
	if err != nil {
		return Calculation{}, err
	}
	result, err := s.calculateExpression(expression)
	if err != nil {
		return Calculation{}, err
	}

	calc.Expression = expression
	calc.Result = result

	if err := s.repo.UpdateCalculation(calc); err != nil {
		return Calculation{}, err
	}

	return calc, nil
}

func (s *calcService) DeleteCalculation(id string) error {
	if err := validateID(id); err != nil {
		return err
	}
	return s.repo.DeleteCalculation(id)
}
