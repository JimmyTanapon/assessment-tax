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

var TaxLevel = []TaxBracket{
	{MaxIncome: floatPtr(150000), TaxRate: 0, Description: "0-150,000"},
	{MaxIncome: floatPtr(500000), TaxRate: 0.1, Description: "150,001-500,000"},
	{MaxIncome: floatPtr(1000000), TaxRate: 0.15, Description: "500,001-1,000,000"},
	{MaxIncome: floatPtr(2000000), TaxRate: 0.2, Description: "1,000,001-2,000,000"},
	{MaxIncome: floatPtr(2000001), TaxRate: 0.35, Description: "2,000,001 ขึ้นไป"},
}

func floatPtr(f float64) *float64 {
	return &f
}

func TaxHandler(c echo.Context) error {
	var incomeDetails IncomeDetails
	if err := c.Bind(&incomeDetails); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	taxAmount := calculateTax(incomeDetails.TotalIncome)
	return c.JSON(http.StatusOK, TaxResponse{Tax: taxAmount})

}

func calculateTax(income float64) float64 {

	taxReduction := 60000.0
	var taxAmount float64

	// Calculate taxable income after deduction
	taxableIncome := income - taxReduction
	for index, tax := range TaxLevel {
		if index > len(TaxLevel) {
			continue
		}
		adjustindex := index - 1
		if adjustindex < 0 {
			adjustindex = 0
		}
		taxAmount += tax.CalculateTaxRate(taxableIncome, *TaxLevel[adjustindex].MaxIncome)
	}

	return taxAmount
}
