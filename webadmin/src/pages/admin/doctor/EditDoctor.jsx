import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import styles from "../../../styles/admin/doctor/EditDoctor.module.css";
import { getDoctor, updateDoctor } from "../../../services/doctors";

export default function EditDoctor() {
  const { id } = useParams();
  const nav = useNavigate();

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [msg, setMsg] = useState("");
  const [err, setErr] = useState("");

  // ฟิลด์ฟอร์ม
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName]   = useState("");
  const [email, setEmail]         = useState(""); // แสดง read-only (จาก username)

  // ดึงข้อมูลหมอตาม id
  useEffect(() => {
    let cancelled = false;
    (async () => {
      setErr(""); setMsg(""); setLoading(true);
      try {
        const res = await getDoctor(id);     // GET /admin/doctors/:id  -> { data: {...} }
        const d   = res?.data || {};
        if (!cancelled) {
          setFirstName(d.first_name || "");
          setLastName(d.last_name || "");
          setEmail(d.username || d.email || ""); // DTO ของ BE ใช้ username เป็นอีเมล
        }
      } catch (e) {
        if (!cancelled) setErr(e.message || "ไม่พบข้อมูลแพทย์");
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => { cancelled = true; };
  }, [id]);

  const save = async (e) => {
    e.preventDefault();
    setErr(""); setMsg(""); setSaving(true);
    try {
      // PUT /admin/doctors/:id  -> body: { username?, first_name?, last_name?, password? }
      const payload = {
        first_name: firstName.trim(),
        last_name: lastName.trim(),
        // ไม่ให้แก้อีเมลจากหน้านี้: ถ้าจะให้แก้ ค่อยเพิ่ม input/ส่ง username ด้วย
      };
      await updateDoctor(id, payload);
      setMsg("บันทึกสำเร็จ");
      // ถ้าต้องการกลับทันที: nav(-1)
    } catch (e) {
      setErr(e.message || "บันทึกไม่สำเร็จ");
    } finally {
      setSaving(false);
    }
  };

  if (loading) return <div className={styles.page}>กำลังโหลด...</div>;
  if (err && !firstName && !lastName) return <div className={styles.page}>{err}</div>;

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

          {/* ปุ่มไปหน้า reset password */}
          <button
            type="button"
            className={styles.reset}
            onClick={() => nav(`/admin/${id}/reset-password`)}
          >
            reset password →
          </button>

          <button className={styles.submit} type="submit" disabled={saving}>
            {saving ? "Saving…" : "Save Changes"}
          </button>
        </form>
      </div>
    </div>
  );
}
