package tax

import (
	"strconv"
)

type InputErrorMeassager struct {
	Message    string
	Valitation bool
}

func valitationInpunt(input IncomeDetails, v TaxDiscountType) InputErrorMeassager {
	var message InputErrorMeassager
	if input.TotalIncome < 0 || input.WHT < 0 {
		message.Message = "TotalIncome and WHT  cannot be less than 0"
		message.Valitation = false
		return message
	}

	for _, i := range input.Allowances {
		if i.AllowanceType == v.Donation.Discount_Type {
			if i.Amount >= 0.0 {
				return InputErrorMeassager{Message: "", Valitation: true}
			} else {
				message.Message = "Amount is below minimum for Donation"
				message.Valitation = false
				return message
			}
		} else if i.AllowanceType == v.Kreceipt.Discount_Type {
			if i.Amount >= 0.0 {
				return InputErrorMeassager{Message: "", Valitation: true}
			} else {
				message.Message = "Amount is below minimum for K-receipt"
				message.Valitation = false
				return message
			}
		} else {
			message.Message = "No valid allowance type found"
			message.Valitation = false
		}
	}

	return message

}

func valitationSetingInpunt(input Amount, tdt string, v TaxDiscountType) InputErrorMeassager {

	var message InputErrorMeassager
	message.Valitation = true
	if tdt == "" {
		message.Message = "ใส่ประเภทของ ส่วนลดหยอนด้วยครับ!"
		message.Valitation = false
		return message
	}
	if tdt == v.Personal.Discount_Type {
		if input.Amount < v.Personal.Min_discount_value {
			message.Message = "input ที่ใส่มา ต้องมากกว่า 0 เเละ ต้อง มากกว่าขั้นต่ำที่กำหนดไว้"
			message.Valitation = false

		}
		if input.Amount > v.Personal.Max_discount_value {
			message.Message = "input ที่ใส่มา ต้อง น้อยกว่าค่าสูงสุดที่กำหนดไว้"
			message.Valitation = false

		}
		return message
	}

	return message
}
func valitationCsvFile(file [][]string) InputErrorMeassager {
	var message InputErrorMeassager
	message.Valitation = true

	for i, h := range file[:1] {
		if h[i] == "" {
			message.Message = "head ไม่เท่ากันค่าว่าง เเละเรียงตามนี้เสมอ  totalIncome,wht,donation "
			message.Valitation = false
			return message
		}
		if h[0] != "totalIncome" || h[1] != "wht" || h[2] != "donation" {
			message.Message = "ตำเเหน่งของ head ต้องเรียงตามนี้ totalIncome,wht,donation"
			message.Valitation = false
			return message
		}
	}
	for _, r := range file[1:] {
		t, _ := strconv.ParseFloat(r[0], 64)
		w, _ := strconv.ParseFloat(r[1], 64)
		d, _ := strconv.ParseFloat(r[2], 64)
		if t < 0 && w < 0 && d < 0 {
			message.Message = "ค่า Totalincome, wht, donation ต้องมากกว่า 0 เเละต้องไม่เป็นค่าว่าง หรือ ตัวหนังสือ"
			message.Valitation = false
			return message
		}

	}

	return message
}
