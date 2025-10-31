// src/pages/admin/doctor/DoctorList.jsx
import { useEffect, useState } from 'react'
import styles from '../../../styles/admin/doctor/DoctorList.module.css'
import { listDoctors, deleteDoctor } from '../../../services/doctors'

export default function DoctorList() {
  const [rows, setRows] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  const fetchData = async () => {
    setLoading(true); setError("")
    try {
      const res = await listDoctors() // GET /admin/doctors
      const list = Array.isArray(res?.data) ? res.data : []
      setRows(list)
    } catch (err) {
      setError(err.message || "โหลดรายชื่อไม่สำเร็จ")
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [])

  const onDelete = async (id) => {
    if (!confirm('ยืนยันการลบรายการนี้?')) return
    try {
      await deleteDoctor(id) // DELETE /admin/doctors/:id
      setRows(prev => prev.filter(x => x.id !== id))
    } catch (err) {
      alert(err.message || "ลบไม่สำเร็จ")
    }
  }

  const displayName = (d) =>
    [d.first_name, d.last_name].filter(Boolean).join(" ") || d.name || "-"

  return (
    <div>
      {/* หัวข้ออย่างเดียว (เอาปุ่ม Add Doctor ออก) */}
      <div className={styles.headerRow}>
        <h2 className={styles.title}>Doctors</h2>
      </div>

      {error && <div className={styles.error}>{error}</div>}
      {loading ? (
        <div className={styles.loading}>กำลังโหลด...</div>
      ) : (
        <div className={styles.tableWrap}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th style={{ width: '8%'  }}>#</th>
                <th style={{ width: '42%' }}>Name</th>
                <th style={{ width: '40%' }}>Username</th>
                <th style={{ width: '10%' }}># Action</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((d, i) => (
                <tr key={d.id ?? i}>
                  <td>{i + 1}</td>
                  <td>{displayName(d)}</td>
                  <td>{d.username}</td>
                  <td className={styles.actions}>
                    <button
                      className={styles.delete}
                      onClick={() => onDelete(d.id)}
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
              {rows.length === 0 && (
                <tr>
                  <td colSpan={4} style={{ textAlign: 'center', padding: 12 }}>
                    ไม่พบข้อมูล
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
