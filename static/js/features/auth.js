// Filename: static/js/features/auth.js

import { state } from "../core/state.js";
import { emitter } from "../core/event-emitter.js";

export function renderLogin() {
  const app = document.querySelector("#app");
  if (!app) return;

  app.innerHTML = `
    <div class="login-container" style="max-width: 400px; margin: 60px auto; padding: 40px; background: white; border-radius: 12px; box-shadow: 0 10px 25px rgba(0,0,0,0.05); animation: fadeIn 0.4s ease-out;">
        <h2 style="margin-top: 0; color: #4a90e2; text-align: center; font-weight: 700;">Welcome Back</h2>
        <p style="text-align: center; color: #666; margin-bottom: 30px;">Access the National Heritage Library</p>
        
        ${state.error ? `<div style="background: #fff5f5; color: #e53e3e; padding: 12px; border-radius: 8px; margin-bottom: 20px; font-size: 0.9em; border: 1px solid #fed7d7; animation: shake 0.4s ease-in-out;">! ${state.error}</div>` : ""}

        <form id="login-form">
            <div style="margin-bottom: 20px;">
                <label style="display: block; margin-bottom: 8px; font-weight: 600; font-size: 0.85em; text-transform: uppercase; letter-spacing: 0.05em; color: #666;">Email Address</label>
                <input type="email" id="login-email" required style="width: 100%; padding: 12px; border: 1.5px solid #eee; border-radius: 8px; box-sizing: border-box; font-size: 1em; transition: border-color 0.2s;" placeholder="e.g. user@clms.bz">
            </div>
            <div style="margin-bottom: 30px;">
                <label style="display: block; margin-bottom: 8px; font-weight: 600; font-size: 0.85em; text-transform: uppercase; letter-spacing: 0.05em; color: #666;">Password</label>
                <input type="password" id="login-password" required style="width: 100%; padding: 12px; border: 1.5px solid #eee; border-radius: 8px; box-sizing: border-box; font-size: 1em; transition: border-color 0.2s;" placeholder="••••••••">
            </div>
            <button type="submit" ${state.loading ? "disabled" : ""} style="width: 100%; padding: 14px; background: #4a90e2; color: white; border: none; border-radius: 8px; font-weight: 600; cursor: pointer; transition: all 0.2s; font-size: 1em; box-shadow: 0 4px 12px rgba(74, 144, 226, 0.2);">
                ${state.loading ? "Authenticating..." : "Sign In"}
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
    <button id="logout-btn" style="padding: 8px 16px; background: rgba(255,255,255,0.2); color: white; border: 1px solid rgba(255,255,255,0.4); border-radius: 6px; cursor: pointer; font-size: 0.85em; font-weight: 600; transition: all 0.2s; backdrop-filter: blur(4px);">
        Sign Out
    </button>
  `;

  const btn = document.querySelector("#logout-btn");
  if (btn) {
    btn.addEventListener("click", () => {
      emitter.emit("auth:logoutRequested");
    });
    
    btn.onmouseover = () => btn.style.background = "rgba(255,255,255,0.3)";
    btn.onmouseout = () => btn.style.background = "rgba(255,255,255,0.2)";
  }
}
