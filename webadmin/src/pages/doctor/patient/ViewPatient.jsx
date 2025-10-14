import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/patient/ViewPatient.module.css'
import { getPatient } from '../../../services/patients'

const GENDERS = ['ชาย', 'หญิง']

function calcAgeYMD(ymd) {
  if (!ymd) return ''
  const dob = new Date(`${ymd}T00:00:00`)
  if (isNaN(dob)) return ''
  const t = new Date()
  let a = t.getFullYear() - dob.getFullYear()
  const m = t.getMonth() - dob.getMonth()
  if (m < 0 || (m === 0 && t.getDate() < dob.getDate())) a--
  return String(Math.max(a, 0))
}

export default function ViewPatient() {
  const { id } = useParams()
  const nav = useNavigate()

  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // เก็บ DTO ต้นฉบับ (จาก BE)
  const [rec, setRec] = useState(null)

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoading(true); setError('')
      try {
        const res = await getPatient(id) // -> { data: {...} }
        if (cancelled) return
        if (!res?.data) throw new Error('ไม่พบข้อมูลผู้ป่วย')
        setRec(res.data)
      } catch (e) {
        if (!cancelled) setError(e.message || 'โหลดข้อมูลไม่สำเร็จ')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => { cancelled = true }
  }, [id])

  // map DTO -> UI fields
  const view = useMemo(() => {
    if (!rec) return null
    return {
      firstName: rec.first_name || '',
      lastName:  rec.last_name  || '',
      birthDay:  rec.birth_day  || '',           // "YYYY-MM-DD" จาก DTO
      gender:    rec.gender     || '',
      idcard:    rec.id_card_number || '',
      phone:     rec.phone_number  || '',
    }
  }, [rec])

  const age = useMemo(() => calcAgeYMD(view?.birthDay || ''), [view])

  if (loading) return <div className={styles.page}>กำลังโหลด...</div>
  if (error)   return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h2 className={styles.title}>View Patient</h2>
        <button className={styles.back} onClick={()=>nav(-1)}>← Back</button>
      </div>
      <hr className={styles.hr} />
      <div className={styles.card}>
        <div className={styles.error}>{error}</div>
      </div>
    </div>
  )
  if (!view)   return <div className={styles.page}>ไม่พบข้อมูลผู้ป่วย</div>

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h2 className={styles.title}>View Patient</h2>
        <button className={styles.back} onClick={()=>nav(-1)}>← Back</button>
      </div>
      <hr className={styles.hr} />

      <div className={styles.card}>
        <div className={styles.row}>
          <label className={styles.label}>
            <span>First Name</span>
            <input className={styles.input} value={view.firstName} disabled />
          </label>
          <label className={styles.label}>
            <span>Birth Day</span>
            <input className={styles.input} type="date" value={view.birthDay} disabled />
          </label>
        </div>

        <div className={styles.row}>
          <label className={styles.label}>
            <span>Last Name</span>
            <input className={styles.input} value={view.lastName} disabled />
          </label>
          <label className={styles.label}>
            <span>Age</span>
            <div className={styles.inputGroup}>
              <input className={styles.input} value={age} disabled />
              <span className={styles.suffix}>ปี</span>
            </div>
          </label>
        </div>

        <div className={styles.row}>
          <label className={styles.label}>
            <span>ID Card Number</span>
            <input className={styles.input} value={view.idcard} disabled />
          </label>
          <label className={styles.label}>
            <span>Gender</span>
            <select className={styles.select} value={view.gender} disabled>
              {[view.gender || '-', ...GENDERS.filter(g => g !== view.gender)].map(g => (
                <option key={g}>{g}</option>
              ))}
            </select>
          </label>
        </div>

        <div className={styles.row}>
          <label className={styles.label} style={{gridColumn:'1 / -1'}}>
            <span>Phone Number</span>
            <input className={styles.input} value={view.phone} disabled />
          </label>
        </div>
      </div>
    </div>
  )
}
