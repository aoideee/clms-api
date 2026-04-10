// Filename: static/js/core/state.js

export const state = {
    // Book data
    books: [],

    // Pagination state
    currentPage: 1,
    pageSize: 10,
    hasNextPage: true,
    hasPreviousPage: false,
    totalPages: 1,
    totalBooks: 0,

    // UI state
    loading: false,
    error: null,
}