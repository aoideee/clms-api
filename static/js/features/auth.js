// Filename: static/js/features/auth.js

import { state } from "../core/state.js";
import { emitter } from "../core/event-emitter.js";

export function renderLogin() {
  const app = document.querySelector("#app");
  if (!app) return;

  app.innerHTML = `
    <div class="login-container">
        <h2>Welcome Back</h2>
        <p class="login-subtitle">Access the National Heritage Library</p>

        ${state.error ? `<div class="alert-error">&#9888; ${state.error}</div>` : ""}

        <form id="login-form">
            <div class="form-group">
                <label class="form-label" for="login-email">Email Address</label>
                <input type="email" id="login-email" class="form-input" required placeholder="e.g. user@clms.bz">
            </div>
            <div class="form-group">
                <label class="form-label" for="login-password">Password</label>
                <input type="password" id="login-password" class="form-input" required placeholder="••••••••">
            </div>
            <button type="submit" class="btn-primary" ${state.loading ? "disabled" : ""}>
                ${state.loading ? "Authenticating…" : "Sign In"}
            </button>
        </form>
    </div>
  `;

  const form = document.querySelector("#login-form");
  if (form) {
    form.addEventListener("submit", (e) => {
      e.preventDefault();
      const email = document.querySelector("#login-email").value;
      const password = document.querySelector("#login-password").value;
      emitter.emit("auth:loginRequested", { email, password });
    });
  }
}

export function renderHeader() {
  const container = document.querySelector("#header-actions");
  if (!container) return;

  if (!state.token) {
    container.innerHTML = "";
    return;
  }

  container.innerHTML = `
    <button id="logout-btn">Sign Out</button>
  `;

  const btn = document.querySelector("#logout-btn");
  if (btn) {
    btn.addEventListener("click", () => {
      emitter.emit("auth:logoutRequested");
    });
  }
}
