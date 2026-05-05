// Filename: static/js/features/render-books.js

import { state } from "../core/state.js";
import { emitter } from "../core/event-emitter.js";

// Render books to the DOM
export function renderBooks() {
  const app = document.querySelector("#app");
  if (!app) return;

  // 1. Handle Loading
  if (state.loading) {
    app.innerHTML = `<div class="status-loading"><span class="spinner">Loading the CLMS catalog…</span></div>`;
    return;
  }

  // 2. Handle Error
  if (state.error) {
    app.innerHTML = `<div class="status-error">&#9888; ${state.error}</div>`;
    return;
  }

  // 3. Handle Empty State
  if (!state.books || state.books.length === 0) {
    app.innerHTML = `<div class="status-empty">No books found in the catalog.</div>`;
    return;
  }

  // 4. The Happy Path — Render the book catalog
  const bookCards = state.books
    .map(
      (book) => `
        <div class="book-card">
            <h3>${book.title}</h3>
            <p class="book-meta"><strong>ISBN:</strong> ${book.isbn}</p>
            <p class="book-desc">${book.description || "No description available."}</p>
        </div>
    `,
    )
    .join("");

  // 5. Pagination Controls
  const controls = `
        <div class="pagination">
            <button id="prev-btn" class="btn-page" ${state.currentPage === 1 ? "disabled" : ""}>← Previous</button>
            <span class="pagination-label">Page ${state.currentPage} of ${state.totalPages}</span>
            <button id="next-btn" class="btn-page" ${!state.hasNextPage ? "disabled" : ""}>Next →</button>
        </div>
    `;

  // Write the DOM exactly once
  app.innerHTML = `<div>${bookCards}</div>${controls}`;

  // 6. Attach Event Listeners (Upward Flow)
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