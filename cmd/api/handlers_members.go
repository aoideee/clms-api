// Filename: handlers_members.go

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"errors"

	"github.com/aoideee/clms-api/internal/data"
	"github.com/aoideee/clms-api/internal/validator"
)

// createMemberHandler handles POST requests to register a new library member.
func (app *application) createMemberHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		DOB           string `json:"dob"`
		PhoneNumber   string `json:"phone_number"`
		Email         string `json:"email"`
		Address       string `json:"address"`
		AccountStatus string `json:"account_status"`
	}

	err := app.readJSON(r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	member := &data.Member{
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		DOB:           input.DOB,
		PhoneNumber:   input.PhoneNumber,
		Email:         input.Email,
		Address:       input.Address,
		AccountStatus: input.AccountStatus,
	}

	// 1. Initialize a new Validator instance
	v := validator.New()

	// 2. Run the strict BNLSIS checks on the new application
	if data.ValidateMember(v, member); !v.Valid() {
		// If the form is invalid, hand it back to the patron with the exact errors!
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// 3. If the form is perfect, insert the member into the database
	err = app.models.Members.Insert(member)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/members/%d", member.ID))
	totalMembersRegistered.Add(1)

	err = app.writeJSON(w, http.StatusCreated, envelope{"member": member}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMemberHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		app.errorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
		return
	}

	member, err := app.models.Members.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"member": member}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMemberHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		app.errorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
		return
	}

	member, err := app.models.Members.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.errorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		FirstName     *string `json:"first_name"`
		LastName      *string `json:"last_name"`
		DOB           *string `json:"dob"`
		PhoneNumber   *string `json:"phone_number"`
		Email         *string `json:"email"`
		Address       *string `json:"address"`
		AccountStatus *string `json:"account_status"`
	}

	err = app.readJSON(r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if input.FirstName != nil {
		member.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		member.LastName = *input.LastName
	}
	if input.DOB != nil {
		member.DOB = *input.DOB
	}
	if input.PhoneNumber != nil {
		member.PhoneNumber = *input.PhoneNumber
	}
	if input.Email != nil {
		member.Email = *input.Email
	}
	if input.Address != nil {
		member.Address = *input.Address
	}
	if input.AccountStatus != nil {
		member.AccountStatus = *input.AccountStatus
	}

	// Re-run the Inspector to make sure the updates didn't break our rules!
	v := validator.New()
	if data.ValidateMember(v, member); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Members.Update(member)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"member": member}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMemberHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		app.errorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
		return
	}

	err = app.models.Members.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "member successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}