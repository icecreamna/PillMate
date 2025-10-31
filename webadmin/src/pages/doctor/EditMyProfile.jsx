// src/pages/doctor/EditMyProfile.jsx
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import styles from "../../styles/admin/doctor/EditDoctor.module.css"; // reuse สไตล์เดิม
import { getMyDoctor, updateMyDoctor } from "../../services/doctors";

export default function EditMyProfile() {
  const nav = useNavigate();

  const [loading, setLoading] = useState(true);
  const [saving,  setSaving]  = useState(false);
  const [msg,     setMsg]     = useState("");
  const [err,     setErr]     = useState("");

  const [firstName, setFirstName] = useState("");
  const [lastName,  setLastName]  = useState("");
  const [email,     setEmail]     = useState(""); // read-only (username/email จาก BE)

  useEffect(() => {
    let cancelled = false;
    (async () => {
      setErr(""); setMsg(""); setLoading(true);
      try {
        // GET /doctor/me -> { data: { first_name, last_name, username ... } }
        const res = await getMyDoctor();
        const d   = res?.data || {};
        if (!cancelled) {
          setFirstName(d.first_name || "");
          setLastName(d.last_name || "");
          setEmail(d.username || d.email || "");
        }
      } catch (e) {
        if (!cancelled) setErr(e?.message || "ไม่พบข้อมูลผู้ใช้");
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => { cancelled = true; };
  }, []);

  const save = async (e) => {
    e.preventDefault();
    setErr(""); setMsg(""); setSaving(true);
    try {
      // PUT /doctor/me -> { username?, first_name?, last_name? } (เราไม่แก้ username ที่หน้านี้)
      await updateMyDoctor({
        first_name: firstName.trim(),
        last_name : lastName.trim(),
      });
      setMsg("บันทึกสำเร็จ");
    } catch (e) {
      setErr(e?.message || "บันทึกไม่สำเร็จ");
    } finally {
      setSaving(false);
    }
  };

  if (loading) return <div className={styles.page}>กำลังโหลด...</div>;
  if (err && !firstName && !lastName) return <div className={styles.page}>{err}</div>;

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h2 className={styles.title}>Edit My Profile</h2>
        <button className={styles.back} onClick={() => nav(-1)}>
          ← Back
        </button>
      </div>
      <hr className={styles.hr} />

      <div className={styles.card}>
        {/* แสดงอีเมล/ยูสเซอร์เนม read-only */}
        <div className={styles.email}>{email}</div>

        {err && <div className={styles.error}>{err}</div>}
        {msg && <div className={styles.success}>{msg}</div>}

        <form onSubmit={save} className={styles.form}>
          <label className={styles.label}>
            <span>First Name</span>
            <input
              className={styles.input}
              value={firstName}
              onChange={(e) => setFirstName(e.target.value)}
            />
          </label>

          <label className={styles.label}>
            <span>Last Name</span>
            <input
              className={styles.input}
              value={lastName}
              onChange={(e) => setLastName(e.target.value)}
            />
          </label>

          {/* ปุ่มไปหน้าเปลี่ยนรหัสผ่านของหมอเอง */}
          <button
            type="button"
            className={styles.reset}
            onClick={() => nav(`/doc/change-password`)}
          >
            change password →
          </button>

          <button className={styles.submit} type="submit" disabled={saving}>
            {saving ? "Saving…" : "Save Changes"}
          </button>
        </form>
      </div>
    </div>
  );
}
