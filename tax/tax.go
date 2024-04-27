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
			if discount.Amount > dic.Donation.Discount_value {
				totalDiscount += dic.Donation.Discount_value
			} else {
				totalDiscount += discount.Amount
			}

		}
		if discount.AllowanceType == dic.Kreceipt.Discount_Type {
			if discount.Amount > dic.Kreceipt.Discount_value {
				totalDiscount += dic.Kreceipt.Discount_value
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
	taxRefund := 0.00
	maxIncomLevel := *TaxLevel[len(TaxLevel)-1].MaxIncome
	var taxLevelRespose = []TaxLevelRespose{}
	for index, tax := range TaxLevel {
		if index > len(TaxLevel) {
			continue
		}
		adjustindex := index - 1
		if adjustindex < 0 {
			adjustindex = 0
		}

		tlr := TaxLevelRespose{
			Level:      TaxLevel[index].Description,
			TaxinLevel: tax.CalculateTaxRate(taxableIncome, *TaxLevel[adjustindex].MaxIncome, maxIncomLevel),
		}
		taxAmount += tlr.TaxinLevel
		taxLevelRespose = append(taxLevelRespose, tlr)
	}
	result := math.Round((taxAmount-income.WHT)*100) / 100
	if result < 0 {
		taxRefund = -1 * result
	}

	var taxResponse = TaxResponse{
		Tax:       result,
		Level:     taxLevelRespose,
		TaxRefund: taxRefund,
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

type Taxes struct {
	TotalIncome float64
	Tax         float64
	TaxRefund   float64
}

func CalculateTaxCsv(taxesCsv []IncomeDetails, dic TaxDiscountType) []Taxes {
	var taxes []Taxes

	for _, tax := range taxesCsv {
		amount := tax.CalculateTax(dic)

		tax := Taxes{
			TotalIncome: tax.TotalIncome,
			Tax:         amount.Tax,
			TaxRefund:   amount.TaxRefund,
		}

		taxes = append(taxes, tax)

	}

	return taxes

}
