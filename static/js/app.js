// Filename: static/js/app.js

import { emitter } from "./core/event-emitter.js";
import { state } from "./core/state.js";
import { DataService } from "./core/data-service.js";
import { renderBooks } from "./features/render-books.js";
import { renderLogin, renderHeader } from "./features/auth.js";

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

// --- OBSERVERS FOR AUTH DATA FLOW ---

emitter.on("auth:loginRequested", async ({ email, password }) => {
  state.loading = true;
  state.error = null;
  renderLogin();

  try {
    const data = await DataService.login(email, password);
    
    // Save token to state and storage
    state.token = data.authentication_token.token;
    localStorage.setItem("clms_token", state.token);
    
    state.loading = false;
    
    // Update header to show logout button
    renderHeader();
    
    // Switch to books view
    emitter.emit("books:pageRequested", 1);
  } catch (err) {
    state.loading = false;
    state.error = "Invalid email or password."; 
    renderLogin();
  }
});

emitter.on("auth:logoutRequested", () => {
  // Clear everything
  state.token = null;
  state.books = [];
  localStorage.removeItem("clms_token");

  // Re-render
  renderHeader();
  renderLogin();
});

// --- BOOT ---

// Initialize header
renderHeader();

// If we don't have a token, show login. Otherwise, fetch books.
if (!state.token) {
  renderLogin();
} else {
  emitter.emit("books:pageRequested", 1);
}