// src/pages/doctor/appointment/ViewAppointment.jsx
import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/appointment/ViewAppointment.module.css'
import { getPatient } from '../../../services/patients'
import { listAppointments } from '../../../services/appointments'

// ===== helpers แสดงวันที่/เวลาแบบไทย (ล็อกโซน Bangkok เสมอ) =====
const BANGKOK_TZ = 'Asia/Bangkok'

function formatDateTH(d) {
  const dt = new Date(d)
  return dt.toLocaleDateString('th-TH', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
    timeZone: BANGKOK_TZ,
  })
}

function formatTimeTH(d) {
  const dt = new Date(d)
  return dt.toLocaleTimeString('th-TH', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false,         // บังคับ 24 ชั่วโมง
    timeZone: BANGKOK_TZ,  // ล็อกเป็นเวลาไทย
  })
}

// รวม date + time เป็น ISO ของกรุงเทพฯ ถ้า input ไม่มี timezone
function toBangkokISO(dateYMD, timeHM = '00:00') {
  if (!dateYMD) return ''
  const hasTZ = /[zZ]|[+-]\d{2}:\d{2}$/.test(timeHM) || /[zZ]|[+-]\d{2}:\d{2}$/.test(dateYMD)
  // รองรับ "HH:mm:ss" หรือ "HH:mm"
  const hhmmss = timeHM.length === 5 ? `${timeHM}:00` : timeHM
  if (hasTZ) {
    // ผู้ให้บริการส่ง timezone มาแล้ว ไม่แตะต้อง
    return `${dateYMD}T${hhmmss}`
  }
  return `${dateYMD}T${hhmmss}+07:00`
}

export default function ViewAppointment() {
  const nav = useNavigate()
  const { id: patientId } = useParams() // /doc/appointment/view/:id

  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [patient, setPatient] = useState(null)
  const [rows, setRows] = useState([])

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoading(true); setError('')

      try {
        // 1) โหลดผู้ป่วย เพื่อเอา id_card_number + ชื่อ
        const pRes = await getPatient(patientId) // -> { data: {...} }
        const p = pRes?.data
        if (!p) throw new Error('ไม่พบข้อมูลผู้ป่วย')
        if (cancelled) return
        setPatient(p)

        // 2) โหลดนัดหมายด้วย q = id_card_number
        const aRes = await listAppointments({ q: p.id_card_number }) // -> { data: [...] }
        const arr = Array.isArray(aRes?.data) ? aRes.data : []

        // 3) map -> rows [{ id, at, note }]
        const mapped = arr.map(ap => {
          // รองรับหลายชื่อฟิลด์จาก BE
          const dateYMD = ap.appointment_date || ap.date || ap.appointed_date || ''
          const rawTime  = ap.appointment_time || ap.time || ap.appointed_time || '00:00'
          const at = ap.at || ap.appointed_at || (dateYMD ? toBangkokISO(dateYMD, rawTime) : '')

          return {
            id: ap.id,
            at,
            note: ap.note || '',
          }
        })

        // เรียงใหม่สุดก่อน (ตามค่ามิลลิวินาที)
        mapped.sort((a, b) => new Date(b.at).getTime() - new Date(a.at).getTime())

        if (!cancelled) setRows(mapped)
      } catch (e) {
        if (!cancelled) setError(e.message || 'โหลดข้อมูลไม่สำเร็จ')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()

    return () => { cancelled = true }
  }, [patientId])

  const patientName = useMemo(() => {
    if (!patient) return `Patient #${patientId}`
    const n = [patient.first_name, patient.last_name].filter(Boolean).join(' ')
    return n || `Patient #${patientId}`
  }, [patient, patientId])

  return (
    <div>
      <div className={styles.header}>
        <h2 className={styles.title}>Appointment</h2>
        <button className={styles.back} onClick={()=>nav(-1)}>← Back</button>
      </div>

      {patient && (
        <div className={styles.meta}>
          <div className={styles.badge}>ผู้ป่วย</div>
          <div className={styles.name}>{patientName}</div>
        </div>
      )}

      {loading ? (
        <div className={styles.loading}>กำลังโหลด...</div>
      ) : error ? (
        <div className={styles.error}>{error}</div>
      ) : (
        <div className={styles.tableWrap}>
          <table className={styles.table}>
            <thead>
              <tr>
                <th style={{width:'8%'}}>#</th>
                <th style={{width:'22%'}}>วันที่</th>
                <th style={{width:'20%'}}>เวลา</th>
                <th>หมายเหตุ</th>
                <th style={{width:'18%'}}></th>
              </tr>
            </thead>
            <tbody>
              {rows.length === 0 ? (
                <tr>
                  <td colSpan={5} style={{textAlign:'center', color:'#6b7280', height:56}}>
                    ยังไม่มีการนัดหมาย
                  </td>
                </tr>
              ) : (
                rows.map((r, i) => {
                  const isLatest = i === 0
                  return (
                    <tr key={r.id} className={isLatest ? styles.latestRow : undefined}>
                      <td>{i+1}</td>
                      <td>{r.at ? formatDateTH(r.at) : '-'}</td>
                      <td>{r.at ? formatTimeTH(r.at) : '-'}</td>
                      <td className={styles.note}>{r.note || '-'}</td>
                      <td className={styles.statusCell}>
                        <div className={styles.rowActions}>
                          {isLatest && <span className={styles.latestBadge}>ล่าสุด</span>}
                        </div>
                      </td>
                    </tr>
                  )
                })
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
