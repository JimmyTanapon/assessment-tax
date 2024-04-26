package tax

import (
	"log"
	"math"
	"time"
)

type TaxDiscount struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Discount_Type      string    `json:"discount_type"`
	Discount_value     float64   `json:"discount_value"`
	Min_discount_value float64   `json:"min_discount_value"`
	Max_discount_value float64   `json:"max_discount_value"`
	CreatedAt          time.Time `json:"created_at"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type IncomeDetails struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

func (h *Handler) getTaxReduction() TaxDiscountType {
	discounts, err := h.store.Discounts()

	if err != nil {
		log.Fatal("err!", err)
	}

	return discounts
}
func (income IncomeDetails) CalculateTaxDiscount(dic TaxDiscountType) float64 {
	var totalDiscount float64
	for _, discount := range income.Allowances {
		if discount.AllowanceType == dic.Donation.Discount_Type {
			if discount.Amount > dic.Donation.Max_discount_value {
				totalDiscount += dic.Donation.Max_discount_value
			} else {
				totalDiscount += discount.Amount
			}
		}
	}
	return income.TotalIncome - dic.Personal.Discount_value - totalDiscount
}

func (income IncomeDetails) CalculateTax(dic TaxDiscountType) TaxResponse {
	taxableIncome := income.CalculateTaxDiscount(dic)
	var taxAmount float64
	maxIncomLevel := *TaxLevel[len(TaxLevel)-1].MaxIncome
	var taxLevelRespose = []TaxLevelRespose{}
	for index, tax := range TaxLevel {
		log.Println("here", index)
		if index > len(TaxLevel) {
			continue
		}
		adjustindex := index - 1
		if adjustindex < 0 {
			adjustindex = 0
		}
		// if taxableIncome < *TaxLevel[adjustindex].MaxIncome {
		// 	continue
		// }

		tlr := TaxLevelRespose{
			Level:      TaxLevel[index].Description,
			TaxinLevel: tax.CalculateTaxRate(taxableIncome, *TaxLevel[adjustindex].MaxIncome, maxIncomLevel),
		}
		taxAmount += tlr.TaxinLevel
		log.Println("Taxamount:", taxAmount)
		taxLevelRespose = append(taxLevelRespose, tlr)
	}

	var taxResponse = TaxResponse{
		Tax:   math.Round((taxAmount-income.WHT)*100) / 100,
		Level: taxLevelRespose,
	}

	return taxResponse
}

type TaxChart struct {
	MinIncome   float64  `json:"min_income"`
	MaxIncome   *float64 `json:"max_income,omitempty"`
	TaxRate     float64  `json:"tax_rate"`
	Description string   `json:"description"`
}

func (tb TaxChart) CalculateTaxRate(income float64, prvlevel float64, maxlevel float64) float64 {

	if income <= prvlevel {
		return 0
	}
	if income > *tb.MaxIncome && *tb.MaxIncome < maxlevel {
		return math.Round(((*tb.MaxIncome-prvlevel)*(tb.TaxRate))*100) / 100
	}

	return math.Round(((income-prvlevel)*(tb.TaxRate))*100) / 100

}
