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

// ===== ใหม่ (สำหรับ doctor self-service) =====
// POST /auth/doctor/register  — สมัครหมอเอง (public)
export const registerDoctor = ({ username, password, first_name, last_name }) =>
  request("/auth/doctor/register", {
    method: "POST",
    body: {
      username: String(username || "").trim(),
      password: String(password || "").trim(),
      first_name: String(first_name || "").trim(),
      last_name: String(last_name || "").trim(),
    },
  });