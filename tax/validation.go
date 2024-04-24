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
		}
	}
	message.Message = "No valid allowance type found"
	message.Valitation = false

	return message
}
