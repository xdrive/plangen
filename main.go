package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/xdrive/plangen/calc"
)

func main() {
	t, _ := time.Parse(time.RFC3339, "2018-01-01T00:00:01Z")
	plan := calc.CalcAnnuityPlan(500000, 24, 5.0, t)

	res, _ := json.MarshalIndent(plan, "", "  ")
	fmt.Println(string(res))
}
