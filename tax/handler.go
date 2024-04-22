package tax

import (
	"log"
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
	h.taxReduction()
	taxAmount := calculateTax(incomeDetails.TotalIncome)
	return c.JSON(http.StatusOK, TaxResponse{Tax: taxAmount})

}
func calculateTax(income float64) float64 {

	taxReduction := 60000.0
	taxableIncome := income - taxReduction
	var taxAmount float64
	taxLevelindex := 0
	maxIncomLevel := *TaxLevel[len(TaxLevel)-1].MaxIncome

	for index, tax := range TaxLevel {

		if index > len(TaxLevel) {
			continue
		}
		adjustindex := index - 1
		if adjustindex < 0 {
			adjustindex = 0
		}
		if taxableIncome < *TaxLevel[adjustindex].MaxIncome {
			continue
		}
		taxLevelindex += 1
		taxAmount += tax.CalculateTaxRate(taxableIncome, *TaxLevel[adjustindex].MaxIncome, maxIncomLevel)
	}
	log.Println("taxlevle", taxLevelindex)

	return taxAmount
}
func (h *Handler) taxReduction() {
	discount, err := h.store.Discounts()
	if err != nil {
		log.Fatal("err!", err)
	}
	log.Println(discount)

}
