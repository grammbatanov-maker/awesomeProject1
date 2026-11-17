package handlers

import (
	"awesomeProject1/internal/calculationService"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CalculationHandler struct {
	service calculationService.CalculationService
}

func NewCalculationHandler(s calculationService.CalculationService) *CalculationHandler {
	return &CalculationHandler{service: s}
}

func (h *CalculationHandler) GetCalculations(c echo.Context) error {
	calculations, err := h.service.GetAllCalculation()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get calculations"})
	}
	return c.JSON(http.StatusOK, calculations)
}

func (h *CalculationHandler) PostCalculation(c echo.Context) error {
	var req calculationService.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON request"})
	}

	calc, err := h.service.CreateCalculation(req.Expression)
	if err != nil {
		switch {
		case errors.Is(err, calculationService.ErrInvalidExpression):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			// Если ошибка от govaluate (неверное математическое выражение)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not evaluate expression: " + err.Error()})
		}
	}
	return c.JSON(http.StatusCreated, calc)
}

func (h *CalculationHandler) PatchCalculation(c echo.Context) error {
	id := c.Param("id")
	var req calculationService.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON request"})
	}

	updatedCalc, err := h.service.UpdateCalculation(id, req.Expression)
	if err != nil {
		switch {
		case errors.Is(err, calculationService.ErrInvalidID):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid UUID format"})
		case errors.Is(err, calculationService.ErrCalculationNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Calculation not found"})
		case errors.Is(err, calculationService.ErrInvalidExpression):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			// Ошибка парсинга/вычисления выражения
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not evaluate expression: " + err.Error()})
		}
	}
	return c.JSON(http.StatusOK, updatedCalc)
}

func (h *CalculationHandler) DeleteCalculation(c echo.Context) error {
	id := c.Param("id")

	if err := h.service.DeleteCalculation(id); err != nil {
		switch {
		case errors.Is(err, calculationService.ErrInvalidID):
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid UUID format"})
		case errors.Is(err, calculationService.ErrCalculationNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Calculation not found"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete calculation"})
		}
	}
	return c.NoContent(http.StatusNoContent)
}
