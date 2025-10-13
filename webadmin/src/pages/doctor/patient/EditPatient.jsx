import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/patient/EditPatient.module.css'
import { getPatient, updatePatient } from '../../../services/patients'

const GENDERS = ['ชาย', 'หญิง'] // BE รองรับแค่ 2 ค่า

const onlyDigits = (s) => (s || '').replace(/[^\d]/g, '')

function calcAge(ymd) {
  if (!ymd) return ''
  const dob = new Date(`${ymd}T00:00:00`)
  if (isNaN(dob)) return ''
  const t = new Date()
  let a = t.getFullYear() - dob.getFullYear()
  const m = t.getMonth() - dob.getMonth()
  if (m < 0 || (m === 0 && t.getDate() < dob.getDate())) a--
  return String(Math.max(a, 0))
}

// แปลง YYYY-MM-DD -> RFC3339 (+07:00) เพื่อส่งให้ BE
const toBangkokISO = (ymd) => (ymd ? `${ymd}T00:00:00+07:00` : undefined)

export default function EditPatient() {
  const { id } = useParams()
  const nav = useNavigate()

  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')
  const [ok, setOk] = useState('')

  // เก็บ DTO เดิมจาก BE (ไว้เทียบว่าแก้อะไรไปบ้าง)
  const [orig, setOrig] = useState(null)

  // ฟอร์ม UI
  const [form, setForm] = useState({
    firstName: '', lastName: '', birthDay: '', age: '',
    gender: '', idcard: '', phone: ''
  })

  // โหลดข้อมูลจริงตาม id
  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoading(true); setError(''); setOk('')
      try {
        const res = await getPatient(id) // -> { data: {...} }
        if (cancelled) return
        const p = res?.data
        if (!p) throw new Error('ไม่พบข้อมูลผู้ป่วย')

        setOrig(p)
        const mapped = {
          firstName: p.first_name || '',
          lastName:  p.last_name  || '',
          birthDay:  p.birth_day  || '',          // "YYYY-MM-DD"
          gender:    p.gender     || '',
          idcard:    p.id_card_number || '',
          phone:     p.phone_number   || '',
        }
        mapped.age = calcAge(mapped.birthDay)
        setForm(mapped)
      } catch (e) {
        if (!cancelled) setError(e.message || 'โหลดข้อมูลไม่สำเร็จ')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => { cancelled = true }
  }, [id])

  const onChange = (k, v) => setForm(s => ({ ...s, [k]: v }))

  const validate = () => {
    if (!form.firstName.trim()) return 'กรุณากรอก First Name'
    if (!form.lastName.trim())  return 'กรุณากรอก Last Name'
    if (!form.birthDay)         return 'กรุณาเลือกวันเกิด'
    if (!form.gender)           return 'กรุณาเลือกเพศ'
    const id13 = onlyDigits(form.idcard)
    if (id13 && id13.length !== 13) return 'เลขบัตรประชาชนต้องมี 13 หลัก'
    const tel = onlyDigits(form.phone)
    if (tel && tel.length !== 10)   return 'เบอร์โทรต้องมี 10 หลัก'
    if (form.gender && !GENDERS.includes(form.gender)) return 'เพศไม่ถูกต้อง'
    return ''
  }

  // สร้าง payload ที่ “ส่งเฉพาะฟิลด์ที่เปลี่ยน”
  const buildPatch = () => {
    const patch = {}
    if (!orig) return patch

    if (form.firstName.trim() !== (orig.first_name || ''))
      patch.first_name = form.firstName.trim()

    if (form.lastName.trim() !== (orig.last_name || ''))
      patch.last_name = form.lastName.trim()

    if (form.birthDay !== (orig.birth_day || ''))
      patch.birth_day = toBangkokISO(form.birthDay)

    // ถ้าผู้ใช้แก้เลขบัตร/เบอร์ → ส่งไปตรวจที่ BE (BE จะ validate length + duplicate)
    const id13 = onlyDigits(form.idcard)
    if (id13 !== (orig.id_card_number || ''))
      patch.id_card_number = id13 || undefined

    const tel = onlyDigits(form.phone)
    if (tel !== (orig.phone_number || ''))
      patch.phone_number = tel || undefined

    if (form.gender !== (orig.gender || ''))
      patch.gender = form.gender

    // ลบ key ที่เป็น undefined ออก (กันส่งค่าแปลกไป)
    Object.keys(patch).forEach(k => {
      if (patch[k] === undefined) delete patch[k]
    })

    return patch
  }

  const onSubmit = async (e) => {
    e.preventDefault()
    const v = validate()
    if (v) { setError(v); setOk(''); return }
    const patch = buildPatch()
    if (Object.keys(patch).length === 0) {
      setOk('ไม่มีการเปลี่ยนแปลง'); setError(''); return
    }

    setError(''); setOk(''); setSaving(true)
    try {
      await updatePatient(id, patch) // PUT /doctor/hospital-patients/:id
      setOk('บันทึกสำเร็จ')
      nav('/doc/patients', { replace: true })
    } catch (err) {
      setError(err.message || 'บันทึกไม่สำเร็จ')
    } finally {
      setSaving(false)
    }
  }

  const todayStr = useMemo(() => new Date().toISOString().slice(0,10), [])

  if (loading) return <div className={styles.page}>กำลังโหลด...</div>

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
              <input
                className={styles.input}
                value={form.firstName}
                onChange={e=>onChange('firstName', e.target.value)}
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
                onChange={e=>{
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
                onChange={e=>onChange('lastName', e.target.value)}
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
                onChange={e=>onChange('idcard', onlyDigits(e.target.value))}
                inputMode="numeric"
                placeholder="กรอก 13 หลัก (ไม่ต้องใส่ขีด)"
              />
            </label>

            <label className={styles.label}>
              <span>Gender</span>
              <select
                className={styles.select}
                value={form.gender}
                onChange={e=>onChange('gender', e.target.value)}
              >
                <option value="">Select Gender</option>
                {GENDERS.map(g => <option key={g} value={g}>{g}</option>)}
              </select>
            </label>
          </div>

          <div className={styles.row}>
            <label className={styles.label} style={{gridColumn:'1 / -1'}}>
              <span>Phone Number</span>
              <input
                className={styles.input}
                value={form.phone}
                onChange={e=>onChange('phone', onlyDigits(e.target.value))}
                inputMode="tel"
                placeholder="เช่น 0812345678"
              />
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
