// src/pages/doctor/appointment/AppointmentList.jsx
import { useState, useEffect, useCallback, useMemo } from 'react'
import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/doctor/appointment/AppointmentList.module.css'
import { listPatients } from '../../../services/patients'
import { createAppointment } from '../../../services/appointments'

// helper
function todayStr() {
  const d = new Date()
  const yyyy = d.getFullYear()
  const mm = String(d.getMonth() + 1).padStart(2, '0')
  const dd = String(d.getDate()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd}`
}
function calcAge(birthYMD) {
  if (!birthYMD) return '-'
  const d = new Date(birthYMD)
  if (Number.isNaN(d.getTime())) return '-'
  const now = new Date()
  let age = now.getFullYear() - d.getFullYear()
  const m = now.getMonth() - d.getMonth()
  if (m < 0 || (m === 0 && now.getDate() < d.getDate())) age--
  return age < 0 ? '-' : age
}
const onlyDigits = (s) => (s || '').replace(/[^\d]/g, '')
const toPatientCodeCandidate = (s) => {
  const digits = (s || '').replace(/\D/g, '')
  return digits ? digits.padStart(6, '0') : ''
}

function mapPatientDTO(p){
  return {
    id: p.id,
    code: p.patient_code || '-',
    name: [p.first_name, p.last_name].filter(Boolean).join(' ') || '-',
    idcard: p.id_card_number || '-',
    gender: p.gender || '-',
    age: calcAge(p.birth_day),
    raw: p,
  }
}

export default function AppointmentList() {
  const nav = useNavigate()

  // search state
  const [q, setQ] = useState('')
  const [allRows, setAllRows] = useState([])
  const [selectedRows, setSelectedRows] = useState([]) // ผลลัพธ์ที่แสดง (กรองแล้ว)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // modal state
  const [open, setOpen] = useState(false)
  const [selected, setSelected] = useState(null) // patient
  const [dateVal, setDateVal] = useState(todayStr())
  const [timeVal, setTimeVal] = useState('08:00')
  const [note, setNote] = useState('')
  const [saving, setSaving] = useState(false)

  // โหลดผู้ป่วยทั้งหมดครั้งแรก
  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoading(true); setError('')
      try {
        const res = await listPatients() // GET /doctor/hospital-patients
        const list = Array.isArray(res?.data) ? res.data : []
        if (!cancelled) {
          const mapped = list.map(mapPatientDTO)
          setAllRows(mapped)
          setSelectedRows(mapped)
        }
      } catch (e) {
        if (!cancelled) setError(e.message || 'โหลดข้อมูลไม่สำเร็จ')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => { cancelled = true }
  }, [])

  // ค้นหา Patient Code หรือ เลขบัตร — ถ้าเว้นว่างให้แสดงทั้งหมด
  const onSearch = async () => {
    setError('')
    const raw = q.trim()
    if (!raw) { setSelectedRows(allRows); return }

    const idQ = onlyDigits(raw)
    const codeQ = toPatientCodeCandidate(raw)
    const rawUpper = raw.toUpperCase()

    try {
      setLoading(true)
      const res = await listPatients({ q: raw })
      const list = Array.isArray(res?.data) ? res.data : []

      // ลำดับ exact:
      // 1) patient_code เท่ากับที่พิมพ์ (รองรับตัวอักษรในอนาคต)
      // 2) patient_code เท่ากับเลข 6 หลัก (pad ซ้าย)
      // 3) id_card_number เท่ากับเลขบัตร (normalize ตัวเลข)
      const exact =
        list.find(p => String(p.patient_code || '').toUpperCase() === rawUpper) ||
        list.find(p => String(p.patient_code || '') === codeQ) ||
        list.find(p => onlyDigits(String(p.id_card_number || '')) === idQ)

      if (exact) {
        setSelectedRows([mapPatientDTO(exact)])
      } else {
        setSelectedRows([])
        setError('ไม่พบผู้ป่วย')
      }
    } catch (e) {
      setError(e.message || 'ค้นหาไม่สำเร็จ')
    } finally {
      setLoading(false)
    }
  }
  const onKeyDown = (e) => { if (e.key === 'Enter') { e.preventDefault(); onSearch() } }

  // modal controls
  useEffect(() => {
    if (open) { document.body.style.overflow = 'hidden' } else { document.body.style.overflow = '' }
    return () => { document.body.style.overflow = '' }
  }, [open])

  const openModalFor = (p) => {
    setSelected(p)
    setDateVal(todayStr())
    setTimeVal('08:00')
    setNote('')
    setOpen(true)
  }
  const closeModal = useCallback(() => setOpen(false), [])
  useEffect(() => {
    const h = (e) => { if (e.key === 'Escape') closeModal() }
    if (open) window.addEventListener('keydown', h)
    return () => window.removeEventListener('keydown', h)
  }, [open, closeModal])

  // กด “นัดหมาย” -> ยิง POST /doctor/appointments
  const confirmAppointment = async () => {
    if (saving) return
    if (!selected?.raw?.id_card_number) {
      alert('ไม่พบเลขบัตรของผู้ป่วย')
      return
    }
    if (!dateVal) {
      alert('กรุณาเลือกวัน')
      return
    }
    if (!/^\d{2}:\d{2}$/.test(timeVal)) {
      alert('รูปแบบเวลาไม่ถูกต้อง (ต้องเป็น HH:mm)')
      return
    }
    try {
      setSaving(true)
      await createAppointment({
        id_card_number: selected.raw.id_card_number,
        appointment_date: dateVal,           // รูปแบบ YYYY-MM-DD
        appointment_time: timeVal,           // รูปแบบ HH:mm
        note: note.trim() || undefined,
      })
      alert(`นัดหมายเรียบร้อย\nผู้ป่วย: ${selected.name}\nวันที่: ${dateVal}\nเวลา: ${timeVal}\nNote: ${note.trim() || '-'}`)
      setOpen(false)
    } catch (e) {
      alert(e.message || 'สร้างนัดหมายไม่สำเร็จ')
    } finally {
      setSaving(false)
    }
  }

  const rows = useMemo(() => selectedRows, [selectedRows])

  return (
    <div>
      <h2 className={styles.title}>Appointment</h2>

      {/* แถวค้นหา */}
      <div className={styles.searchRow}>
        <div className={styles.searchWrap}>
          <input
            className={styles.searchInput}
            placeholder="ค้นหา Patient Code หรือ เลขบัตรประชาชน"
            value={q}
            onChange={e=>setQ(e.target.value)}
            onKeyDown={onKeyDown}
          />
          <span className={styles.searchIcon}>🔍</span>
        </div>
        <button className={styles.searchBtn} onClick={onSearch} disabled={loading}>
          {loading ? 'กำลังค้นหา...' : 'ค้นหา'}
        </button>
      </div>

      {error && <div className={styles.error}>{error}</div>}

      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th style={{width:'6%'}}>#</th>
              <th style={{width:'26%'}}>Name</th>
              <th style={{width:'16%'}}>Patient Code</th>
              <th style={{width:'22%'}}>IDCardNumber</th>
              <th style={{width:'12%'}}>Gender</th>
              <th style={{width:'10%'}}>Age</th>
              <th># Action</th>
            </tr>
          </thead>
          <tbody>
            {loading ? (
              <tr><td colSpan={7} style={{textAlign:'center', color:'#6b7280', height:56}}>กำลังโหลด...</td></tr>
            ) : rows.length === 0 ? (
              <tr><td colSpan={7} style={{textAlign:'center', color:'#6b7280', height:56}}>ไม่มีข้อมูล</td></tr>
            ) : (
              rows.map((p, i) => (
                <tr key={p.id}>
                  <td>{i + 1}</td>
                  <td>{p.name}</td>
                  <td>{p.code}</td>
                  <td>{p.idcard}</td>
                  <td>{p.gender}</td>
                  <td>{p.age}</td>
                  <td className={styles.actions}>
                    <button className={styles.viewBtn} onClick={()=>nav(`/doc/appointment/view/${p.id}`)}>View</button>
                    <button className={styles.addBtn} onClick={()=>openModalFor(p)}>Add</button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {/* Modal */}
      {open && (
        <div className={styles.modalBackdrop} onClick={closeModal} aria-hidden="true">
          <div className={styles.modal} role="dialog" aria-modal="true" onClick={e=>e.stopPropagation()}>
            <button className={styles.modalClose} onClick={closeModal} aria-label="close">×</button>

            <h3 className={styles.modalTitle}>นัดหมาย</h3>

            <div className={styles.formRow}>
              <div className={styles.formGroup}>
                <label className={styles.label}>วัน</label>
                <input
                  type="date"
                  className={styles.input}
                  value={dateVal}
                  min={todayStr()}
                  onChange={e=>setDateVal(e.target.value)}
                />
              </div>
              <div className={styles.formGroup}>
                <label className={styles.label}>เวลา</label>
                <input
                  type="time"
                  className={styles.input}
                  value={timeVal}
                  onChange={e=>setTimeVal(e.target.value)}
                  lang="th-TH"
                  step="60"
                  inputMode="numeric"
                  pattern="^\\d{2}:\\d{2}$"
                  aria-label="เวลา (24 ชั่วโมง)"
                />
              </div>
            </div>

            <div className={styles.formGroup}>
              <label className={styles.label}>Note</label>
              <textarea
                className={styles.textarea}
                rows={4}
                placeholder="เช่น งดอาหาร 8 ชั่วโมง ก่อนเจาะเลือด"
                value={note}
                onChange={e=>setNote(e.target.value)}
              />
            </div>

            <button className={styles.primaryBtn} onClick={confirmAppointment} disabled={saving}>
              {saving ? 'กำลังสร้างนัดหมาย…' : 'นัดหมาย'}
            </button>
          </div>
        </div>
      )}
    </div>
  )
}
