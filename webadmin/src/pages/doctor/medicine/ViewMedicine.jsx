import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/medicine/ViewMedicine.module.css'
import { getMedicineForDoctor } from '../../../services/medicines'
import { listForms, listUnitsByForm, listInstructions } from '../../../services/initialData'

export default function ViewMedicine(){
  const { id } = useParams()
  const nav = useNavigate()

  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [med, setMed] = useState(null)

  // อ้างอิงชื่อ (ไว้แปลง id -> label ถ้า DTO ไม่มี name ติดมา)
  const [forms, setForms] = useState([])
  const [units, setUnits] = useState([])
  const [instructions, setInstructions] = useState([])

  // โหลดข้อมูลยา + อ้างอิง (forms, instructions)
  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoading(true); setError('')
      try {
        const [resMed, fs, ins] = await Promise.all([
          getMedicineForDoctor(id),
          listForms(),
          listInstructions()
        ])
        if (cancelled) return
        const data = resMed?.data
        if (!data) throw new Error('ไม่พบข้อมูลยา')
        setMed(data)
        setForms(Array.isArray(fs) ? fs : [])
        setInstructions(Array.isArray(ins) ? ins : [])
      } catch (e) {
        if (!cancelled) setError(e.message || 'โหลดข้อมูลไม่สำเร็จ')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => { cancelled = true }
  }, [id])

  // เมื่อมี form_id -> โหลด units สำหรับฟอร์มนั้น
  useEffect(() => {
    let cancelled = false
    const fid = Number(med?.form_id)
    if (!fid) { setUnits([]); return }
    ;(async () => {
      try {
        const res = await listUnitsByForm(fid) // { form_id, units: [{id, unit_name}] }
        if (!cancelled) setUnits(Array.isArray(res?.units) ? res.units : [])
      } catch {
        if (!cancelled) setUnits([])
      }
    })()
    return () => { cancelled = true }
  }, [med?.form_id])

  // ตัวช่วยหยิบชื่อแบบยืดหยุ่น
  const formName = useMemo(() => {
    if (!med) return ''
    if (med.form_name) return med.form_name
    const f = forms.find(x => String(x.id) === String(med.form_id))
    return f?.form_name || f?.name || ''
  }, [med, forms])

  const unitName = useMemo(() => {
    if (!med) return ''
    if (med.unit_name) return med.unit_name
    const u = units.find(x => String(x.id) === String(med.unit_id))
    return u?.unit_name || u?.name || ''
  }, [med, units])

  const instructionName = useMemo(() => {
    if (!med) return ''
    if (med.instruction_name) return med.instruction_name
    const ins = instructions.find(x => String(x.id) === String(med.instruction_id))
    return ins?.instruction_name || ins?.name || ''
  }, [med, instructions])

  if (loading) return <div className={styles.page}>กำลังโหลด...</div>
  if (error)   return <div className={styles.page}>{error}</div>
  if (!med)    return <div className={styles.page}>ไม่พบข้อมูลยา</div>

  const statusLabel = (med.med_status?.toLowerCase() === 'inactive') ? 'Inactive' : 'Active'

  const safe = (v) => (v == null || v === '' ? '-' : v)

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h2 className={styles.title}>View Medicine</h2>
        <button className={styles.back} onClick={()=>nav(-1)}>← Back</button>
      </div>
      <hr className={styles.hr} />

      <div className={styles.card}>
        <div className={styles.form}>
          <div className={styles.col}>
            <label className={styles.label}>
              <span>Medicine Name</span>
              <input className={styles.input} value={safe(med.med_name)} disabled />
            </label>
            <label className={styles.label}>
              <span>Generic Name</span>
              <input className={styles.input} value={safe(med.generic_name)} disabled />
            </label>
            <label className={styles.label}>
              <span>Properties</span>
              <textarea className={styles.textarea} rows={4} value={safe(med.properties)} disabled />
            </label>
            <label className={styles.label}>
              <span>Strength</span>
              <input className={styles.input} value={safe(med.strength)} disabled />
            </label>
          </div>

          <div className={styles.col}>
            <label className={styles.label}>
              <span>Form</span>
              <input className={styles.input} value={safe(formName)} disabled />
            </label>
            <label className={styles.label}>
              <span>Unit</span>
              <input className={styles.input} value={safe(unitName)} disabled />
            </label>
            <label className={styles.label}>
              <span>Instruction</span>
              <input className={styles.input} value={safe(instructionName)} disabled />
            </label>
            <label className={styles.label}>
              <span>Status</span>
              <input className={styles.input} value={statusLabel} disabled />
            </label>
          </div>
        </div>
      </div>
    </div>
  )
}
