import { useState, useEffect, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/doctor/appointment/AppointmentList.module.css'

const MOCK_PATIENTS = [
  { id:1, name:'สมชาย ใจดี', idcard:'1234567890100', gender:'ชาย', age:30 },
  { id:2, name:'สมหญิง ใจร้าย', idcard:'1234567890101', gender:'หญิง', age:40 },
  { id:3, name:'สมหมาย ใจบุญ', idcard:'1234567890102', gender:'ชาย', age:55 },
]

function todayStr() {
  const d = new Date()
  const yyyy = d.getFullYear()
  const mm = String(d.getMonth() + 1).padStart(2, '0')
  const dd = String(d.getDate()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd}`
}

export default function AppointmentList() {
  const nav = useNavigate()

  const [q, setQ] = useState('')
  const [results, setResults] = useState(MOCK_PATIENTS)   // เริ่มต้นแสดงทุกคน
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  // ------ modal state ------
  const [open, setOpen] = useState(false)
  const [selected, setSelected] = useState(null) // patient ที่จะนัด
  const [dateVal, setDateVal] = useState(todayStr())
  const [timeVal, setTimeVal] = useState('08:00')
  const [note, setNote] = useState('')

  // ป้องกัน scroll พื้นหลังเมื่อเปิด modal
  useEffect(() => {
    if (open) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }
    return () => { document.body.style.overflow = '' }
  }, [open])

  const onSearch = async () => {
    setError('')
    const qtrim = q.trim()
    try {
      setLoading(true)
      await new Promise(r => setTimeout(r, 150)) // mock latency
      if (!qtrim) { setResults(MOCK_PATIENTS); return }
      const found = MOCK_PATIENTS.filter(p => p.idcard === qtrim)
      if (found.length === 0) {
        setResults([])
        setError('ไม่พบผู้ป่วยตามเลขบัตรนี้')
      } else {
        setResults(found)
      }
    } finally {
      setLoading(false)
    }
  }

  const onKeyDown = (e) => { if (e.key === 'Enter') { e.preventDefault(); onSearch() } }

  const openModalFor = (p) => {
    setSelected(p)
    setDateVal(todayStr())
    setTimeVal('08:00')
    setNote('')
    setOpen(true)
  }

  const closeModal = useCallback(() => setOpen(false), [])

  // ปิดด้วย ESC
  useEffect(() => {
    const h = (e) => { if (e.key === 'Escape') closeModal() }
    if (open) window.addEventListener('keydown', h)
    return () => window.removeEventListener('keydown', h)
  }, [open, closeModal])

  const confirmAppointment = () => {
    // TODO: เรียก API POST /api/appointments ที่นี่
    console.log('create appointment', {
      patientId: selected?.id,
      date: dateVal,
      time: timeVal,
      note,
    })
    alert(`(mock) นัดหมายเรียบร้อย\nผู้ป่วย: ${selected?.name}\nวันที่: ${dateVal}\nเวลา: ${timeVal}\nNote: ${note || '-'}`)
    setOpen(false)
  }

  return (
    <div>
      <h2 className={styles.title}>Appointment</h2>

      {/* แถวค้นหา */}
      <div className={styles.searchRow}>
        <div className={styles.searchWrap}>
          <input
            className={styles.searchInput}
            placeholder="ค้นหาเลขบัตรประชาชน"
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
              <th style={{width:'32%'}}>Name</th>
              <th style={{width:'26%'}}>IDCardNumber</th>
              <th style={{width:'12%'}}>Gender</th>
              <th style={{width:'10%'}}>Age</th>
              <th># Action</th>
            </tr>
          </thead>
          <tbody>
            {results.length === 0 ? (
              <tr>
                <td colSpan={6} style={{textAlign:'center', color:'#6b7280', height:56}}>ไม่มีข้อมูล</td>
              </tr>
            ) : (
              results.map((p, i) => (
                <tr key={p.id}>
                  <td>{i + 1}</td>
                  <td>{p.name}</td>
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

      {/* ---------- Modal ---------- */}
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

            <button className={styles.primaryBtn} onClick={confirmAppointment}>นัดหมาย</button>
          </div>
        </div>
      )}
    </div>
  )
}
