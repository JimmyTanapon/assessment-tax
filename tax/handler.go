package tax

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Err struct {
	Message string `json:"message"`
}
type TaxResponse struct {
	Tax float64 `json:"tax"`
}

func TaxHanler(c echo.Context) error {
	var incomeDetails IncomeDetails
	if err := c.Bind(&incomeDetails); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	taxAmount := calculateTax(incomeDetails.TotalIncome)
	return c.JSON(http.StatusOK, TaxResponse{Tax: taxAmount})

}

func calculateTaxRate(income float64) float64 {
	var tax float64
	if income <= 150000 {
		tax = 0
	} else if income <= 500000 {
		tax = (income - 150000) * 0.10
	} else if income <= 1000000 {
		tax = (500000-150000)*0.10 + (income-500000)*0.15
	} else if income <= 2000000 {
		tax = (500000-150000)*0.10 + (1000000-500000)*0.15 + (income-1000000)*0.20
	} else {
		tax = (500000-150000)*0.10 + (1000000-500000)*0.15 + (2000000-1000000)*0.20 + (income-2000000)*0.35
	}
	return tax
}
func calculateTax(income float64) float64 {
	taxReduction := 60000.0
	// Calculate taxable income after deduction
	taxableIncome := income - taxReduction
	tax := calculateTaxRate(taxableIncome)
	// Calculate tax amount using tax rate

	return tax
}
