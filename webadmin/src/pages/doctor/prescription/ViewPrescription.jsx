import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/prescription/ViewPrescription.module.css'
import { getPatient } from '../../../services/patients'
import { listPrescriptions } from '../../../services/prescriptions'
import { getMedicineForDoctor } from '../../../services/medicines'
import { getForm, getUnit } from '../../../services/initialData'
import { getDoctorPublic } from '../../../services/doctors'

// ==== helpers ====
function formatDateTH(d) {
  const dt = new Date(d)
  return dt.toLocaleDateString('th-TH', { day: 'numeric', month: 'short', year: 'numeric' })
}
function formatTimeTH(d) {
  const dt = new Date(d)
  return dt.toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' })
}
/** YYYY-MM-DD -> DD/MM/YYYY (สำหรับแสดงผล) */
function isoToDDMMYYYY(iso) {
  const s = String(iso || '')
  if (!/^\d{4}-\d{2}-\d{2}$/.test(s)) return ''
  const [y, m, d] = s.split('-')
  return `${d}/${m}/${y}`
}
/** YYYY-MM-DD -> YYYY-MM-DD (end+1 วัน) */
function addOneDayISO(iso) {
  const s = String(iso || '')
  if (!/^\d{4}-\d{2}-\d{2}$/.test(s)) return ''
  const dt = new Date(`${s}T00:00:00`)
  if (Number.isNaN(dt.getTime())) return ''
  dt.setDate(dt.getDate() + 1)
  const y = dt.getFullYear()
  const m = String(dt.getMonth() + 1).padStart(2, '0')
  const d = String(dt.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

export default function ViewPrescription() {
  const nav = useNavigate()
  const { id: patientId } = useParams()

  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [patient, setPatient] = useState(null)
  const [data, setData] = useState([])

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoading(true); setError('')

      try {
        // 1) ผู้ป่วย
        const pRes = await getPatient(patientId)
        const p = pRes?.data
        if (!p) throw new Error('ไม่พบข้อมูลผู้ป่วย')
        if (cancelled) return
        setPatient(p)

        // 2) ใบสั่งยาทั้งหมด (q = id_card_number)
        const rxRes = await listPrescriptions({ q: p.id_card_number })
        const list = Array.isArray(rxRes?.data) ? rxRes.data : []
        if (cancelled) return

        // ===== เตรียม lookup ยา (medicine_info) =====
        const uniqueMedIDs = [...new Set(
          list.flatMap(rx => (Array.isArray(rx.items) ? rx.items : [])).map(it => it.medicine_info_id)
        )].filter(Boolean)

        const medMap = new Map()
        await Promise.all(uniqueMedIDs.map(async (mid) => {
          try {
            const mres = await getMedicineForDoctor(mid)
            const m = mres?.data
            if (m) medMap.set(mid, m)
          } catch {/* noop */}
        }))

        // ===== เตรียม lookup หมอ =====
        const uniqueDoctorIDs = [...new Set(
          list.map(rx => rx.doctor_id ?? rx.doctorId).filter(Boolean)
        )]
        const doctorMap = new Map()
        await Promise.all(uniqueDoctorIDs.map(async (did) => {
          try {
            const dres = await getDoctorPublic(did)
            const d = dres?.data
            if (d) {
              const name = [d.first_name, d.last_name].filter(Boolean).join(' ') || (d.username || '-')
              doctorMap.set(did, name)
            }
          } catch {/* noop */}
        }))

        // ===== เก็บ form_id / unit_id ที่ต้องใช้ =====
        const formIDs = new Set()
        const unitIDs = new Set()

        for (const rx of list) {
          for (const it of (rx.items || [])) {
            if (it.form_id) formIDs.add(it.form_id)
            if (it.unit_id) unitIDs.add(it.unit_id)

            const mi = it.medicine_info_id ? medMap.get(it.medicine_info_id) : null
            if (mi?.form_id) formIDs.add(mi.form_id)
            if (mi?.unit_id) unitIDs.add(mi.unit_id)
          }
        }

        // ===== ดึงชื่อ form/unit =====
        const formMap = new Map()
        const unitMap = new Map()

        await Promise.all([
          ...[...formIDs].map(async (fid) => {
            try {
              const fres = await getForm(fid)
              const fObj = fres?.data || fres
              const name = fObj?.form_name || fObj?.name || ''
              if (name) formMap.set(fid, name)
            } catch {/* noop */}
          }),
          ...[...unitIDs].map(async (uid) => {
            try {
              const ures = await getUnit(uid)
              const uObj = ures?.data || ures
              const name = uObj?.unit_name || uObj?.name || ''
              if (name) unitMap.set(uid, name)
            } catch {/* noop */}
          }),
        ])

        // 3) map ใบสั่งยา + รายการ
        const mapped = list.map(rx => {
          const did = rx.doctor_id ?? rx.doctorId
          const doctorName = (did && doctorMap.get(did)) || rx.doctor_name || rx.doctorName || '-'
          const orderedAt =
            rx.ordered_at || rx.orderedAt || rx.created_at || rx.createdAt || new Date().toISOString()

          return {
            id: rx.id,
            orderedAt,
            doctorName,
            items: (Array.isArray(rx.items) ? rx.items : []).map(it => {
              const mi = it.medicine_info_id ? medMap.get(it.medicine_info_id) : null

              const medName =
                it.med_name || it.MedName ||
                mi?.med_name || mi?.MedName || '-'
              const generic =
                it.generic_name || it.GenericName ||
                mi?.generic_name || mi?.GenericName || '-'
              const strength =
                it.strength || mi?.strength || '-'

              const formName =
                it.form_name ||
                (it.form_id ? formMap.get(it.form_id) : '') ||
                mi?.form_name ||
                (mi?.form_id ? formMap.get(mi.form_id) : '') ||
                '-'

              const unitName =
                it.unit_name ||
                (it.unit_id ? unitMap.get(it.unit_id) : '') ||
                mi?.unit_name ||
                (mi?.unit_id ? unitMap.get(mi.unit_id) : '') ||
                ''

              const rawDosePerTime = it.amount_per_time ?? it.dosePerTime
              const dosePerTime =
                (rawDosePerTime !== undefined && rawDosePerTime !== null && String(rawDosePerTime) !== '')
                  ? `${rawDosePerTime}${unitName ? ` ${unitName}` : (formName && formName !== '-' ? ` ${formName}` : '')}`
                  : '-'

              const rawTimesPerDay = it.times_per_day ?? it.timesPerDay
              const timesPerDay =
                (rawTimesPerDay !== undefined && rawTimesPerDay !== null && String(rawTimesPerDay) !== '')
                  ? `${rawTimesPerDay} ครั้ง`
                  : '-'

              // ====== ฟิลด์ใหม่: start_date / end_date / expire_date(คำนวณ) / note ======
              const startISO = it.start_date ?? it.startDate ?? null
              const endISO   = it.end_date   ?? it.endDate   ?? null
              const note     = it.note ?? it.Note ?? ''

              const expireISO = endISO ? addOneDayISO(endISO) : null

              return {
                medName, generic, strength, form: formName,
                dosePerTime, timesPerDay,
                startISO, endISO, expireISO, note
              }
            })
          }
        })

        if (!cancelled) setData(mapped)
      } catch (e) {
        if (!cancelled) setError(e.message || 'โหลดข้อมูลไม่สำเร็จ')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()

    return () => { cancelled = true }
  }, [patientId])

  const latestAt = useMemo(() => {
    if (!data.length) return null
    return Math.max(...data.map(rx => new Date(rx.orderedAt).getTime()))
  }, [data])

  const groups = useMemo(() => {
    const map = new Map()
    for (const rx of data) {
      const key = formatDateTH(rx.orderedAt)
      if (!map.has(key)) map.set(key, [])
      map.get(key).push(rx)
    }
    const sorted = [...map.entries()].sort(([a],[b]) => {
      const ad = new Date(a); const bd = new Date(b)
      return bd - ad
    })
    for (const [, arr] of sorted) {
      arr.sort((a,b) => new Date(b.orderedAt) - new Date(a.orderedAt))
    }
    return sorted
  }, [data])

  return (
    <div>
      <div className={styles.header}>
        <h2 className={styles.title}>View Prescription</h2>
        <button className={styles.back} onClick={() => nav(-1)}>← Back</button>
      </div>

      {patient && (
        <div className={styles.patientBar}>
          <div><strong>ชื่อผู้ป่วย:</strong> {[patient.first_name, patient.last_name].filter(Boolean).join(' ') || '-'}</div>
          <div><strong>Patient Code:</strong> {patient.patient_code || '-'}</div>
          <div><strong>เลขบัตร:</strong> {patient.id_card_number || '-'}</div>
        </div>
      )}

      {loading ? (
        <div className={styles.empty}>กำลังโหลด...</div>
      ) : error ? (
        <div className={styles.error}>{error}</div>
      ) : groups.length === 0 ? (
        <div className={styles.empty}>ยังไม่มีประวัติการสั่งยา</div>
      ) : (
        groups.map(([dateLabel, rxs]) => (
          <section key={dateLabel} className={styles.section}>
            <div className={styles.dateBadge}>{dateLabel}</div>

            {rxs.map((rx) => {
              const isLatest = latestAt && (new Date(rx.orderedAt).getTime() === latestAt)
              return (
                <div
                  key={rx.id}
                  className={`${styles.card} ${isLatest ? styles.cardLatest : ''}`}
                >
                  <div className={styles.cardHeader}>
                    <div className={styles.rxMeta}>
                      <span className={styles.time}>{formatTimeTH(rx.orderedAt)}</span>
                      <span className={styles.dot}>•</span>
                      <span className={styles.doctor}>แพทย์ผู้สั่ง: {rx.doctorName || '-'}</span>
                      {isLatest && <span className={styles.latestTag}>ล่าสุด</span>}
                    </div>
                  </div>

                  <div className={styles.tableWrap}>
                    <table className={styles.table}>
                      <thead>
                        <tr>
                          <th style={{width:'6%'}}>#</th>
                          <th style={{width:'23%'}}>MedName</th>
                          <th style={{width:'23%'}}>GenericName</th>
                          <th style={{width:'12%'}}>Strength</th>
                          <th style={{width:'12%'}}>Form</th>
                          <th style={{width:'24%'}}>Dosage & Period</th>
                        </tr>
                      </thead>
                      <tbody>
                        {rx.items.map((it, idx) => (
                          <tr key={idx}>
                            <td>
                              <input type="checkbox" checked readOnly className={styles.checkbox}/>
                            </td>
                            <td>{it.medName}</td>
                            <td>{it.generic}</td>
                            <td>{it.strength}</td>
                            <td>{it.form}</td>
                            <td className={styles.dose}>
                              <div>ครั้งละ <strong>{it.dosePerTime}</strong> | วันละ <strong>{it.timesPerDay}</strong></div>

                              {(it.startISO || it.endISO) && (
                                <div className={styles.period}>
                                  {it.startISO && <span>เริ่ม {isoToDDMMYYYY(it.startISO)}</span>}
                                  {it.endISO && <span className={styles.sep}>สิ้นสุด {isoToDDMMYYYY(it.endISO)}</span>}
                                  {it.expireISO && <span className={styles.expire}>หมดอายุ {isoToDDMMYYYY(it.expireISO)}</span>}
                                </div>
                              )}

                              {it.note && (
                                <div className={styles.noteLine}>
                                  <span className={styles.noteLabel}>note:</span> {it.note}
                                </div>
                              )}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              )
            })}
          </section>
        ))
      )}
    </div>
  )
}
