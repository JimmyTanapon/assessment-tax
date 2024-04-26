package tax

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
