package calc

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

// AnnuityPlanPayment represents a single interval in annuity loan payback plan
type AnnuityPlanPayment struct {
	Date                          time.Time `json:"date"`
	PaymentAmount                 money     `json:"borrowerPaymentAmount"`
	PrincipalAmount               money     `json:"principal"`
	InterestAmount                money     `json:"interest"`
	InitialOutstandingPrincipal   money     `json:"initialOutstandingPrincipal"`
	RemainingOutstandingPrincipal money     `json:"remainingOutstandingPrincipal"`
}

// MarshalJSON customizes the json serialization with date formating
func (app *AnnuityPlanPayment) MarshalJSON() ([]byte, error) {
	type parent AnnuityPlanPayment
	return json.Marshal(&struct {
		Date string `json:"date"`
		*parent
	}{
		Date:   app.Date.Format(time.RFC3339),
		parent: (*parent)(app),
	})
}

// money type represents money value in cents. So 1.23 will be represented as 123
type money int64

func (m money) MarshalJSON() ([]byte, error) {
	val := fmt.Sprintf("\"%.2f\"", float64(m)/100.0)
	return []byte(val), nil
}

// AnnuityPlan calculates payback plan for annuity loan. The loan amount is provided in
// cents. Interest rate is in percents per year
func AnnuityPlan(amountCents int64, durationMonths int, annualInterestRatePercent float64, startDate time.Time) []AnnuityPlanPayment {
	plan := []AnnuityPlanPayment{}

	outstandingPrincipal := money(amountCents)
	paymentAmount := calcAnnuityPayment(outstandingPrincipal, durationMonths, annualInterestRatePercent/12)
	for i := 0; i < durationMonths; i++ {
		date := startDate.AddDate(0, 1*i, 0)
		interest := calcPaymentInterestPart(outstandingPrincipal, annualInterestRatePercent)
		if i == durationMonths-1 {
			paymentAmount = outstandingPrincipal + interest
		}
		remainingPrincipal := outstandingPrincipal - paymentAmount + interest
		planPayment := AnnuityPlanPayment{
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

func calcAnnuityPayment(amountCents money, durationMonths int, monthlyInterestRatePercent float64) money {
	rate := monthlyInterestRatePercent / 100
	monthlyAmount := float64(amountCents) * rate / (1 - math.Pow((1+rate), -float64(durationMonths)))

	return money(math.Round(monthlyAmount))
}

func calcPaymentInterestPart(outstandingPrincipal money, annualInterestRatePercent float64) money {
	return money(math.Round(annualInterestRatePercent * 30 * float64(outstandingPrincipal) / (360 * 100)))
}
