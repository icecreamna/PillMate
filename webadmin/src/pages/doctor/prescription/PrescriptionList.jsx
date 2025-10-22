// src/pages/doctor/prescription/PrescriptionList.jsx
import { useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/doctor/prescription/PrescriptionList.module.css'
import { listPatients } from '../../../services/patients'

// คำนวณอายุจาก birth_day (YYYY-MM-DD)
function calcAge(birthYMD){
  if(!birthYMD) return '-'
  const d = new Date(birthYMD)
  if (Number.isNaN(d.getTime())) return '-'
  const now = new Date()
  let age = now.getFullYear() - d.getFullYear()
  const m = now.getMonth() - d.getMonth()
  if (m < 0 || (m === 0 && now.getDate() < d.getDate())) age--
  return age < 0 ? '-' : age
}

// map DTO -> shape ที่ตารางนี้ใช้ (เพิ่ม patient_code)
function mapPatientDTO(p){
  return {
    id: p.id,
    code: p.patient_code || '-',
    name: [p.first_name, p.last_name].filter(Boolean).join(' ') || '-',
    idcard: p.id_card_number || '-',
    gender: p.gender || '-',
    age: calcAge(p.birth_day),
  }
}

const onlyDigits = (s) => (s || '').replace(/[^\d]/g, '')
const toPatientCodeCandidate = (s) => {
  const digits = (s || '').replace(/\D/g, '')
  return digits ? digits.padStart(6, '0') : ''
}

export default function PrescriptionList() {
  const nav = useNavigate()

  const [q, setQ] = useState('')
  const [allRows, setAllRows] = useState([])     // ทั้งหมดจาก API
  const [selected, setSelected] = useState(null) // ผลค้นหาแบบตรง
  const [error, setError] = useState('')

  const [loadingInit, setLoadingInit] = useState(true)  // โหลดครั้งแรก
  const [searching, setSearching]   = useState(false)   // กดค้นหา

  // โหลดรายชื่อผู้ป่วยทั้งหมดครั้งแรก
  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoadingInit(true); setError('')
      try {
        const res = await listPatients() // GET /doctor/hospital-patients
        const list = Array.isArray(res?.data) ? res.data : []
        if (!cancelled) {
          setAllRows(list.map(mapPatientDTO))
          setSelected(null)
        }
      } catch (e) {
        if (!cancelled) setError(e.message || 'โหลดข้อมูลไม่สำเร็จ')
      } finally {
        if (!cancelled) setLoadingInit(false)
      }
    })()
    return () => { cancelled = true }
  }, [])

  // ค้นหา patient_code หรือ เลขบัตร (normalize ก่อน แล้ว exact match ที่ฝั่ง FE)
  const onSearch = async () => {
    setError('')
    const raw = (q || '').trim()
    if (!raw) { setSelected(null); return }

    const idQ = onlyDigits(raw)
    const codeQ = toPatientCodeCandidate(raw)
    const rawUpper = raw.toUpperCase()

    try {
      setSearching(true)
      const res = await listPatients({ q: raw }) // BE จะ filter ตาม q ในหลายฟิลด์
      const list = Array.isArray(res?.data) ? res.data : []

      // ลำดับการหา exact:
      // 1) patient_code เท่ากับค่าที่ผู้ใช้พิมพ์ (เผื่ออนาคตมี prefix)
      // 2) patient_code เท่ากับตัวเลข 6 หลัก (pad ซ้าย)
      // 3) id_card_number เท่ากับตัวเลข (normalize)
      const exact = list.find(p => String(p.patient_code || '').toUpperCase() === rawUpper)
                 || list.find(p => String(p.patient_code || '') === codeQ)
                 || list.find(p => onlyDigits(String(p.id_card_number || '')) === idQ)
      if (!exact) {
        setSelected(null)
        setError('ไม่พบผู้ป่วย')
        return
      }
      setSelected(mapPatientDTO(exact))
    } catch (e) {
      setError(e.message || 'ค้นหาไม่สำเร็จ')
    } finally {
      setSearching(false)
    }
  }

  const onKeyDown = (e) => { if (e.key === 'Enter') { e.preventDefault(); onSearch() } }

  // ล้าง error ทันทีเมื่อผู้ใช้พิมพ์ใหม่
  const onChangeQ = (v) => { setQ(v); if (error) setError('') }

  // แถวที่จะแสดง
  const results = useMemo(() => selected ? [selected] : allRows, [selected, allRows])

  return (
    <div>
      <h2 className={styles.title}>Prescription</h2>

      {/* แถวค้นหา */}
      <div className={styles.searchRow}>
        <div className={styles.searchWrap}>
          <input
            className={styles.searchInput}
            placeholder="ค้นหา Patient Code หรือ เลขบัตรประชาชน"
            value={q}
            onChange={e => onChangeQ(e.target.value)}
            onKeyDown={onKeyDown}
            inputMode="text"
          />
          <span className={styles.searchIcon}>🔍</span>
        </div>
        <button
          className={styles.searchBtn}
          onClick={onSearch}
          disabled={searching || loadingInit}
          title="ค้นหาด้วย Patient Code หรือ เลขบัตรประชาชน แบบตรงตัว"
        >
          {searching ? 'กำลังค้นหา...' : 'ค้นหา'}
        </button>
      </div>

      {error && <div className={styles.error}>{error}</div>}

      {/* ตารางผลลัพธ์ */}
      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th style={{width:'6%'}}>#</th>
              <th style={{width:'24%'}}>Name</th>
              <th style={{width:'16%'}}>Patient Code</th>
              <th style={{width:'22%'}}>IDCardNumber</th>
              <th style={{width:'12%'}}>Gender</th>
              <th style={{width:'10%'}}>Age</th>
              <th style={{width:'15%'}}># Action</th>
            </tr>
          </thead>
          <tbody>
            {loadingInit ? (
              <tr>
                <td colSpan={7} style={{textAlign:'center', color:'#6b7280', height:56}}>
                  กำลังโหลด...
                </td>
              </tr>
            ) : results.length === 0 ? (
              <tr>
                <td colSpan={7} style={{textAlign:'center', color:'#6b7280', height:56}}>
                  ไม่มีข้อมูล
                </td>
              </tr>
            ) : (
              results.map((r, i) => (
                <tr key={r.id}>
                  <td>{i + 1}</td>
                  <td>{r.name}</td>
                  <td>{r.code}</td>
                  <td>{r.idcard}</td>
                  <td>{r.gender}</td>
                  <td>{r.age}</td>
                  <td className={styles.actions}>
                    <button className={styles.viewBtn} onClick={()=>nav(`/doc/prescription/view/${r.id}`)}>view</button>
                    <button className={styles.addBtn}  onClick={()=>nav(`/doc/prescription/add/${r.id}`)}>Add</button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
