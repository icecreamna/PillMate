// src/services/appointments.js
import { request } from "../lib/http";

const base = "/doctor/appointments";

// CREATE
// body:
// {
//   id_card_number: "1101700203452",
//   appointment_date: "YYYY-MM-DD",
//   appointment_time: "HH:mm",
//   note?: "string"
// }
export function createAppointment(payload) {
  return request(`${base}`, { method: "POST", body: payload });
}

// LIST (ใหม่สุดก่อน)
// GET /doctor/appointments?q=&date_from=YYYY-MM-DD&date_to=YYYY-MM-DD
export function listAppointments({ q = "", date_from, date_to } = {}) {
  const qs = new URLSearchParams();
  if (q) qs.set("q", q);
  if (date_from) qs.set("date_from", date_from);
  if (date_to) qs.set("date_to", date_to);
  const suffix = qs.toString() ? `?${qs.toString()}` : "";
  return request(`${base}${suffix}`); // -> { data: [...] }
}

// GET ONE
export function getAppointment(id) {
  return request(`${base}/${id}`); // -> { data: {...} }
}

// UPDATE (partial)
// body: { appointment_date?, appointment_time?, note? }
export function updateAppointment(id, payload) {
  return request(`${base}/${id}`, { method: "PUT", body: payload });
}

// DELETE
export function deleteAppointment(id) {
  return request(`${base}/${id}`, { method: "DELETE" });
}
