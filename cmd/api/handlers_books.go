// Filename: handlers_books.go

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"errors"


	"github.com/aoideee/clms-api/internal/data"
)

// createBookHandler handles POST requests to add a new book to the database.
func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	// We declare an anonymous struct to hold the expected data from the request body.
	var input struct {
		Title           string `json:"title"`
		ISBN            string `json:"isbn"`
		Publisher       string `json:"publisher"`
		PublicationYear int32  `json:"publication_year"`
		MinimumAge      int32  `json:"minimum_age"`
		Description     string `json:"description"`
	}

	// Use our new helper to decode the request body into the input struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// Copy the values from the input struct to a new Book struct.
	book := &data.Book{
		Title:           input.Title,
		ISBN:            input.ISBN,
		Publisher:       input.Publisher,
		PublicationYear: input.PublicationYear,
		MinimumAge:      input.MinimumAge,
		Description:     input.Description,
	}

	// Hand the book to our Model to be inserted into the database!
	err = app.models.Books.Insert(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// When sending a 201 Created response, standard practice is to include a Location
	// header pointing to the URL of the newly created resource.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/books/%d", book.ID))

	// Write the JSON response with a 201 Created status code.
	err = app.writeJSON(w, http.StatusCreated, envelope{"book": book}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showBookHandler handles GET requests to retrieve a single book.
func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path. Go 1.22+ makes this very easy with r.PathValue()
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		app.errorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
		return
	}

	// Call our new Get method
	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send the book back to the visitor as a beautifully formatted JSON response
	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateBookHandler handles PATCH requests to modify an existing book.
func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the ID from the URL
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		app.errorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
		return
	}

	// 2. Fetch the existing book from the database
	book, err := app.models.Books.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// 3. Define our input struct. We use pointers here (*string, *int32) so we can tell
	// the difference between a visitor providing an empty value, and not providing the field at all.
	var input struct {
		Title           *string `json:"title"`
		ISBN            *string `json:"isbn"`
		Publisher       *string `json:"publisher"`
		PublicationYear *int32  `json:"publication_year"`
		MinimumAge      *int32  `json:"minimum_age"`
		Description     *string `json:"description"`
	}

	// 4. Decode the JSON request
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// 5. Check which fields were actually provided in the request, and update the book.
	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.ISBN != nil {
		book.ISBN = *input.ISBN
	}
	if input.Publisher != nil {
		book.Publisher = *input.Publisher
	}
	if input.PublicationYear != nil {
		book.PublicationYear = *input.PublicationYear
	}
	if input.MinimumAge != nil {
		book.MinimumAge = *input.MinimumAge
	}
	if input.Description != nil {
		book.Description = *input.Description
	}

	// 6. Save the updated book back to the database
	err = app.models.Books.Update(book)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// 7. Send the updated book back to the visitor
	err = app.writeJSON(w, http.StatusOK, envelope{"book": book}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteBookHandler handles DELETE requests to remove a book from the catalog.
func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the ID from the URL
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		app.errorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
		return
	}

	// 2. Instruct the model to delete the record
	err = app.models.Books.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// 3. Return a 200 OK status with a simple confirmation message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "book successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}