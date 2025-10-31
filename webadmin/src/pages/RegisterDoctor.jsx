// src/pages/RegisterDoctor.jsx
import { useState } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import styles from "../styles/auth/Login.module.css";
import { registerDoctor } from "../services/auth"; // เขียน service ด้านล่าง

export default function RegisterDoctor() {
  const nav = useNavigate();
  const { search } = useLocation();
  const params = new URLSearchParams(search);
  const returnTo = params.get("returnTo");

  const [form, setForm] = useState({
    username: "",
    password: "",
    confirm_password: "",
    first_name: "",
    last_name: "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError]     = useState("");

  const onChange = (e) => setForm({ ...form, [e.target.name]: e.target.value });

  const onSubmit = async (e) => {
    e.preventDefault();
    setError("");

    // ✅ เช็คให้ password ตรงกับ confirm_password ก่อน
    if (form.password.trim() !== form.confirm_password.trim()) {
      setError("รหัสผ่านและยืนยันรหัสผ่านไม่ตรงกัน");
      return;
    }

    setLoading(true);
    try {
      // ส่งเฉพาะ password จริง (ไม่ส่ง confirm_password)
      await registerDoctor({
        username: form.username,
        password: form.password,
        first_name: form.first_name,
        last_name: form.last_name,
      });
      // สมัครเสร็จ พาไปหน้า Login พร้อม returnTo เดิม (ถ้ามี)
      nav(`/login${returnTo ? `?returnTo=${encodeURIComponent(returnTo)}` : ""}`, { replace: true });
    } catch (err) {
      setError(err?.message || "Register error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.page}>
      <h1 className={styles.heading}>Register (Doctor)</h1>
      <div className={styles.card}>
        {error && <div className={styles.error}>{error}</div>}
        <form onSubmit={onSubmit} className={styles.form}>
          <label className={styles.label}>
            <span>Username</span>
            <input
              name="username"
              className={styles.input}
              value={form.username}
              onChange={onChange}
              autoComplete="username"
            />
          </label>

          <label className={styles.label}>
            <span>Password</span>
            <input
              name="password"
              type="password"
              className={styles.input}
              value={form.password}
              onChange={onChange}
              autoComplete="new-password"
            />
          </label>

          {/* ✅ ช่องยืนยันรหัสผ่าน */}
          <label className={styles.label}>
            <span>Confirm password</span>
            <input
              name="confirm_password"
              type="password"
              className={styles.input}
              value={form.confirm_password}
              onChange={onChange}
              autoComplete="new-password"
            />
          </label>

          <label className={styles.label}>
            <span>First name</span>
            <input
              name="first_name"
              className={styles.input}
              value={form.first_name}
              onChange={onChange}
            />
          </label>

          <label className={styles.label}>
            <span>Last name</span>
            <input
              name="last_name"
              className={styles.input}
              value={form.last_name}
              onChange={onChange}
            />
          </label>

          <button className={styles.button} type="submit" disabled={loading}>
            {loading ? "กำลังสมัคร..." : "Create account"}
          </button>
        </form>
      </div>
    </div>
  );
}
