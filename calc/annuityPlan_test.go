package calc

import "testing"

func TestCalcAnnuityPaymentRequirementExample(t *testing.T) {
	payment := calcAnnuityPayment(money(500000), 24, 5.0/12)
	if payment != 21936 {
		t.Error("Payment amount should be 21936 cents, got", payment)
	}
}

func TestCalcPaymentInterestPartRequirementExample(t *testing.T) {
	interest := calcPaymentInterestPart(money(500000), 5.0)
	if interest != 2083 {
		t.Error("Interest amount should be 2083 cents, got:", interest)
	}
}
