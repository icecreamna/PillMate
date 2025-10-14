// src/services/patients.js
import { request } from "../lib/http";

const base = "/doctor/hospital-patients";

// list (รองรับ q ค้นหา)
export function listPatients({ q = "" } = {}) {
  const qs = q ? `?q=${encodeURIComponent(q)}` : "";
  return request(`${base}${qs}`); // expected -> { data: [...] }
}

// get one
export function getPatient(id) {
  return request(`${base}/${id}`);
}

// create
export function createPatient(payload) {
  // payload ตัวอย่าง:
  // {
  //   id_card_number, first_name, last_name, phone_number,
  //   birth_day: "2000-01-15T00:00:00+07:00",
  //   gender
  // }
  return request(`${base}`, {
    method: "POST",
    body: payload, // ส่ง object ตรง ๆ
  });
}

// update (partial)
export function updatePatient(id, payload) {
  return request(`${base}/${id}`, {
    method: "PUT",
    body: payload, // ส่ง object ตรง ๆ
  });
}

// delete (soft)
export function deletePatient(id) {
  return request(`${base}/${id}`, { method: "DELETE" });
}
