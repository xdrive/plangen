package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/xdrive/plangen/calc"
)

type generateRequest struct {
	LoanAmount  string
	Duration    int
	NominalRate string
	StartDate   string
}

func (gr *generateRequest) validate() error {
	if gr.Duration <= 0 {
		return errors.New("duration request param is missing/invalid")
	}

	loanAmount, err := strconv.ParseInt(gr.LoanAmount, 10, 64)
	if err != nil {
		return fmt.Errorf("loanAmount request param is invalid: %v", err)
	}
	if loanAmount <= 0 {
		return fmt.Errorf("incorrect loan amount: %v", loanAmount)
	}

	nominalRate, err := strconv.ParseFloat(gr.NominalRate, 64)
	if err != nil {
		return fmt.Errorf("nominalRate request param is invalid: %v", err)
	}
	if nominalRate <= 0 {
		return fmt.Errorf("incorrect nominalRate: %v", nominalRate)
	}

	_, err = time.Parse(time.RFC3339, gr.StartDate)
	if err != nil {
		return fmt.Errorf("startDate request param is invalid: %v", err)
	}

	return nil
}

func generatePlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var gr generateRequest
	err := decoder.Decode(&gr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error parsing json request body:", err)
		return
	}
	if err = gr.validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error validating request params:", err)
		return
	}
	log.Println("request params:", gr)

	// parse errors already checked in request validate method
	loanAmount, _ := strconv.ParseInt(gr.LoanAmount, 10, 64)
	nominalRate, _ := strconv.ParseFloat(gr.NominalRate, 64)
	startDate, _ := time.Parse(time.RFC3339, gr.StartDate)

	plan := calc.CalcAnnuityPlan(loanAmount*100, gr.Duration, nominalRate, startDate)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(plan)
}
