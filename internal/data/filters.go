// Filename: internal/data/filters.go

package data

import (
	"math"
)

// Filters holds the strict rules for how we search and organize our library records.
type Filters struct {
	Page         int      // Which page of results the visitor wants
	PageSize     int      // How many records to show per page
	Sort         string   // The column we are sorting by (e.g., "title" or "-title" for descending)
	SortSafelist []string // A strict list of columns we allow visitors to sort by
}

// Metadata holds the pagination details so the visitor knows exactly where they are in the catalog.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// calculateMetadata is our mathematical helper. It figures out how many pages exist in total.
func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{} // Return an empty struct if there are no books
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}