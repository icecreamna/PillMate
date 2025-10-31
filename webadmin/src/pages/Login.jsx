// src/pages/Login.jsx
import { useState } from "react";
import { useNavigate, useLocation, Link } from "react-router-dom"; // ✅ เพิ่ม Link
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
      const data = await login(username, password);
      const role  = data?.role || data?.user?.role || "";
      const email = data?.user?.email || username || "";

      localStorage.setItem("role", role);
      localStorage.setItem("email", email);

      if (returnTo) {
        nav(returnTo, { replace: true });
        return;
      }
      if (role === "superadmin") nav("/admin", { replace: true });
      else if (role === "doctor") nav("/doc", { replace: true });
      else nav("/", { replace: true });
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

        {/* ✅ ปุ่ม Register for Doctor */}
        <div className={styles.actions} style={{ marginTop: 12, textAlign: "center" }}>
          <span className={styles.muted}>ยังไม่มีบัญชีหมอ?</span>{" "}
          <Link
            to={`/register/doctor${returnTo ? `?returnTo=${encodeURIComponent(returnTo)}` : ""}`}
            className={styles.linkButton /* ถ้ายังไม่มี คลาสนี้ให้ใช้ styles.button แทน */}
          >
            Register for Doctor
          </Link>
        </div>
      </div>
    </div>
  );
}
