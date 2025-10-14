import { request } from "../lib/http";

// LIST: GET /admin/medicine-infos
export const listMedicines = () =>
  request(`/admin/medicine-infos`); // -> { data: [...] }

// GET ONE: GET /admin/medicine-info/:id
export const getMedicine = (id) =>
  request(`/admin/medicine-info/${id}`); // -> { data: {...} }

// CREATE: POST /admin/medicine-info
export const createMedicine = (payload) =>
  request(`/admin/medicine-info`, { method: "POST", body: payload }); // -> { message, data }

// UPDATE: PUT /admin/medicine-info/:id
export const updateMedicine = (id, payload) =>
  request(`/admin/medicine-info/${id}`, { method: "PUT", body: payload }); // -> { message, data }

// DELETE: DELETE /admin/medicine-info/:id
export const deleteMedicine = (id) =>
  request(`/admin/medicine-info/${id}`, { method: "DELETE" }); // -> { message }

// ----- Doctor (ใหม่) -----
export function listMedicinesForDoctor() {
  return request("/doctor/medicine-infos");
}
export function getMedicineForDoctor(id) {
  return request(`/doctor/medicine-info/${id}`);
}