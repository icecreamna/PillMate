import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/patient/EditPatient.module.css'

const GENDERS = ['ชาย','หญิง','ไม่ระบุ']

function calcAge(dobStr) {
  if (!dobStr) return ''
  const dob = new Date(dobStr)
  if (isNaN(dob)) return ''
  const t = new Date()
  let a = t.getFullYear() - dob.getFullYear()
  const m = t.getMonth() - dob.getMonth()
  if (m < 0 || (m === 0 && t.getDate() < dob.getDate())) a--
  return String(Math.max(a,0))
}

export default function EditPatient(){
  const { id } = useParams()
  const nav = useNavigate()
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')
  const [ok, setOk] = useState('')

  const [form, setForm] = useState({
    firstName:'', lastName:'', birthDay:'', age:'',
    gender:'', idcard:'', phone:''
  })

  // โหลดข้อมูลจริงตาม id (mock ตอนนี้)
  useEffect(() => {
    // TODO: fetch(`/api/patients/${id}`)
    const mock = {
      firstName:'สมชาย', lastName:'โชติ',
      birthDay:'1995-08-24', gender:'ชาย',
      idcard:'1234567890100', phone:'0811234567'
    }
    setForm({ ...mock, age: calcAge(mock.birthDay) })
  }, [id])

  // เปลี่ยนค่าในฟอร์ม
  const onChange = (k, v) => setForm(s => ({ ...s, [k]: v }))

  const validate = () => {
    if (!form.firstName.trim()) return 'กรุณากรอก First Name'
    if (!form.lastName.trim())  return 'กรุณากรอก Last Name'
    if (!form.birthDay)         return 'กรุณาเลือกวันเกิด'
    if (!form.gender)           return 'กรุณาเลือกเพศ'
    return ''
  }

  const onSubmit = async (e) => {
    e.preventDefault()
    const v = validate()
    if (v) { setError(v); setOk(''); return }
    setError(''); setOk(''); setSaving(true)
    try {
      // TODO: PUT /api/patients/:id  (ส่ง form ที่มี age ซึ่งคำนวณจาก birthDay)
      // await fetch(`/api/patients/${id}`, { method:'PUT', headers:{'Content-Type':'application/json'}, body: JSON.stringify(form) })
      await new Promise(r => setTimeout(r, 600)) // mock
      setOk('บันทึกสำเร็จ')
      nav('/doc/patients', { replace:true })
    } catch(err){
      setError(err.message || 'บันทึกไม่สำเร็จ')
    } finally {
      setSaving(false)
    }
  }

  const todayStr = new Date().toISOString().slice(0,10)

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h2 className={styles.title}>Edit Patient</h2>
        <button className={styles.back} onClick={()=>nav(-1)}>← Back</button>
      </div>
      <hr className={styles.hr} />

      <div className={styles.card}>
        {error && <div className={styles.error}>{error}</div>}
        {ok && <div className={styles.success}>{ok}</div>}

        <form className={styles.form} onSubmit={onSubmit}>
          <div className={styles.row}>
            <label className={styles.label}>
              <span>First Name</span>
              <input className={styles.input}
                     value={form.firstName}
                     onChange={e=>onChange('firstName', e.target.value)} />
            </label>

            <label className={styles.label}>
              <span>Birth Day</span>
              <input className={styles.input} type="date" max={todayStr}
                     value={form.birthDay}
                     onChange={e=>{
                       const v = e.target.value
                       onChange('birthDay', v)
                       onChange('age', calcAge(v))
                     }} />
            </label>
          </div>

          <div className={styles.row}>
            <label className={styles.label}>
              <span>Last Name</span>
              <input className={styles.input}
                     value={form.lastName}
                     onChange={e=>onChange('lastName', e.target.value)} />
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
              <input className={styles.input}
                     value={form.idcard}
                     onChange={e=>onChange('idcard', e.target.value)} />
            </label>

            <label className={styles.label}>
              <span>Gender</span>
              <select className={styles.select}
                      value={form.gender}
                      onChange={e=>onChange('gender', e.target.value)}>
                <option value="">Select Gender</option>
                {GENDERS.map(g => <option key={g} value={g}>{g}</option>)}
              </select>
            </label>
          </div>

          <div className={styles.row}>
            <label className={styles.label} style={{gridColumn:'1 / -1'}}>
              <span>Phone Number</span>
              <input className={styles.input}
                     value={form.phone}
                     onChange={e=>onChange('phone', e.target.value)} />
            </label>
          </div>

          <div className={styles.actions}>
            <button className={styles.submit} type="submit" disabled={saving}>
              {saving ? 'Saving…' : 'Save Changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
