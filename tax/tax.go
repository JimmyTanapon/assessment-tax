package tax

import (
	"log"
	"time"
)

type TaxDiscount struct {
	ID                 int       `json:"id"`
	Discount_Type      string    `json:"discount_type"`
	Discount_value     float64   `json:"discount_value"`
	Min_discount_value float64   `json:"min_discount_value"`
	Max_discount_value float64   `json:"max_discount_value"`
	CreatedAt          time.Time `json:"created_at"`
}
type TaxDiscountType struct {
	PersonalDeduction TaxDiscount
	Donation          TaxDiscount
	Kreceipt          TaxDiscount
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

func (h *Handler) taxReduction() TaxDiscountType {
	discounts, err := h.store.Discounts()
	var personalDeduction, donation, kreceipt TaxDiscount

	if err != nil {
		log.Fatal("err!", err)
	}
	for _, discount := range discounts {
		switch discount.Discount_Type {
		case "personalDeduction":
			personalDeduction = discount
		case "donation":
			donation = discount
		case "k-receipt":
			kreceipt = discount
		default:
			log.Fatal("err!", discount)
		}
	}
	var typeofTaxDicounst = TaxDiscountType{
		PersonalDeduction: personalDeduction,
		Donation:          donation,
		Kreceipt:          kreceipt,
	}

	return typeofTaxDicounst
}

func (h *Handler) CalculateTax(income IncomeDetails) float64 {
	taxType := h.taxReduction()
	taxableIncome := income.TotalIncome - taxType.PersonalDeduction.Discount_value
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

	return taxAmount - income.WHT
}

type TaxBracket struct {
	MinIncome   float64  `json:"min_income"`
	MaxIncome   *float64 `json:"max_income,omitempty"`
	TaxRate     float64  `json:"tax_rate"`
	Description string   `json:"description"`
}

func (tb TaxBracket) CalculateTaxRate(income float64, prvlevel float64, maxlevel float64) float64 {

	if income <= prvlevel || tb.MaxIncome == nil {
		return 0
	}
	if income > *tb.MaxIncome && income < maxlevel {
		return (*tb.MaxIncome - prvlevel) * (tb.TaxRate)
	}
	return (income - prvlevel) * (tb.TaxRate)
}
