// Filename: handlers_loans.go

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aoideee/clms-api/internal/data"
)

// createLoanHandler handles the checkout process.
func (app *application) createLoanHandler(w http.ResponseWriter, r *http.Request) {
	// The scanner only needs two pieces of information: Who is borrowing, and What copy.
	var input struct {
		CopyID   int64 `json:"copy_id"`
		MemberID int64 `json:"member_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// BNLSIS Business Logic: Standard checkout period is 14 days!
	// We dynamically calculate the due date starting from right now.
	dueDate := time.Now().AddDate(0, 0, 14)

	loan := &data.Loan{
		CopyID:   input.CopyID,
		MemberID: input.MemberID,
		DueDate:  dueDate,
	}

	// Hand the form to the database to officially record the checkout
	err = app.models.Loans.Insert(loan)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Note: In a fully fleshed-out system, we would also run a query here to 
	// update the Copy's status from "Available" to "Checked Out".

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/loans/%d", loan.ID))
	totalBooksLoaned.Add(1)

	// Return the receipt to the patron!
	err = app.writeJSON(w, http.StatusCreated, envelope{"loan": loan}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}