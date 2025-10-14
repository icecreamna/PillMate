// src/services/doctors.js
import { request } from "../lib/http";

/** ===== Admin endpoints ===== */

// GET /admin/doctors?q=&page=&page_size=
export const listDoctors = ({ q = "", page = 1, page_size = 20 } = {}) => {
  const params = new URLSearchParams({ q, page: String(page), page_size: String(page_size) });
  return request(`/admin/doctors?${params.toString()}`); // -> { data: [...] }
};

// GET /admin/doctors/:id
export const getDoctor = (id) =>
  request(`/admin/doctors/${id}`); // -> { data: {...} }

// POST /admin/doctors
// body: { username, password, first_name, last_name }
export const createDoctor = (payload) =>
  request(`/admin/doctors`, { method: "POST", body: payload }); // -> { message, data }

// PUT /admin/doctors/:id
// body: { username?, first_name?, last_name?, password? }
export const updateDoctor = (id, payload) =>
  request(`/admin/doctors/${id}`, { method: "PUT", body: payload }); // -> { message, data }

// DELETE /admin/doctors/:id
export const deleteDoctor = (id) =>
  request(`/admin/doctors/${id}`, { method: "DELETE" }); // -> { message }

// PATCH /admin/doctors/:id/password
// body: { new_password }
export const resetDoctorPassword = (id, new_password) =>
  request(`/admin/doctors/${id}/password`, { method: "PATCH", body: { new_password } }); // -> { message, data }

/** ===== Doctor (public/role=doctor) endpoints ===== */

// GET /doctor/me  — ดึงข้อมูลหมอจาก token ที่ล็อกอินอยู่
export const getMyDoctor = () =>
  request(`/doctor/me`); // -> { data: {...} }

// GET /doctor/doctors/:id — อ่านข้อมูลหมอตาม id (read-only)
export const getDoctorPublic = (id) =>
  request(`/doctor/doctors/${id}`); // -> { data: {...} }
