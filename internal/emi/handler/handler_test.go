package handler

import (
	"encoding/json"
	"testing"
)

func TestAndroidEmiPayload(t *testing.T) {
	raw := `[
	  {
	    "amount": "9000",
	    "monthsPaid": "3",
	    "name": "ICICI EMI 1",
	    "startDate": "07/10/2026",
	    "tenure": "9",
	    "type": "expense"
	  }
	]`

	var reqs []emiReq
	if err := json.Unmarshal([]byte(raw), &reqs); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(reqs) != 1 {
		t.Fatalf("got %d reqs", len(reqs))
	}

	e, err := reqs[0].toModel("user-1")
	if err != nil {
		t.Fatal(err)
	}
	if e.Amount != 9000 {
		t.Fatalf("amount %v", e.Amount)
	}
	if e.Type != "expense" {
		t.Fatalf("type %s", e.Type)
	}
	if e.Name == nil || *e.Name != "ICICI EMI 1" {
		t.Fatalf("name %+v", e.Name)
	}
	if e.Description == nil || *e.Description != "ICICI EMI 1" {
		t.Fatalf("description %+v", e.Description)
	}
	if e.Tenure == nil || *e.Tenure != 9 {
		t.Fatalf("tenure %+v", e.Tenure)
	}
	if e.MonthsPaid == nil || *e.MonthsPaid != 3 {
		t.Fatalf("monthsPaid %+v", e.MonthsPaid)
	}
	if e.DateStr != "2026-10-07" {
		t.Fatalf("date %s", e.DateStr)
	}
	if e.StartDate != "2026-10-07" {
		t.Fatalf("startDate %s", e.StartDate)
	}
}
