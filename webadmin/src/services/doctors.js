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

// GET /doctor/me — ดึงข้อมูลหมอจาก token ที่ล็อกอินอยู่
export const getMyDoctor = () =>
  request(`/doctor/me`); // -> { data: {...} }

// PUT /doctor/me — หมอแก้โปรไฟล์ตัวเอง
// payload: { first_name?, last_name?, username? }  (BE จะไม่รับ password ที่ endpoint นี้)
export const updateMyDoctor = (payload = {}) =>
  request(`/doctor/me`, { method: "PUT", body: payload }); // -> { message, data }

// PATCH /doctor/me/password — หมอเปลี่ยนรหัสผ่านตัวเอง
// body: { old_password, new_password }
export const changeMyPassword = (old_password, new_password) =>
  request(`/doctor/me/password`, {
    method: "PATCH",
    body: {
      old_password: String(old_password || ""),
      new_password: String(new_password || ""),
    },
  }); // -> { message, data }

// GET /doctor/doctors/:id — อ่านข้อมูลหมอตาม id (read-only)
export const getDoctorPublic = (id) =>
  request(`/doctor/doctors/${id}`); // -> { data: {...} }
