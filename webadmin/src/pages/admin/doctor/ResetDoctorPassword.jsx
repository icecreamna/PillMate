import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import styles from "../../../styles/admin/doctor/ResetDoctorPassword.module.css";

export default function ResetDoctorPassword() {
  const { id } = useParams();
  const nav = useNavigate();
  // ในของจริงโหลดอีเมลจาก API ด้วย id; ตอนนี้ mock ไว้
  const email = `doctor_${id}@pillmate.com`;

  const [pw1, setPw1] = useState("");
  const [pw2, setPw2] = useState("");
  const [show, setShow] = useState(false);
  const [loading, setLoading] = useState(false);
  const [msg, setMsg] = useState("");
  const [err, setErr] = useState("");

  const validate = () => {
    if (!pw1) return "กรุณากรอกรหัสผ่านใหม่";
    if (pw1 !== pw2) return "รหัสผ่านยืนยันไม่ตรงกัน";
    return "";
  };

  const onSubmit = async (e) => {
    e.preventDefault();
    const v = validate();
    if (v) {
      setErr(v);
      return;
    }
    setErr("");
    setMsg("");
    setLoading(true);
    try {
      // TODO: เรียก API จริง:
      // await fetch(`/api/doctors/${id}/password`, { method:'PUT', headers:{'Content-Type':'application/json'}, body: JSON.stringify({ password: pw1 }) })
      await new Promise((r) => setTimeout(r, 600)); // mock
      setMsg("เปลี่ยนรหัสผ่านสำเร็จ");
      // nav(-1) // ถ้าต้องการย้อนกลับอัตโนมัติ
    } catch (e) {
      setErr(e.message || "เปลี่ยนรหัสผ่านไม่สำเร็จ");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h2 className={styles.title}>Edit Doctor</h2>
        <button className={styles.back} onClick={() => nav(-1)}>
          ← Back
        </button>
      </div>
      <hr className={styles.hr} />

      <div className={styles.card}>
        <div className={styles.email}>
          <button
            className={styles.circle}
            onClick={() => nav(-1)}
            aria-label="back"
          >
            ←
          </button>
          <span>{email}</span>
        </div>

        {err && <div className={styles.error}>{err}</div>}
        {msg && <div className={styles.success}>{msg}</div>}

        <form onSubmit={onSubmit} className={styles.form}>
          <label className={styles.label}>
            <span>New Password</span>
            <div className={styles.pwWrap}>
              <input
                className={styles.input}
                type={show ? "text" : "password"}
                value={pw1}
                onChange={(e) => setPw1(e.target.value)}
              />
              <button
                type="button"
                className={styles.toggle}
                onClick={() => setShow((s) => !s)}
              >
                {show ? "Hide" : "Show"}
              </button>
            </div>
          </label>

          <label className={styles.label}>
            <span>Confirm New Password</span>
            <input
              className={styles.input}
              type={show ? "text" : "password"}
              value={pw2}
              onChange={(e) => setPw2(e.target.value)}
            />
          </label>

          <button className={styles.submit} type="submit" disabled={loading}>
            {loading ? "Processing…" : "Change Password"}
          </button>
        </form>
      </div>
    </div>
  );
}
