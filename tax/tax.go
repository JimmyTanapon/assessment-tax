package tax

import "time"

type TaxDiscount struct {
	ID                 int       `json:"id"`
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
