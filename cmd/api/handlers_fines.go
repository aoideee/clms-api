package main

import (
	"fmt"
	"net/http"

	"github.com/aoideee/clms-api/internal/data"
)

// createFineHandler processes new fees and fines for a member's account.
func (app *application) createFineHandler(w http.ResponseWriter, r *http.Request) {
	// The Cashier needs to know who to charge, what type of fee it is, and the amount.
	// We use a pointer for LoanID because registration fees are not tied to a specific book!
	var input struct {
		LoanID     *int64  `json:"loan_id"`
		MemberID   int64   `json:"member_id"`
		FineType   string  `json:"fine_type"`
		Amount     float64 `json:"amount"`
		PaidStatus bool    `json:"paid_status"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	fine := &data.Fine{
		LoanID:     input.LoanID,
		MemberID:   input.MemberID,
		FineType:   input.FineType,
		Amount:     input.Amount,
		PaidStatus: input.PaidStatus,
	}

	// Hand the transaction to the database
	err = app.models.Fines.Insert(fine)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/fines/%d", fine.ID))
	totalFinesCreated.Add(1)

	// Return the receipt!
	err = app.writeJSON(w, http.StatusCreated, envelope{"fine": fine}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}