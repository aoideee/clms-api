// Filename: static/js/core/data-service.js

import { emitter } from "./event-emitter.js";
import { state } from "./state.js";

const API_BASE = "http://localhost:4000/v1";

// Data service for making API requests
export const DataService = {
  // Helper function for making API requests
  async request(endpoint, options = {}) {
    const url = `${API_BASE}/${endpoint}`;

    // Add Authorization header if we have a token in state
    if (state.token) {
      options.headers = {
        ...options.headers,
        "Authorization": `Bearer ${state.token}`
      };
    }

    const response = await fetch(url, options);
    if (!response.ok) {
      throw new Error(`API request failed: ${response.statusText}`);
    }
    return response.json();
  },

  // Fetch all books with pagination
  async fetchBooks(page = 1, limit = 10) {
    try {
      const data = await this.request(`books?page=${page}&limit=${limit}`);
      return data;
    } catch (error) {
      console.error("Error fetching books:", error);
      throw error;
    }
  },

  // Login a user and retrieve a token
  async login(email, password) {
    try {
      const data = await this.request("tokens/authentication", {
        method: "POST",
        body: JSON.stringify({ email, password }),
        headers: { "Content-Type": "application/json" }
      });
      return data;
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    }
  }
};