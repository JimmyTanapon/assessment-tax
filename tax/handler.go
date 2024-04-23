package tax

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	store Storer
}

type Err struct {
	Message string `json:"message"`
}
type TaxResponse struct {
	Tax float64 `json:"tax"`
}
type Storer interface {
	Discounts() ([]TaxDiscount, error)
}

func New(db Storer) *Handler {
	return &Handler{store: db}
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

func (h *Handler) TaxHandler(c echo.Context) error {
	var incomeDetails IncomeDetails
	if err := c.Bind(&incomeDetails); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	taxAmount := h.CalculateTax(incomeDetails)
	return c.JSON(http.StatusOK, TaxResponse{Tax: taxAmount})

}
