// src/pages/doctor/EditDoctor.jsx
import { useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import styles from "../../../styles/admin/doctor/EditDoctor.module.css";

const mock = [
  {
    id: 1,
    first_name: "Doctor",
    last_name: "A",
    email: "doctor_a@pillmate.com",
  },
  {
    id: 2,
    first_name: "Doctor",
    last_name: "B",
    email: "doctor_b@pillmate.com",
  },
  {
    id: 3,
    first_name: "Doctor",
    last_name: "C",
    email: "doctor_c@pillmate.com",
  },
];

export default function EditDoctor() {
  const { id } = useParams();
  const nav = useNavigate();

  const current = useMemo(
    () => mock.find((m) => String(m.id) === String(id)),
    [id]
  );
  const [firstName, setFirstName] = useState(current?.first_name || "");
  const [lastName, setLastName] = useState(current?.last_name || "");
  const [loading, setLoading] = useState(false);
  const [msg, setMsg] = useState("");
  const [err, setErr] = useState("");

  const save = async (e) => {
    e.preventDefault();
    setErr("");
    setMsg("");
    setLoading(true);
    try {
      // TODO: PUT /api/doctors/:id  (ส่ง {first_name,last_name})
      await new Promise((r) => setTimeout(r, 600));
      setMsg("บันทึกสำเร็จ");
      // nav(-1) // ถ้าต้องการกลับทันที
    } catch (e) {
      setErr(e.message || "บันทึกไม่สำเร็จ");
    } finally {
      setLoading(false);
    }
  };

  const resetPassword = async () => {
    setErr("");
    setMsg("");
    setLoading(true);
    try {
      // TODO: POST /api/doctors/:id/reset-password
      await new Promise((r) => setTimeout(r, 500));
      setMsg("ส่งลิงก์รีเซ็ตรหัสผ่านแล้ว");
    } catch (e) {
      setErr(e.message || "ทำรายการไม่สำเร็จ");
    } finally {
      setLoading(false);
    }
  };

  if (!current) return <div className={styles.page}>ไม่พบข้อมูลแพทย์</div>;

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
        <div className={styles.email}>{current.email}</div>

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
          // ปุ่มไปหน้า reset password
          <button
            className={styles.reset}
            onClick={() => nav(`/admin/${id}/reset-password`)}
          >
            reset password →
          </button>
          <button className={styles.submit} type="submit" disabled={loading}>
            {loading ? "Saving…" : "Save Changes"}
          </button>
        </form>
      </div>
    </div>
  );
}
