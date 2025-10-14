// src/pages/doctor/patient/PatientList.jsx
import { useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/doctor/patient/PatientList.module.css'
import { listPatients, deletePatient } from '../../../services/patients'

// เก็บเฉพาะตัวเลข (กันกรณีพิมพ์มีขีด/ช่องว่าง)
const normIdCard = (s) => (s || '').replace(/[^\d]/g, '')

// คำนวณอายุจาก birth_day (ISO string)
function calcAge(birthISO) {
  if (!birthISO) return '-'
  const d = new Date(birthISO)
  if (Number.isNaN(d.getTime())) return '-'
  const now = new Date()
  let age = now.getFullYear() - d.getFullYear()
  const m = now.getMonth() - d.getMonth()
  if (m < 0 || (m === 0 && now.getDate() < d.getDate())) age--
  return age < 0 ? '-' : age
}

// แปลง DTO → รูปแบบที่ UI เดิมใช้
function mapPatientDTO(p) {
  return {
    raw: p,
    id: p.id,
    name: [p.first_name, p.last_name].filter(Boolean).join(' ') || '-',
    idcard: p.id_card_number || '-',
    gender: p.gender || '-',
    age: calcAge(p.birth_day),
  }
}

export default function PatientList(){
  const nav = useNavigate()

  const [query, setQuery] = useState('')
  const [allRows, setAllRows] = useState([])   // เก็บทั้งหมดจาก API
  const [selected, setSelected] = useState(null) // เก็บผลค้นหาแบบตรงด้วยเลขบัตร
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // โหลดทั้งหมดครั้งแรก
  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoading(true); setError('')
      try {
        const res = await listPatients() // GET /doctor/hospital-patients
        const list = Array.isArray(res?.data) ? res.data : []
        if (!cancelled) {
          setAllRows(list.map(mapPatientDTO))
          setSelected(null)
        }
      } catch (e) {
        if (!cancelled) setError(e.message || 'โหลดรายชื่อผู้ป่วยไม่สำเร็จ')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => { cancelled = true }
  }, [])

  // ค้นหาด้วยเลขบัตร (normalize ก่อน แล้ว exact match ที่ฝั่ง FE)
  const handleSearch = async () => {
    const q = normIdCard(query.trim())
    if (!q) { setSelected(null); return }
    try {
      const res = await listPatients({ q }) // GET /doctor/hospital-patients?q=...
      const list = Array.isArray(res?.data) ? res.data : []
      // เลือก “เลขบัตรตรงเป๊ะ” ก่อน (normalize ทั้งสองฝั่ง) ถ้าไม่เจอ fallback ตัวแรก
      const exact = list.find(p => normIdCard(String(p.id_card_number)) === q) || list[0]
      if (exact) setSelected(mapPatientDTO(exact))
      else alert('ไม่พบหมายเลขบัตรประชาชนนี้')
    } catch (e) {
      alert(e.message || 'ค้นหาไม่สำเร็จ')
    }
  }

  const onKeyDown = (e) => { if (e.key === 'Enter') handleSearch() }

  const onDelete = async (id) => {
    if (!confirm('ยืนยันลบ?')) return
    try {
      await deletePatient(id)
      // อัปเดตรายการปัจจุบัน
      setAllRows(prev => prev.filter(x => x.id !== id))
      setSelected(prev => (prev && prev.id === id ? null : prev))
    } catch (e) {
      alert(e.message || 'ลบไม่สำเร็จ')
    }
  }

  // รายการที่จะแสดง
  const rows = useMemo(() => {
    return selected ? [selected] : allRows
  }, [selected, allRows])

  return (
    <div>
      <div className={styles.headerRow}>
        <h2 className={styles.title}>Patients</h2>
        <div className={styles.headerActions}>
          <input
            className={styles.search}
            placeholder="ค้นหาหมายเลขบัตรประชาชน"
            value={query}
            onChange={e=>setQuery(e.target.value)}
            onKeyDown={onKeyDown}
          />
          <button className={styles.searchBtn} onClick={handleSearch}>ค้นหา</button>
          <button className={styles.addBtn} onClick={()=>nav('/doc/patients/add')}>+ Add Patient</button>
        </div>
      </div>

      {error && <div className={styles.error}>{error}</div>}
      {loading ? (
        <div className={styles.loading}>กำลังโหลด...</div>
      ) : (
        <div className={styles.tableWrap}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th style={{width:'6%'}}>#</th>
                <th style={{width:'20%'}}>Name</th>
                <th style={{width:'20%'}}>IDCardNumber</th>
                <th style={{width:'12%'}}>Gender</th>
                <th style={{width:'10%'}}>Age</th>
                <th style={{width:'25%'}}># Action</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((p,i)=>(
                <tr key={p.id}>
                  <td>{i+1}</td>
                  <td>{p.name}</td>
                  <td>{p.idcard}</td>
                  <td>{p.gender}</td>
                  <td>{p.age}</td>
                  <td className={styles.actions}>
                    <button className={styles.view} onClick={()=>nav(`/doc/patients/${p.id}`)}>View</button>
                    <button className={styles.edit} onClick={()=>nav(`/doc/patients/${p.id}/edit`)}>Edit</button>
                    <button className={styles.delete} onClick={()=>onDelete(p.id)}>Delete</button>
                  </td>
                </tr>
              ))}
              {rows.length === 0 && (
                <tr><td colSpan={6} style={{textAlign:'center', opacity:.7, padding:'12px'}}>ไม่พบข้อมูล</td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
