import { request } from "../lib/http";

// POST /admin/login
export const login = (username, password) =>
  request("/admin/login", {
    method: "POST",
    body: { username: username.trim(), password: password.trim() },
  });

// POST /admin/logout
export const logout = () =>
  request("/admin/logout", { method: "POST" });

// GET /admin/me
export const me = () =>
  request("/admin/me");
