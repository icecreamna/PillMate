import { useEffect, useState } from "react";
import { Navigate } from "react-router-dom";
import { me } from "../services/auth";

function Loading() {
  return <div style={{ padding: 16 }}>กำลังตรวจสอบสิทธิ์...</div>;
}

/**
 * ตัวครอบ Route ที่ต้องล็อกอิน
 * - เรียก /admin/me ทุกครั้งที่เข้า เพื่อยืนยันเซสชันจากคุกกี้ HttpOnly
 * - needRole: "superadmin" | "doctor"
 */
export default function RequireAuth({ children, needRole }) {
  const [st, setSt] = useState({ loading: true, allow: false });

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const data = await me();
        const role  = data?.role || data?.user?.role || "";
        const email = data?.user?.email || "";
        if (role)  localStorage.setItem("role", role);
        if (email) localStorage.setItem("email", email);
        if (!cancelled) setSt({ loading: false, allow: !needRole || role === needRole });
      } catch {
        if (!cancelled) setSt({ loading: false, allow: false });
      }
    })();
    return () => { cancelled = true; };
  }, [needRole]);

  if (st.loading) return <Loading />;
  if (!st.allow)  return <Navigate to="/login" replace />;
  return children;
}
