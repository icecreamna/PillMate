import { useState } from "react";
import { useNavigate } from "react-router-dom";
import styles from "../../../styles/admin/doctor/AddDoctor.module.css";
import { createDoctor } from "../../../services/doctors";

export default function AddDoctor() {
  const nav = useNavigate();
  const [form, setForm] = useState({
    first_name: "",
    last_name: "",
    username: "",
    password: "",
    confirm: "",
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const onChange = (k, v) => setForm((s) => ({ ...s, [k]: v }));

  const validate = () => {
    if (!form.first_name.trim()) return "กรุณากรอก First Name";
    if (!form.last_name.trim()) return "กรุณากรอก Last Name";
    if (!form.username.trim()) return "กรุณากรอก Username (อีเมล)";
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.username.trim()))
      return "รูปแบบอีเมลไม่ถูกต้อง";
    if (form.password == null || form.password === "")
      return "กรุณากรอก Password";
    if (form.password !== form.confirm) return "ยืนยันรหัสผ่านไม่ตรงกัน";
    return "";
  };

  const onSubmit = async (e) => {
    e.preventDefault();
    const v = validate();
    if (v) {
      setError(v);
      return;
    }
    setError("");
    setLoading(true);
    try {
      // เรียก API จริง: POST /admin/doctors
      await createDoctor({
        first_name: form.first_name.trim(),
        last_name: form.last_name.trim(),
        username: form.username.trim(),
        password: form.password, // ส่งจริงไปให้ BE hash เอง
      });
      // บันทึกสำเร็จ → กลับไปหน้า list ของ admin
      nav("/admin", { replace: true });
    } catch (err) {
      setError(err.message || "บันทึกไม่สำเร็จ");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h2 className={styles.title}>Add Doctor</h2>
        <button className={styles.back} onClick={() => nav(-1)}>
          ← Back
        </button>
      </div>
      <hr className={styles.hr} />

      <div className={styles.card}>
        {error && <div className={styles.error}>{error}</div>}

        <form onSubmit={onSubmit} className={styles.form}>
          <label className={styles.label}>
            <span>First Name</span>
            <input
              className={styles.input}
              value={form.first_name}
              onChange={(e) => onChange("first_name", e.target.value)}
            />
          </label>

          <label className={styles.label}>
            <span>Last Name</span>
            <input
              className={styles.input}
              value={form.last_name}
              onChange={(e) => onChange("last_name", e.target.value)}
            />
          </label>

          <label className={styles.label}>
            <span>Username</span>
            <input
              className={styles.input}
              value={form.username}
              onChange={(e) => onChange("username", e.target.value)}
              placeholder="doctor_x@pillmate.com"
            />
          </label>

          <label className={styles.label}>
            <span>Password</span>
            <input
              type="password"
              className={styles.input}
              value={form.password}
              onChange={(e) => onChange("password", e.target.value)}
            />
          </label>

          <label className={styles.label}>
            <span>Confirm Password</span>
            <input
              type="password"
              className={styles.input}
              value={form.confirm}
              onChange={(e) => onChange("confirm", e.target.value)}
            />
          </label>

          <button className={styles.submit} type="submit" disabled={loading}>
            {loading ? "Saving…" : "Add Doctor"}
          </button>
        </form>
      </div>
    </div>
  );
}
