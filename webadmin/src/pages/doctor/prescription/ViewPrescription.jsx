// src/pages/doctor/prescription/ViewPrescription.jsx
import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/prescription/ViewPrescription.module.css'
import { getPatient } from '../../../services/patients'
import { listPrescriptions } from '../../../services/prescriptions'
import { getMedicineForDoctor } from '../../../services/medicines'
import { getForm, getUnit } from '../../../services/initialData'   // <- เพิ่ม getUnit
import { getDoctorPublic } from '../../../services/doctors'

// helper แสดงวันที่/เวลาแบบไทย
function formatDateTH(d) {
  const dt = new Date(d)
  return dt.toLocaleDateString('th-TH', { day: 'numeric', month: 'short', year: 'numeric' })
}
function formatTimeTH(d) {
  const dt = new Date(d)
  return dt.toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' })
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
        // 1) โหลดผู้ป่วย
        const pRes = await getPatient(patientId)
        const p = pRes?.data
        if (!p) throw new Error('ไม่พบข้อมูลผู้ป่วย')
        if (cancelled) return
        setPatient(p)

        // 2) โหลดใบสั่งยาทั้งหมดของผู้ป่วย ด้วย q = id_card_number
        const rxRes = await listPrescriptions({ q: p.id_card_number })
        const list = Array.isArray(rxRes?.data) ? rxRes.data : []
        if (cancelled) return

        // ===== เตรียม lookup ยา (medicine_info) =====
        const uniqueMedIDs = [...new Set(
          list.flatMap(rx => (Array.isArray(rx.items) ? rx.items : [])).map(it => it.medicine_info_id)
        )].filter(Boolean)

        const medMap = new Map() // id -> raw medicine_info
        await Promise.all(uniqueMedIDs.map(async (mid) => {
          try {
            const mres = await getMedicineForDoctor(mid)   // GET /doctor/medicine-info/:id
            const m = mres?.data
            if (m) medMap.set(mid, m)
          } catch {/* noop */}
        }))

        // ===== เตรียม lookup หมอ =====
        const uniqueDoctorIDs = [...new Set(
          list.map(rx => rx.doctor_id ?? rx.doctorId).filter(Boolean)
        )]
        const doctorMap = new Map() // id -> "ชื่อ นามสกุล"
        await Promise.all(uniqueDoctorIDs.map(async (did) => {
          try {
            const dres = await getDoctorPublic(did) // GET /doctor/doctors/:id
            const d = dres?.data
            if (d) {
              const name = [d.first_name, d.last_name].filter(Boolean).join(' ') || (d.username || '-')
              doctorMap.set(did, name)
            }
          } catch {/* noop */}
        }))

        // ===== เก็บ form_id / unit_id ที่ต้องใช้ (จาก item และ medicine_info) =====
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

        // ===== ดึงชื่อ form/unit ตาม id =====
        const formMap = new Map() // id -> form_name
        const unitMap = new Map() // id -> unit_name

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

        // 3) map ใบสั่งยา → เติมชื่อยา/รูปแบบ/หน่วย/ความแรง + ชื่อหมอ + ประกอบข้อความ
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

              // ---- ชื่อยา/generic/strength ----
              const medName =
                it.med_name || it.MedName ||
                mi?.med_name || mi?.MedName || '-'
              const generic =
                it.generic_name || it.GenericName ||
                mi?.generic_name || mi?.GenericName || '-'
              const strength =
                it.strength || mi?.strength || '-'

              // ---- form ----
              const formName =
                it.form_name ||
                (it.form_id ? formMap.get(it.form_id) : '') ||
                mi?.form_name ||
                (mi?.form_id ? formMap.get(mi.form_id) : '') ||
                '-'

              // ---- unit ----
              const unitName =
                it.unit_name ||
                (it.unit_id ? unitMap.get(it.unit_id) : '') ||
                mi?.unit_name ||
                (mi?.unit_id ? unitMap.get(mi.unit_id) : '') ||
                '' // ไม่มีจริง ๆ ปล่อยว่าง

              // ---- dosage text ----
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

              return {
                medName,
                generic,
                strength,
                form: formName,
                dosePerTime,   // “ครั้งละ X หน่วย/ฟอร์ม”
                timesPerDay,   // “วันละ N ครั้ง”
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

  // หาว่าอันไหนล่าสุด (เอาตามเวลา orderedAt มากสุด)
  const latestAt = useMemo(() => {
    if (!data.length) return null
    return Math.max(...data.map(rx => new Date(rx.orderedAt).getTime()))
  }, [data])

  // จัดกลุ่มตามวันที่ (ใหม่ → เก่า) และเรียงเวลาในวัน (ใหม่ → เก่า)
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
                          <th style={{width:'26%'}}>MedName</th>
                          <th style={{width:'26%'}}>GenericName</th>
                          <th style={{width:'14%'}}>Strength</th>
                          <th style={{width:'12%'}}>Form</th>
                          <th style={{width:'30%'}}>Dosage</th>
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
                              ครั้งละ <strong>{it.dosePerTime}</strong> | วันละ <strong>{it.timesPerDay}</strong>
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
