package tax

import (
	"net/http"

	"github.com/JimmyTanapon/assessment-tax/helper"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	store Storer
}
type TaxDiscountType struct {
	Personal TaxDiscount
	Donation TaxDiscount
	Kreceipt TaxDiscount
}
type Err struct {
	Message string `json:"message"`
}
type TaxResponse struct {
	Tax   float64           `json:"tax"`
	Level []TaxLevelRespose `json:"taxLevel"`
}
type TaxLevelRespose struct {
	Level      string  `json:"level"`
	TaxinLevel float64 `json:"tax"`
}
type Storer interface {
	Discounts() (TaxDiscountType, error)
	SettingDeductionWithType(t string, amount float64) (string, float64, error)
}
type Amount struct {
	Amount float64
}
type TResponse struct {
	Type   string
	Amount float64
}

func New(db Storer) *Handler {
	return &Handler{store: db}
}

var TaxLevel = []TaxChart{
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
	discount := h.getTaxReduction()
	v := valitationInpunt(incomeDetails, discount)
	if !v.Valitation {
		return c.JSON(http.StatusBadRequest, v.Message)
	}

	taxAmount := incomeDetails.CalculateTax(discount)
	// return c.JSON(http.StatusOK, TaxResponse{Tax: taxAmount})
	return helper.SuccessHandler(c, taxAmount)

}

func (h *Handler) UpDeductionHandler(c echo.Context) error {

	ductionType := c.Param("type")

	var amount Amount
	var a float64
	var s string
	if err := c.Bind(&amount); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	discount := h.getTaxReduction()
	v := valitationSetingInpunt(amount, ductionType, discount)

	if v.Valitation {
		return c.JSON(http.StatusBadRequest, v.Message)
	}

	s, a, err := h.store.SettingDeductionWithType(ductionType, amount.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	var response = TResponse{
		Type:   s,
		Amount: a,
	}

	return helper.SuccessHandler(c, response)

}
