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

// map DTO -> shape ที่ตารางนี้ใช้
function mapPatientDTO(p){
  return {
    id: p.id,
    name: [p.first_name, p.last_name].filter(Boolean).join(' ') || '-',
    idcard: p.id_card_number || '-',
    gender: p.gender || '-',
    age: calcAge(p.birth_day),
  }
}

const onlyDigits = (s) => (s || '').replace(/[^\d]/g, '')

export default function PrescriptionList() {
  const nav = useNavigate()

  const [q, setQ] = useState('')
  const [allRows, setAllRows] = useState([])     // ทั้งหมดจาก API
  const [selected, setSelected] = useState(null) // ผลค้นหาแบบตรงเลขบัตร
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

  // ค้นหาเลขบัตร “ตรงตัว” ถ้าว่างให้โชว์ทั้งหมด
  const onSearch = async () => {
    setError('')
    const qdigits = onlyDigits(q.trim())
    if (!qdigits) { setSelected(null); return }

    try {
      setSearching(true)
      const res = await listPatients({ q: qdigits }) // BE จะ filter ฝั่งเซิร์ฟเวอร์
      const list = Array.isArray(res?.data) ? res.data : []
      // เลือกคนที่เลขบัตรตรงก่อน ถ้าไม่มีให้ว่าง
      const exact = list.find(p => String(p.id_card_number) === qdigits)
      if (!exact) setError('ไม่พบผู้ป่วย')
      setSelected(exact ? mapPatientDTO(exact) : null)
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
            placeholder="ค้นหาเลขบัตรประชาชน"
            value={q}
            onChange={e => onChangeQ(e.target.value)}
            onKeyDown={onKeyDown}
            inputMode="numeric"
          />
          <span className={styles.searchIcon}>🔍</span>
        </div>
        <button
          className={styles.searchBtn}
          onClick={onSearch}
          disabled={searching || loadingInit}
          title="ค้นหาตามเลขบัตรประชาชนแบบตรงตัว"
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
              <th style={{width:'32%'}}>Name</th>
              <th style={{width:'26%'}}>IDCardNumber</th>
              <th style={{width:'12%'}}>Gender</th>
              <th style={{width:'12%'}}>Age</th>
              <th style={{width:'15%'}}># Action</th>
            </tr>
          </thead>
          <tbody>
            {loadingInit ? (
              <tr>
                <td colSpan={6} style={{textAlign:'center', color:'#6b7280', height:56}}>
                  กำลังโหลด...
                </td>
              </tr>
            ) : results.length === 0 ? (
              <tr>
                <td colSpan={6} style={{textAlign:'center', color:'#6b7280', height:56}}>
                  ไม่มีข้อมูล
                </td>
              </tr>
            ) : (
              results.map((r, i) => (
                <tr key={r.id}>
                  <td>{i + 1}</td>
                  <td>{r.name}</td>
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
