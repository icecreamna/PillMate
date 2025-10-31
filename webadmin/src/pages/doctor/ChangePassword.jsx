// src/pages/doctor/ChangePassword.jsx
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import styles from "../../styles/admin/doctor/ResetDoctorPassword.module.css"; // reuse สไตล์เดิม
import { getMyDoctor, changeMyPassword } from "../../services/doctors";

export default function ChangePassword() {
  const nav = useNavigate();

  const [email, setEmail] = useState("me@pillmate.com");

  const [oldPw, setOldPw] = useState("");
  const [newPw, setNewPw] = useState("");
  const [confirmPw, setConfirmPw] = useState("");

  const [show, setShow] = useState(false);     // toggle โชว์/ซ่อน ทุกช่องรหัส
  const [loading, setLoading] = useState(false);
  const [msg, setMsg] = useState("");
  const [err, setErr] = useState("");

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const res = await getMyDoctor(); // GET /doctor/me -> { data:{ username,... } }
        const d = res?.data || {};
        if (!cancelled) setEmail(d.username || d.email || "me@pillmate.com");
      } catch {
        /* ใช้ค่า default ต่อไป */
      }
    })();
    return () => { cancelled = true; };
  }, []);

  const validate = () => {
    if (!oldPw.trim()) return "กรุณากรอกรหัสผ่านเดิม";
    if (!newPw.trim()) return "กรุณากรอกรหัสผ่านใหม่";
    if (newPw !== confirmPw) return "รหัสผ่านยืนยันไม่ตรงกัน";
    if (newPw === oldPw) return "รหัสผ่านใหม่ต้องแตกต่างจากรหัสเดิม";
    return "";
  };

  const onSubmit = async (e) => {
    e.preventDefault();
    const v = validate();
    if (v) { setErr(v); return; }
    setErr(""); setMsg(""); setLoading(true);
    try {
      // PATCH /doctor/me/password { old_password, new_password }
      await changeMyPassword(oldPw, newPw);
      setMsg("เปลี่ยนรหัสผ่านสำเร็จ");
      setOldPw(""); setNewPw(""); setConfirmPw("");
    } catch (e) {
      setErr(e?.message || "เปลี่ยนรหัสผ่านไม่สำเร็จ");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h2 className={styles.title}>Change Password</h2>
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
          {/* Old password */}
          <label className={styles.label}>
            <span>Old Password</span>
            <div className={styles.pwWrap}>
              <input
                className={styles.input}
                type={show ? "text" : "password"}
                value={oldPw}
                onChange={(e) => setOldPw(e.target.value)}
                autoComplete="current-password"
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

          {/* New password */}
          <label className={styles.label}>
            <span>New Password</span>
            <div className={styles.pwWrap}>
              <input
                className={styles.input}
                type={show ? "text" : "password"}
                value={newPw}
                onChange={(e) => setNewPw(e.target.value)}
                autoComplete="new-password"
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

          {/* Confirm new password */}
          <label className={styles.label}>
            <span>Confirm New Password</span>
            <input
              className={styles.input}
              type={show ? "text" : "password"}
              value={confirmPw}
              onChange={(e) => setConfirmPw(e.target.value)}
              autoComplete="new-password"
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
