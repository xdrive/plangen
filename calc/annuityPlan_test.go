package calc

import "testing"

func TestCalcAnnuityPaymentRequirementExample(t *testing.T) {
	payment := calcAnnuityPayment(500000, 24, 5.0/12)
	if payment != 21936 {
		t.Error("Payment amount should be 21936 cents, got", payment)
	}
}
