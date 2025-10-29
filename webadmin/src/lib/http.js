// src/lib/http.js
const API_BASE = "http://localhost:8080"; // เปลี่ยนตอน deploy
const CREDENTIALS = "include";

function isFormData(v) {
  return typeof FormData !== "undefined" && v instanceof FormData;
}

async function parseResponse(res) {
  const text = await res.text();
  if (!text) return null;
  try { return JSON.parse(text); } catch { return { raw: text }; }
}

function buildHeaders(baseHeaders = {}, body) {
  const headers = { ...baseHeaders };
  // ถ้าเป็น FormData ห้ามตั้ง Content-Type เอง ให้ browser ใส่ boundary ให้
  if (isFormData(body)) {
    if ("Content-Type" in headers) delete headers["Content-Type"];
  } else {
    headers["Content-Type"] = headers["Content-Type"] || "application/json";
  }
  return headers;
}

function buildBody(body, headers) {
  if (body == null) return undefined;

  // ส่งตรงถ้าเป็น string หรือ FormData
  if (typeof body === "string" || isFormData(body)) return body;

  // ถ้าเป็น object และ content-type เป็น json → stringify
  const ct = (headers["Content-Type"] || "").toLowerCase();
  if (ct.includes("application/json")) return JSON.stringify(body);

  // อย่างอื่น ปล่อยตามเดิม
  return body;
}

// ยิง HTTP แบบรวมศูนย์ + ดัก 401 ให้เด้ง /login อัตโนมัติ
export async function request(path, { method = "GET", body, headers } = {}) {
  const finalHeaders = buildHeaders(headers, body);
  const finalBody = buildBody(body, finalHeaders);

  let res;
  try {
    res = await fetch(API_BASE + path, {
      method,
      headers: finalHeaders,
      credentials: CREDENTIALS,
      body: finalBody,
    });
  } catch (netErr) {
    console.error("NETWORK ERROR →", netErr);
    throw new Error("ไม่สามารถเชื่อมต่อเซิร์ฟเวอร์ได้");
  }

  // แยก response
  const data = await parseResponse(res);

  if (res.status === 401) {
    // เคลียร์ฝั่ง FE แล้วเด้งไปล็อกอิน (แนบ returnTo ให้ด้วย)
    localStorage.removeItem("role");
    localStorage.removeItem("email");
    const returnTo = window.location.pathname + window.location.search;
    window.location.href = "/login?returnTo=" + encodeURIComponent(returnTo);
    throw new Error(data?.error || data?.message || "unauthorized");
  }

  if (!res.ok) {
    console.error("HTTP", res.status, method, path, "→", data);
    const msg = data?.error || data?.message || data?.raw || res.statusText || "Request failed";
    throw new Error(msg);
  }

  return data;
}
