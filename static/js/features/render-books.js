// Filename: static/js/features/render-books.js

import { state } from "../core/state.js";
import { emitter } from "../core/event-emitter.js";

// Render books to the DOM
export function renderBooks() {
  const app = document.querySelector("#app");
  if (!app) return;

  // 1. Handle Loading
  if (state.loading) {
    app.innerHTML = `<div class="status-loading"><span class="spinner">Loading the CLMS catalog...</span></div>`;
    return;
  }

  // 2. Handle Error
  if (state.error) {
    app.innerHTML = `<div class="status-error">! ${state.error}</div>`;
    return;
  }

  // 3. Handle Empty State
  if (!state.books || state.books.length === 0) {
    app.innerHTML = `<div class="status-empty">No books found in the catalog.</div>`;
    return;
  }

  // 4. The Happy Path - Render the Belizean Literature
  const bookCards = state.books
    .map(
      (book) => `
        <div class="book-card" style="border: 1px solid #ddd; padding: 15px; margin-bottom: 10px; border-radius: 5px; background: white;">
            <h3 style="margin-top: 0; color: #4a90e2;">${book.title}</h3>
            <p style="margin: 5px 0; font-size: 0.9em; color: #666;"><strong>ISBN:</strong> ${book.isbn}</p>
            <p style="margin: 5px 0;">${book.description || "No description available."}</p>
        </div>
    `,
    )
    .join("");

  // 5. Pagination Controls
  const controls = `
        <div class="pagination" style="display: flex; gap: 10px; justify-content: center; margin-top: 20px;">
            <button id="prev-btn" ${state.currentPage === 1 ? "disabled" : ""}>Previous</button>
            <span>Page ${state.currentPage}</span>
            <button id="next-btn" ${!state.hasNextPage ? "disabled" : ""}>Next</button>
        </div>
    `;

  // Write the DOM exactly once
  app.innerHTML = `<div>${bookCards}</div>${controls}`;

  // 6. Attach Event Listeners to the buttons (The "Upward" Flow)
  const prevBtn = document.querySelector("#prev-btn");
  const nextBtn = document.querySelector("#next-btn");

  if (prevBtn) {
    prevBtn.addEventListener("click", () => {
      emitter.emit("books:pageRequested", state.currentPage - 1);
    });
  }

  if (nextBtn) {
    nextBtn.addEventListener("click", () => {
      emitter.emit("books:pageRequested", state.currentPage + 1);
    });
  }
}