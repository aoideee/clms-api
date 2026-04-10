// Filename: static/js/app.js

import { emitter } from "./core/event-emitter.js";
import { state } from "./core/state.js";
import { DataService } from "./core/data-service.js";
import { renderBooks } from "./features/render-books.js";

// --- OBSERVERS FOR DOWNWARD DATA FLOW (API -> State -> Screen) ---

emitter.on("books:loading", () => {
  state.loading = true;
  state.error = null;
  renderBooks();
});

emitter.on("books:loaded", (books) => {
  state.books = books;
  state.loading = false;

  // Logic to disable "Next" button if we get fewer books than a full page
  state.hasNextPage = books.length === state.pageSize;

  renderBooks();
});

emitter.on("books:error", (msg) => {
  state.error = msg;
  state.loading = false;
  renderBooks();
});

// --- OBSERVER FOR UPWARD DATA FLOW (Screen -> State -> API) ---

emitter.on("books:pageRequested", (newPage) => {
  // 1. Update the state with the new page number
  state.currentPage = newPage;

  // 2. Tell the DataService to go fetch that specific page
  DataService.fetchBooks(newPage);
});

// --- BOOT ---

// Kick off the very first fetch for Page 1 when the app loads
emitter.emit("books:pageRequested", 1);