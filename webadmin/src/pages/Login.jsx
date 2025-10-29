// src/pages/Login.jsx
import { useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import styles from "../styles/auth/Login.module.css";
import { login } from "../services/auth";

export default function Login() {
  const nav = useNavigate();
  const { search } = useLocation();
  const params = new URLSearchParams(search);
  const returnTo = params.get("returnTo");

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading,  setLoading]  = useState(false);
  const [error,    setError]    = useState("");

  const onSubmit = async (e) => {
    e.preventDefault();
    setError(""); 
    setLoading(true);

    try {
      // เรียก BE: ตั้งคุกกี้ admin_jwt ที่ฝั่งเซิร์ฟเวอร์
      const data = await login(username, password);
      // โครงตอบกลับคาดหวัง: { user:{..., role, email}, role }
      const role  = data?.role || data?.user?.role || "";
      const email = data?.user?.email || username || "";

      // เก็บเฉพาะข้อมูลโชว์/redirect (ไม่เก็บ token ใน FE เพราะใช้ HttpOnly cookie)
      localStorage.setItem("role", role);
      localStorage.setItem("email", email);

      // มี returnTo มาก็พากลับก่อน
      if (returnTo) {
        nav(returnTo, { replace: true });
        return;
      }

      // เปลี่ยนเส้นทางตามบทบาท
      if (role === "superadmin") {
        nav("/admin", { replace: true });
      } else if (role === "doctor") {
        nav("/doc", { replace: true });
      } else {
        nav("/", { replace: true });
      }
    } catch (err) {
      setError(err?.message || "Login error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.page}>
      <h1 className={styles.heading}>Login</h1>
      <div className={styles.card}>
        <p className={styles.note}>กรุณาเข้าสู่ระบบเพื่อใช้งาน</p>

        {error && <div className={styles.error}>{error}</div>}

        <form onSubmit={onSubmit} className={styles.form}>
          <label className={styles.label}>
            <span>Username</span>
            <input
              className={styles.input}
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              autoComplete="username"
            />
          </label>

          <label className={styles.label}>
            <span>Password</span>
            <input
              className={styles.input}
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete="current-password"
            />
          </label>

          <button className={styles.button} type="submit" disabled={loading}>
            {loading ? "กำลังเข้าสู่ระบบ..." : "Login"}
          </button>
        </form>
      </div>
    </div>
  );
}
