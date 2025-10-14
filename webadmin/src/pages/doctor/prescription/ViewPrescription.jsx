// src/pages/doctor/prescription/ViewPrescription.jsx
import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/prescription/ViewPrescription.module.css'
import { getPatient } from '../../../services/patients'
import { listPrescriptions } from '../../../services/prescriptions'
import { getMedicineForDoctor } from '../../../services/medicines'
import { getForm } from '../../../services/initialData'
import { getDoctorPublic } from '../../../services/doctors'   // <<< เพิ่ม: ดึงชื่อหมอ

// helper แสดงวันที่/เวลาแบบไทย
function formatDateTH(d) {
  const dt = new Date(d)
  return dt.toLocaleDateString('th-TH', { day:'numeric', month:'short', year:'numeric' })
}
function formatTimeTH(d) {
  const dt = new Date(d)
  return dt.toLocaleTimeString('th-TH', { hour:'2-digit', minute:'2-digit' })
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
        const pRes = await getPatient(patientId) // -> { data: { id_card_number, ... } }
        const p = pRes?.data
        if (!p) throw new Error('ไม่พบข้อมูลผู้ป่วย')
        if (cancelled) return
        setPatient(p)

        // 2) โหลดใบสั่งยาทั้งหมดของผู้ป่วย ด้วย q = id_card_number
        const rxRes = await listPrescriptions({ q: p.id_card_number }) // -> { data: [ { id, doctor_id, items:[{medicine_info_id,...}] } ] }
        const list = Array.isArray(rxRes?.data) ? rxRes.data : []
        if (cancelled) return

        // ==== เตรียม lookup ยา ====
        const uniqueMedIDs = [...new Set(
          list.flatMap(rx => (Array.isArray(rx.items) ? rx.items : [])).map(it => it.medicine_info_id)
        )].filter(Boolean)

        const medMap = new Map() // id -> { med_name, generic_name, strength, form_name, ... }
        await Promise.all(uniqueMedIDs.map(async (mid) => {
          try {
            const mres = await getMedicineForDoctor(mid)   // GET /doctor/medicine-info/:id
            const m = mres?.data
            if (m) {
              // ถ้ามี form_name อยู่แล้วก็ใช้เลย ถ้าไม่มีแต่มี form_id → ไปดึงชื่อฟอร์มเพิ่ม
              let formName = m.form_name || m.form || ''
              if (!formName && m.form_id) {
                try {
                  const fres = await getForm(m.form_id) // GET /form/:id
                  // บาง impl อาจคืนเป็น { form_name: 'ยาเม็ด', ... } หรือ { data:{...} }
                  const fObj = fres?.data || fres
                  formName = fObj?.form_name || fObj?.name || ''
                } catch {/* noop */}
              }

              medMap.set(mid, {
                med_name: m.med_name || m.MedName || '-',
                generic_name: m.generic_name || m.GenericName || '-',
                strength: m.strength || '-',
                form_name: formName || '-',
                unit_name: m.unit_name || m.unit || '',
                instruction_name: m.instruction_name || m.instruction || ''
              })
            }
          } catch {
            // ถ้าดึงไม่ได้ ปล่อย '-' ให้ fallback ตอน map
          }
        }))

        // ==== เตรียม lookup หมอ ====
        const uniqueDoctorIDs = [...new Set(
          list.map(rx => rx.doctor_id || rx.doctorId).filter(Boolean)
        )]
        const doctorMap = new Map() // id -> "ชื่อ นามสกุล"
        await Promise.all(uniqueDoctorIDs.map(async (did) => {
          try {
            const dres = await getDoctorPublic(did) // GET /doctor/doctors/:id  -> { data:{ first_name, last_name, ... } }
            const d = dres?.data
            if (d) {
              const name = [d.first_name, d.last_name].filter(Boolean).join(' ') || (d.username || '-')
              doctorMap.set(did, name)
            }
          } catch {
            // ถ้าดึงไม่ได้ ไม่เป็นไร จะโชว์ '-'
          }
        }))

        // 3) map ใบสั่งยา → เติมชื่อยา/รูปแบบ/ความแรง และชื่อหมอ ลงใน items / meta
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
              const info = medMap.get(it.medicine_info_id) || {}
              return {
                medName: info.med_name || '-',
                generic: info.generic_name || '-',
                strength: info.strength || '-',
                form: info.form_name || '-',
                dosePerTime: it.amount_per_time ?? it.dosePerTime ?? '-',
                timesPerDay: it.times_per_day ?? it.timesPerDay ?? '-',
              }
            })
          }
        })

        setData(mapped)
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
