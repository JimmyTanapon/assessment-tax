package tax

import (
	"reflect"
	"testing"
	"time"
)

var donation = TaxDiscount{
	ID:                 1,
	Name:               "donation",
	Discount_Type:      "donation",
	Discount_value:     100000.00,
	Min_discount_value: 0.0,
	Max_discount_value: 100000,
	CreatedAt:          time.Now(),
}
var personal = TaxDiscount{
	ID:                 1,
	Name:               "personalDeduction",
	Discount_Type:      "personal",
	Discount_value:     60000.00,
	Min_discount_value: 10000.00,
	Max_discount_value: 100000.00,
	CreatedAt:          time.Now(),
}
var kreceipt = TaxDiscount{
	ID:                 1,
	Name:               "kreceipt",
	Discount_Type:      "k-receipt",
	Discount_value:     50000.00,
	Min_discount_value: 0.0,
	Max_discount_value: 100000.00,
	CreatedAt:          time.Now(),
}
var TaxDiscountStup = TaxDiscountType{
	Kreceipt: kreceipt,
	Personal: personal,
	Donation: donation,
}

func TestValitationInpunt(t *testing.T) {
	t.Run("totalincome less than 0", func(t *testing.T) {
		income := IncomeDetails{
			TotalIncome: -100000,
			WHT:         0,
			Allowances: []Allowance{
				{AllowanceType: "Housing", Amount: 2000},
				{AllowanceType: "Transport", Amount: 1500},
			},
		}

		var want = InputErrorMeassager{
			Message:    "TotalIncome and WHT  cannot be less than 0",
			Valitation: false,
		}
		got := valitationInpunt(income, TaxDiscountStup)

		if want != got {
			t.Errorf("Want %v, got %v", want, got)
		}

	})
	t.Run("Wht less than 0", func(t *testing.T) {
		income := IncomeDetails{
			TotalIncome: 100000,
			WHT:         -100000,
			Allowances: []Allowance{
				{AllowanceType: "Housing", Amount: 2000},
				{AllowanceType: "Transport", Amount: 1500},
			},
		}

		var want = InputErrorMeassager{
			Message:    "TotalIncome and WHT  cannot be less than 0",
			Valitation: false,
		}
		got := valitationInpunt(income, TaxDiscountStup)

		if want != got {
			t.Errorf("Want %v, got %v", want, got)
		}

	})
	t.Run("Wht more than totalincom", func(t *testing.T) {

		income := IncomeDetails{
			TotalIncome: 500000,
			WHT:         600000,
			Allowances: []Allowance{
				{AllowanceType: "Housing", Amount: 2000},
				{AllowanceType: "Transport", Amount: 1500},
			},
		}

		var want = InputErrorMeassager{
			Message:    "ข้อมูล wht ที่จะถูกส่งเข้ามาคำนวน ไม่สามารถมีค่าน้อยกว่า 0 หรือมากกว่ารายรับได้",
			Valitation: false,
		}
		got := valitationInpunt(income, TaxDiscountStup)

		if want != got {
			t.Errorf("Want %v, got %v", want, got)
		}

	})
	t.Run("invalid AllowanceType", func(t *testing.T) {

		income := IncomeDetails{
			TotalIncome: 500000,
			WHT:         0.0,
			Allowances: []Allowance{
				{AllowanceType: "Housing", Amount: 2000},
				{AllowanceType: "Transport", Amount: 1500},
			},
		}

		var want = InputErrorMeassager{
			Message:    "No valid allowance type found",
			Valitation: false,
		}
		got := valitationInpunt(income, TaxDiscountStup)

		if want != got {
			t.Errorf("Want %v, got %v", want, got)
		}

	})
	t.Run("valid AllowanceType", func(t *testing.T) {

		income := IncomeDetails{
			TotalIncome: 500000,
			WHT:         0.0,
			Allowances: []Allowance{
				{AllowanceType: "donation", Amount: 2000},
				{AllowanceType: "k-receipt", Amount: 1500},
			},
		}

		var want = InputErrorMeassager{
			Message:    "",
			Valitation: true,
		}
		got := valitationInpunt(income, TaxDiscountStup)

		if want != got {
			t.Errorf("Want %v, got %v", want, got)
		}

	})
	t.Run("donation input less than 0", func(t *testing.T) {

		income := IncomeDetails{
			TotalIncome: 500000,
			WHT:         0.0,
			Allowances: []Allowance{
				{AllowanceType: "donation", Amount: -10000.00},
				{AllowanceType: "k-receipt", Amount: 10000.00},
			},
		}

		var want = InputErrorMeassager{
			Message:    "Amount is below minimum for Donation",
			Valitation: false,
		}
		got := valitationInpunt(income, TaxDiscountStup)
		if want != got {
			t.Errorf("Want %v, got %v", want, got)
		}

	})
	t.Run("k-receipt input less than 0", func(t *testing.T) {

		income := IncomeDetails{
			TotalIncome: 500000,
			WHT:         0.0,
			Allowances: []Allowance{
				{AllowanceType: "donation", Amount: 10000.00},
				{AllowanceType: "k-receipt", Amount: -10000.00},
			},
		}

		var want = InputErrorMeassager{
			Message:    "Amount is below minimum for K-receipt",
			Valitation: false,
		}
		got := valitationInpunt(income, TaxDiscountStup)
		if want != got {
			t.Errorf("Want %v, got %v", want, got)
		}

	})

}
func TestValitationSetingInpunt(t *testing.T) {
	var cases = []struct {
		name             string
		input            Amount
		texdeductiontype string

		want InputErrorMeassager
	}{
		{name: "texdeduction input  is null  should be massage and false", input: Amount{Amount: 70000.00}, texdeductiontype: "", want: InputErrorMeassager{Message: "ระบุประเภทของส่วนลดหย่อนไม่ถูกต้อง!", Valitation: false}},
		{name: "texdeduction input is invalid should be massage and false", input: Amount{Amount: 70000.00}, texdeductiontype: "supershy", want: InputErrorMeassager{Message: "ระบุประเภทของส่วนลดหย่อนไม่ถูกต้อง!", Valitation: false}},
		{name: "personal input valid should be  true", input: Amount{Amount: 70000.00}, texdeductiontype: "personal", want: InputErrorMeassager{Message: "", Valitation: true}},
		{name: "personal input  valid but  Amount input is invalid should be massage and false   ", input: Amount{Amount: -70000.00}, texdeductiontype: "personal", want: InputErrorMeassager{Message: "input Personal ที่ใส่มา ต้องมากกว่า10000", Valitation: false}},
		{name: "personal input valid but amount more than specified  should be massage and false   ", input: Amount{Amount: 200000.00}, texdeductiontype: "personal", want: InputErrorMeassager{Message: "input  Personal ที่ใส่มาต้องน้อยกว่า100000", Valitation: false}},
		{name: "personal input  valid but  Amount input is 0 should be massage and false   ", input: Amount{Amount: 0.00}, texdeductiontype: "personal", want: InputErrorMeassager{Message: "input Personal ที่ใส่มา ต้องมากกว่า10000", Valitation: false}},
		{name: "k-receipt input valid should be  true", input: Amount{Amount: 70000.00}, texdeductiontype: "k-receipt", want: InputErrorMeassager{Message: "", Valitation: true}},
		{name: "k-receipt input  valid but  Amount input is invalid should be massage and false   ", input: Amount{Amount: -70000.00}, texdeductiontype: "k-receipt", want: InputErrorMeassager{Message: "input Kreceipt  ที่ใส่มา ต้องมากกว่า0", Valitation: false}},
		{name: "k-receipt input valid but amount more than specified  should be massage and false   ", input: Amount{Amount: 200000.00}, texdeductiontype: "k-receipt", want: InputErrorMeassager{Message: "input Kreceipt ที่ใส่มาต้องน้อยกว่า100000", Valitation: false}},
		{name: "k-receipt input  valid but  Amount input is 0 should be massage and false   ", input: Amount{Amount: 0.00}, texdeductiontype: "k-receipt", want: InputErrorMeassager{Message: "input Kreceipt  ที่ใส่มา ต้องมากกว่า0", Valitation: false}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			got := valitationSetingInpunt(test.input, test.texdeductiontype, TaxDiscountStup)

			if got != test.want {
				t.Errorf("got = %v, want %v", got, test.want)
			}
		})

	}
}
func TestDiscounts(t *testing.T) {
	t.Run(" TotalIncome: 500000 and donation 200000 should be 340,000", func(t *testing.T) {
		income := IncomeDetails{
			TotalIncome: 500000,
			WHT:         0.0,
			Allowances: []Allowance{
				{AllowanceType: "donation", Amount: 200000.00},
			},
		}
		want := 340000.00
		got := income.CalculateTaxDiscount(TaxDiscountStup)
		if want != got {
			t.Errorf("Want %v, got %v", want, got)
		}
	})
	t.Run(" TotalIncome: 500,000 and donation 100,000 and k-receipt  200000.0 should be  290,000", func(t *testing.T) {
		income := IncomeDetails{
			TotalIncome: 500000,
			WHT:         0.0,
			Allowances: []Allowance{
				{AllowanceType: "donation", Amount: 100000.00},
				{AllowanceType: "k-receipt", Amount: 200000.00},
			},
		}
		want := 290000.00
		got := income.CalculateTaxDiscount(TaxDiscountStup)
		if want != got {
			t.Errorf("Want %v, got %v", want, got)
		}
	})

}
func TestCalculateTax(t *testing.T) {
	t.Run("TotalIncome: 500000,donation:100000.00,k-receipt:200000 should return  14000 ", func(t *testing.T) {
		income := IncomeDetails{
			TotalIncome: 500000,
			WHT:         0.0,
			Allowances: []Allowance{
				{AllowanceType: "donation", Amount: 100000.00},
				{AllowanceType: "k-receipt", Amount: 200000.00},
			},
		}
		want := TaxResponse{
			Tax:       14000,
			TaxRefund: 0,
			Level: []TaxLevelRespose{
				{Level: "0-150,000", TaxinLevel: 0},
				{Level: "150,001-500,000", TaxinLevel: 14000},
				{Level: "500,001-1,000,000", TaxinLevel: 0},
				{Level: "1,000,001-2,000,000", TaxinLevel: 0},
				{Level: "2,000,001 ขึ้นไป", TaxinLevel: 0},
			},
		}
		got := income.CalculateTax(TaxDiscountStup) // Assuming TaxDiscountSetup is correctly defined
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Want %v, got %v", want, got)
		}
	})
	t.Run("TotalIncome: 500000,donation:100000.00 should return  19000 ", func(t *testing.T) {
		income := IncomeDetails{
			TotalIncome: 500000,
			WHT:         0.0,
			Allowances: []Allowance{
				{AllowanceType: "donation", Amount: 200000.00},
			},
		}
		want := TaxResponse{
			Tax:       19000,
			TaxRefund: 0,
			Level: []TaxLevelRespose{
				{Level: "0-150,000", TaxinLevel: 0},
				{Level: "150,001-500,000", TaxinLevel: 19000},
				{Level: "500,001-1,000,000", TaxinLevel: 0},
				{Level: "1,000,001-2,000,000", TaxinLevel: 0},
				{Level: "2,000,001 ขึ้นไป", TaxinLevel: 0},
			},
		}
		got := income.CalculateTax(TaxDiscountStup) // Assuming TaxDiscountSetup is correctly defined
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Want %v, got %v", want, got)
		}
	})
	t.Run("TotalIncome: 500000,WHT:25000.00 should return  4000 ", func(t *testing.T) {
		income := IncomeDetails{
			TotalIncome: 500000,
			WHT:         25000.00,
			Allowances: []Allowance{
				{AllowanceType: "donation", Amount: 0.00},
			},
		}
		want := TaxResponse{
			Tax:       4000,
			TaxRefund: 0,
			Level: []TaxLevelRespose{
				{Level: "0-150,000", TaxinLevel: 0},
				{Level: "150,001-500,000", TaxinLevel: 29000},
				{Level: "500,001-1,000,000", TaxinLevel: 0},
				{Level: "1,000,001-2,000,000", TaxinLevel: 0},
				{Level: "2,000,001 ขึ้นไป", TaxinLevel: 0},
			},
		}
		got := income.CalculateTax(TaxDiscountStup) // Assuming TaxDiscountSetup is correctly defined
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Want %v, got %v", want, got)
		}
	})

}
