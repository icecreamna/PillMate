import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/doctor/patient/AddPatient.module.css'
import { createPatient } from '../../../services/patients'

const GENDERS = ['ชาย', 'หญิง']

const onlyDigits = (s) => (s || '').replace(/[^\d]/g, '')

// กัน timezone ด้วยการเติม T00:00:00 แล้วส่งเป็น RFC3339 (+07:00)
function toBangkokISO(d) {
  if (!d) return ''
  return `${d}T00:00:00+07:00`
}

function calcAge(dobStrYMD) {
  if (!dobStrYMD) return ''
  const dob = new Date(`${dobStrYMD}T00:00:00`)
  if (isNaN(dob)) return ''
  const today = new Date()
  let age = today.getFullYear() - dob.getFullYear()
  const m = today.getMonth() - dob.getMonth()
  if (m < 0 || (m === 0 && today.getDate() < dob.getDate())) age--
  return String(Math.max(age, 0))
}

export default function AddPatient() {
  const nav = useNavigate()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [ok, setOk] = useState('')

  const [form, setForm] = useState({
    firstName: '',
    lastName: '',
    birthDay: '',   // "YYYY-MM-DD"
    age: '',
    gender: '',
    idcard: '',
    phone: '',
  })

  const onChange = (k, v) => setForm(s => ({ ...s, [k]: v }))

  const validate = () => {
    if (!form.firstName.trim()) return 'กรุณากรอก First Name'
    if (!form.lastName.trim())  return 'กรุณากรอก Last Name'
    if (!form.birthDay)         return 'กรุณาเลือกวันเกิด'
    if (!form.gender)           return 'กรุณาเลือกเพศ'

    // บังคับเลขบัตร 13 หลัก
    const id13 = onlyDigits(form.idcard)
    if (!id13) return 'กรุณากรอกเลขบัตรประชาชน'
    if (id13.length !== 13) return 'เลขบัตรประชาชนต้องมี 13 หลัก'

    // บังคับเบอร์โทร 10 หลัก
    const tel = onlyDigits(form.phone)
    if (!tel) return 'กรุณากรอกเบอร์โทรศัพท์'
    if (tel.length !== 10) return 'เบอร์โทรต้องมี 10 หลัก'

    // กันเพศผิดค่า
    if (!GENDERS.includes(form.gender)) return 'ค่าเพศไม่ถูกต้อง'

    return ''
  }

  const onSubmit = async (e) => {
    e.preventDefault()
    const v = validate()
    if (v) { setError(v); setOk(''); return }
    setError(''); setOk(''); setLoading(true)
    try {
      const payload = {
        id_card_number: onlyDigits(form.idcard),      // 13 หลัก
        first_name: form.firstName.trim(),
        last_name: form.lastName.trim(),
        phone_number: onlyDigits(form.phone),         // 10 หลัก
        birth_day: toBangkokISO(form.birthDay),       // RFC3339 +07:00
        gender: form.gender,                           // "ชาย" | "หญิง"
      }
       console.log('POST payload:', payload)
      await createPatient(payload) // POST /doctor/hospital-patients
      setOk('เพิ่มผู้ป่วยเรียบร้อย')
      nav('/doc/patients', { replace:true })
    } catch (err) {
      setError(err.message || 'บันทึกไม่สำเร็จ')
    } finally { setLoading(false) }
  }

  const todayStr = new Date().toISOString().slice(0, 10)

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h2 className={styles.title}>Add Patient</h2>
        <button className={styles.back} onClick={() => nav(-1)}>← Back</button>
      </div>
      <hr className={styles.hr} />

      <div className={styles.card}>
        {error && <div className={styles.error}>{error}</div>}
        {ok && <div className={styles.success}>{ok}</div>}

        <form className={styles.form} onSubmit={onSubmit}>
          <div className={styles.row}>
            <label className={styles.label}>
              <span>First Name</span>
              <input
                className={styles.input}
                value={form.firstName}
                onChange={e => onChange('firstName', e.target.value)}
                autoComplete="given-name"
              />
            </label>

            <label className={styles.label}>
              <span>Birth Day</span>
              <input
                className={styles.input}
                type="date"
                max={todayStr}
                value={form.birthDay}
                onChange={e => {
                  const v = e.target.value
                  onChange('birthDay', v)
                  onChange('age', calcAge(v))
                }}
              />
            </label>
          </div>

          <div className={styles.row}>
            <label className={styles.label}>
              <span>Last Name</span>
              <input
                className={styles.input}
                value={form.lastName}
                onChange={e => onChange('lastName', e.target.value)}
                autoComplete="family-name"
              />
            </label>

            <label className={styles.label}>
              <span>Age</span>
              <div className={styles.inputGroup}>
                <input className={styles.input} readOnly value={form.age} />
                <span className={styles.suffix}>ปี</span>
              </div>
            </label>
          </div>

          <div className={styles.row}>
            <label className={styles.label}>
              <span>ID Card Number</span>
              <input
                className={styles.input}
                value={form.idcard}
                onChange={e => onChange('idcard', onlyDigits(e.target.value))}
                inputMode="numeric"
                placeholder="กรอก 13 หลัก (ไม่ต้องใส่ขีด)"
              />
            </label>

            <label className={styles.label}>
              <span>Gender</span>
              <select
                className={styles.select}
                value={form.gender}
                onChange={e => onChange('gender', e.target.value)}
              >
                <option value="">Select Gender</option>
                {GENDERS.map(g => <option key={g} value={g}>{g}</option>)}
              </select>
            </label>
          </div>

          <div className={styles.row}>
            <label className={styles.label} style={{ gridColumn: '1 / -1' }}>
              <span>Phone Number</span>
              <input
                className={styles.input}
                value={form.phone}
                onChange={e => onChange('phone', onlyDigits(e.target.value))}
                inputMode="tel"
                placeholder="เช่น 0812345678"
              />
            </label>
          </div>

          <div className={styles.actions}>
            <button className={styles.submit} type="submit" disabled={loading}>
              {loading ? 'Saving…' : 'Add Patient'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
