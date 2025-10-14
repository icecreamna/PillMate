// src/services/prescriptions.js
import { request } from "../lib/http";

const base = "/doctor/prescriptions";

/**
 * สร้างใบสั่งยา
 * payload ตัวอย่าง:
 * {
 *   id_card_number: "1101700234567",
 *   doctor_id: 7,                 // (optional) ไม่ส่งก็ได้ เซิร์ฟเวอร์จะดึงจาก token
 *   items: [
 *     { medicine_info_id: 12, amount_per_time: "1 เม็ด", times_per_day: "2 ครั้ง" },
 *     { medicine_info_id: 45, amount_per_time: "5 ml",  times_per_day: "3 ครั้ง" }
 *   ],
 *   sync_until: "2025-12-31T00:00:00+07:00", // (optional)
 *   app_sync_status: false                   // (optional)
 * }
 */
export function createPrescription(payload) {
  return request(`${base}`, {
    method: "POST",
    body: payload,
  });
}

/** list (รองรับ q และ doctor_id) */
export function listPrescriptions({ q = "", doctor_id } = {}) {
  const qs = new URLSearchParams();
  if (q) qs.set("q", q);
  if (doctor_id) qs.set("doctor_id", String(doctor_id));
  const suffix = qs.toString() ? `?${qs.toString()}` : "";
  return request(`${base}${suffix}`); // -> { data: [...] }
}

/** get one */
export function getPrescription(id) {
  return request(`${base}/${id}`); // -> { data: {... with items} }
}

/** update (เฉพาะหัวเอกสาร) */
export function updatePrescription(id, payload) {
  // payload: { id_card_number?, doctor_id?, sync_until?, app_sync_status? }
  return request(`${base}/${id}`, {
    method: "PUT",
    body: payload,
  });
}

/** delete */
export function deletePrescription(id) {
  return request(`${base}/${id}`, { method: "DELETE" });
}
