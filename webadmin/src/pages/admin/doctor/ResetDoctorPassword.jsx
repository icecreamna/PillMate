import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import styles from "../../../styles/admin/doctor/ResetDoctorPassword.module.css";
import { getDoctor, resetDoctorPassword } from "../../../services/doctors";

export default function ResetDoctorPassword() {
  const { id } = useParams();
  const nav = useNavigate();

  // แสดงอีเมล/username จริงจาก API (fallback เป็น mock ถ้าโหลดไม่ได้)
  const [email, setEmail] = useState(`doctor_${id}@pillmate.com`);

  const [pw1, setPw1] = useState("");
  const [pw2, setPw2] = useState("");
  const [show, setShow] = useState(false);
  const [loading, setLoading] = useState(false);
  const [msg, setMsg] = useState("");
  const [err, setErr] = useState("");

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const res = await getDoctor(id);   // GET /admin/doctors/:id -> { data:{...} }
        const d = res?.data || {};
        if (!cancelled) setEmail(d.username || d.email || `doctor_${id}@pillmate.com`);
      } catch {
        // เงียบไว้ ใช้ fallback email เดิม
      }
    })();
    return () => { cancelled = true; };
  }, [id]);

  const validate = () => {
    if (!pw1) return "กรุณากรอกรหัสผ่านใหม่";
    if (pw1 !== pw2) return "รหัสผ่านยืนยันไม่ตรงกัน";
    return "";
  };

  const onSubmit = async (e) => {
    e.preventDefault();
    const v = validate();
    if (v) { setErr(v); return; }
    setErr(""); setMsg(""); setLoading(true);
    try {
      // PATCH /admin/doctors/:id/password  -> body: { new_password }
      await resetDoctorPassword(id, pw1);
      setMsg("เปลี่ยนรหัสผ่านสำเร็จ");
      // ถ้าต้องการกลับทันทีหลังสำเร็จ: nav(-1)
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
