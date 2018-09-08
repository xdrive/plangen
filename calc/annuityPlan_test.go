package calc

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

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

func TestCalcAnnuityPlanAllButLastIntervalsWithSameAmount(t *testing.T) {
	duration := 24
	paymentAmount := money(21936)
	plan := AnnuityPlan(500000, duration, 5.0, time.Now())
	for i, pp := range plan {
		if i < duration-1 && pp.PaymentAmount != paymentAmount {
			t.Errorf("Every but last payment interval should have payment amount of %d, got %d for interval %d",
				paymentAmount,
				pp.PaymentAmount,
				i)
		}
	}
}

func TestCalcAnnuityPlanLastInterval(t *testing.T) {
	duration := 5
	plan := AnnuityPlan(10000, duration, 10.0, time.Now())
	//beforeLastInterval := plan[duration-2]
	lastInterval := plan[duration-1]

	if lastInterval.RemainingOutstandingPrincipal != money(0) {
		t.Errorf("RemainingOutstandingPrincipal should be 0 in last interval, got: %d", lastInterval.RemainingOutstandingPrincipal)
	}

	if lastInterval.InitialOutstandingPrincipal != lastInterval.PrincipalAmount {
		t.Errorf("PrincipalAmount and InitialOutstandingPrincipal of last interval should match, got: %d and %d",
			lastInterval.PrincipalAmount,
			lastInterval.InitialOutstandingPrincipal)
	}

	if lastInterval.PaymentAmount != lastInterval.InterestAmount+lastInterval.PrincipalAmount {
		t.Errorf("InterestAmount + PrincipalAmount should equal to PaymentAmount, got: %d + %d and %d",
			lastInterval.InterestAmount,
			lastInterval.PrincipalAmount,
			lastInterval.PaymentAmount)
	}

	if lastInterval.PaymentAmount <= 0 {
		t.Errorf("PaymentAmount should be > 0, got %d", lastInterval.PaymentAmount)
	}
}

func TestMoneyToJson(t *testing.T) {
	m := money(1234)
	res, _ := json.Marshal(m)
	if string(res) != "\"12.34\"" {
		t.Errorf("should be 12.34, got %s", res)
	}
}

func TestAnnuityPlanPaymentToJson(t *testing.T) {
	date := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	paymentAmount := 100.00
	interestAmount := 10.00
	principalAmount := 90.00
	outstandingPrincipal := 900.00
	remainingPrincipal := 810.00

	p := &AnnuityPlanPayment{
		Date:                          date,
		PaymentAmount:                 money(paymentAmount * 100),
		InterestAmount:                money(interestAmount * 100),
		PrincipalAmount:               money(principalAmount * 100),
		InitialOutstandingPrincipal:   money(outstandingPrincipal * 100),
		RemainingOutstandingPrincipal: money(remainingPrincipal * 100),
	}

	expected := fmt.Sprintf(
		"{\"date\":\"%s\",\"borrowerPaymentAmount\":\"%.2f\",\"principal\":\"%.2f\",\"interest\":\"%.2f\",\"initialOutstandingPrincipal\":\"%.2f\",\"remainingOutstandingPrincipal\":\"%.2f\"}",
		date.Format(time.RFC3339),
		paymentAmount,
		principalAmount,
		interestAmount,
		outstandingPrincipal,
		remainingPrincipal)

	res, _ := json.Marshal(p)
	if string(res) != expected {
		t.Error(string(res))
	}
}
