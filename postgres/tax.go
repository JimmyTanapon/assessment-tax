package postgres

import (
	"fmt"
	"time"

	"github.com/JimmyTanapon/assessment-tax/tax"
)

type TaxDiscount struct {
	ID                 int       `postgres:"id"`
	Discount_Type      string    `postgres:"discount_type"`
	Discount_value     float64   `postgres:"discount_value"`
	Min_discount_value float64   `postgres:"min_discount_value"`
	Max_discount_value float64   `postgres:"max_discount_value"`
	CreatedAt          time.Time `postgres:"created_at"`
}
type TaxDiscountType struct {
	PersonalDeduction TaxDiscount
	Donation          TaxDiscount
	Kreceipt          TaxDiscount
}

func (p *Postgres) Discounts() (tax.TaxDiscountType, error) {
	rows, err := p.Db.Query("SELECT * FROM  tax_discount")
	if err != nil {
		return tax.TaxDiscountType{}, err
	}
	defer rows.Close()
	var taxDiscount []tax.TaxDiscount
	for rows.Next() {
		var td TaxDiscount
		err := rows.Scan(
			&td.ID, &td.Discount_Type,
			&td.Discount_value,
			&td.Min_discount_value,
			&td.Max_discount_value,
			&td.CreatedAt,
		)
		if err != nil {
			return tax.TaxDiscountType{}, err
		}
		taxDiscount = append(taxDiscount, tax.TaxDiscount{
			ID:                 td.ID,
			Discount_Type:      td.Discount_Type,
			Discount_value:     td.Discount_value,
			Min_discount_value: td.Min_discount_value,
			Max_discount_value: td.Max_discount_value,
			CreatedAt:          td.CreatedAt,
		})
	}

	result, crateerr := CreateTaxGroup(taxDiscount)
	if crateerr != nil {
		return tax.TaxDiscountType{}, err
	}
	return result, nil
}

func CreateTaxGroup(taxDiscount []tax.TaxDiscount) (tax.TaxDiscountType, error) {
	var personalDeduction, donation, kreceipt tax.TaxDiscount

	for _, discount := range taxDiscount {
		switch discount.Discount_Type {
		case "personalDeduction":
			personalDeduction = discount
		case "donation":
			donation = discount
		case "k-receipt":
			kreceipt = discount
		default:
			return tax.TaxDiscountType{}, fmt.Errorf("Message")
		}
	}
	var typeofTaxDicounst = tax.TaxDiscountType{
		PersonalDeduction: personalDeduction,
		Donation:          donation,
		Kreceipt:          kreceipt,
	}

	return typeofTaxDicounst, nil
}
