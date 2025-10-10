// src/pages/Login.jsx
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import styles from "../styles/auth/Login.module.css";

export default function Login() {
  const nav = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const onSubmit = async (e) => {
    e.preventDefault();
    // TODO: call API จริงแทนด้านล่าง
    await new Promise((r) => setTimeout(r, 400)); // mock
    // หลัง mock login หรือหลังเรียก API สำเร็จ
    localStorage.setItem("auth_token", "demo-token");
    localStorage.setItem("email", username || "super_admin@pillmate.com");
    localStorage.setItem("role", "superadmin"); // ← ใช้ superadmin
    nav("/admin", { replace: true }); // ← ไปโซนผู้ดูแล
  };

  return (
    <div className={styles.page}>
      <h1 className={styles.heading}>Login</h1>
      <div className={styles.card}>
        <p className={styles.note}>กรุณาเข้าสู่ระบบเพื่อใช้งาน</p>
        <form onSubmit={onSubmit} className={styles.form}>
          <label className={styles.label}>
            <span>Username</span>
            <input
              className={styles.input}
              value={username}
              onChange={(e) => setUsername(e.target.value)}
            />
          </label>
          <label className={styles.label}>
            <span>Password</span>
            <input
              className={styles.input}
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
          </label>
          <button className={styles.button} type="submit">
            Login
          </button>
        </form>
      </div>
    </div>
  );
}
