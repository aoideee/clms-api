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

emitter.on("books:pageRequested", async (newPage) => {
  // 1. Update the state with the new page number
  state.currentPage = newPage;

  // 2. Tell the UI to show the loading spinner
  emitter.emit("books:loading");

  try {
    // 3. Fetch the actual data from the Go backend
    const data = await DataService.fetchBooks(newPage);

    // 4. Notify everyone that the data is here!
    // (Assuming backend returns { "books": [...] })
    emitter.emit("books:loaded", data.books);
  } catch (err) {
    // 5. Handle any network or server errors
    emitter.emit("books:error", err.message);
  }
});

// --- BOOT ---

// Kick off the very first fetch for Page 1 when the app loads
emitter.emit("books:pageRequested", 1);