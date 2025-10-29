import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/admin/doctor/DoctorList.module.css'
import { listDoctors, deleteDoctor } from '../../../services/doctors'

export default function DoctorList() {
  const navigate = useNavigate()
  const [rows, setRows] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  const fetchData = async () => {
    setLoading(true); setError("")
    try {
      const res = await listDoctors() // GET /admin/doctors
      // BE คืน { data: [...] }
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
      await deleteDoctor(id)           // DELETE /admin/doctors/:id
      setRows(prev => prev.filter(x => x.id !== id)) // ลบออกจากตารางทันที
    } catch (err) {
      alert(err.message || "ลบไม่สำเร็จ")
    }
  }

  // helper รวมชื่อ (รองรับ DTO ที่เป็น first_name/last_name หรือ name)
  const displayName = (d) =>
    [d.first_name, d.last_name].filter(Boolean).join(" ") || d.name || "-"

  return (
    <div>
      {/* หัวข้อ + ปุ่มอยู่บรรทัดเดียวกัน (สไตล์เดิม) */}
      <div className={styles.headerRow}>
        <h2 className={styles.title}>Doctors</h2>
        <button
          className={styles.addBtn}
          onClick={() => navigate('/admin/add')}
        >
          + Add Doctor
        </button>
      </div>

      {error && <div className={styles.error}>{error}</div>}
      {loading ? (
        <div className={styles.loading}>กำลังโหลด...</div>
      ) : (
        <div className={styles.tableWrap}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th style={{ width: '5%' }}>#</th>
                <th style={{ width: '25%' }}>Name</th>
                <th style={{ width: '40%' }}>Username</th>
                <th style={{ width: '20%' }}>Password</th>
                <th style={{ width: '10%' }}># Action</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((d, i) => (
                <tr key={d.id ?? i}>
                  <td>{i + 1}</td>
                  <td>{displayName(d)}</td>
                  <td>{d.username}</td>
                  <td>{"**********"}</td>
                  <td className={styles.actions}>
                    <button
                      className={styles.edit}
                      onClick={() => navigate(`/admin/${d.id}/edit`)}
                    >
                      Edit
                    </button>
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
                  <td colSpan={5} style={{ textAlign: 'center', padding: 12 }}>
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
