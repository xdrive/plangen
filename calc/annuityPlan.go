package calc

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type AnnuityPlan []AnnuityPlanPayment

type AnnuityPlanPayment struct {
	Date                          time.Time
	PaymentAmount                 int64
	PrincipalAmount               int64
	InterestAmount                int64
	InitialOutstandingPrincipal   int64
	RemainingOutstandingPrincipal int64
}

func (app AnnuityPlanPayment) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Date                          string `json:"date"`
		PaymentAmount                 money  `json:"borrowerPaymentAmount"`
		PrincipalAmount               money  `json:"principal"`
		InterestAmount                money  `json:"interest"`
		InitialOutstandingPrincipal   money  `json:"initialOutstandingPrincipal"`
		RemainingOutstandingPrincipal money  `json:"remainingOutstandingPrincipal"`
	}{
		Date:                          app.Date.Format(time.RFC3339),
		PaymentAmount:                 money(app.PaymentAmount),
		PrincipalAmount:               money(app.PrincipalAmount),
		InterestAmount:                money(app.InterestAmount),
		InitialOutstandingPrincipal:   money(app.InitialOutstandingPrincipal),
		RemainingOutstandingPrincipal: money(app.RemainingOutstandingPrincipal),
	})
}

// money type represents money value in cents. So 1.23 will be represented as 123
type money int64

func (m money) MarshalJSON() ([]byte, error) {
	val := fmt.Sprintf("%.2f", float64(m)/100.0)
	return []byte(val), nil
}

func CalcAnnuityPlan(amountCents int64, durationMonths int, annualInterestRatePercent float64, startDate time.Time) AnnuityPlan {
	plan := AnnuityPlan{}
	paymentAmount := calcAnnuityPayment(amountCents, durationMonths, annualInterestRatePercent/12)
	outstandingPrincipal := amountCents
	var planPayment AnnuityPlanPayment
	for i := 0; i < durationMonths; i++ {
		date := startDate.AddDate(0, 1*i, 0)
		interest := calcPaymentInterestPart(outstandingPrincipal, annualInterestRatePercent)
		if i == durationMonths-1 {
			paymentAmount = outstandingPrincipal + interest
		}
		remainingPrincipal := outstandingPrincipal - paymentAmount + interest
		planPayment = AnnuityPlanPayment{
			Date:                          date,
			PaymentAmount:                 paymentAmount,
			InterestAmount:                interest,
			PrincipalAmount:               paymentAmount - interest,
			InitialOutstandingPrincipal:   outstandingPrincipal,
			RemainingOutstandingPrincipal: remainingPrincipal,
		}
		plan = append(plan, planPayment)
		outstandingPrincipal = remainingPrincipal
	}

	return plan
}

func calcAnnuityPayment(amountCents int64, durationMonths int, monthlyInterestRatePercent float64) int64 {
	rate := monthlyInterestRatePercent / 100
	monthlyAmount := float64(amountCents) * rate / (1 - math.Pow((1+rate), -float64(durationMonths)))

	return int64(math.Round(monthlyAmount))
}

func calcPaymentInterestPart(outstandingPrincipal int64, annualInterestRatePercent float64) int64 {
	return int64(math.Round(annualInterestRatePercent * 30 * float64(outstandingPrincipal) / (360 * 100)))
}
