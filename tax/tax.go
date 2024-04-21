package tax

import "log"

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

func (tb TaxBracket) CalculateTaxRate(income float64, prvlevel float64) float64 {
	log.Println(prvlevel, income)

	if income <= prvlevel || tb.MaxIncome == nil {
		return 0
	}

	if income > *tb.MaxIncome {
		log.Printf("Max incoome %f - MinIncome %f  * TaxRate %f ", *tb.MaxIncome, prvlevel, tb.TaxRate/100)
		return (*tb.MaxIncome - prvlevel) * (tb.TaxRate)
	}
	// log.Printf(" income %f - MinIncome %f  * TaxRate %f ", income, tb.MinIncome-1, tb.TaxRate)

	return (income - prvlevel) * (tb.TaxRate)
}

// func TexResponse(taxs []TaxBracket) {
// 	for _, tax := range taxs {
// 		fmt.Println(tax)
// 	}
// }
