import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/admin/medicine/ViewMedicine.module.css'
import { getMedicine } from '../../../services/medicines'
import { listForms, listUnitsByForm, listInstructions } from '../../../services/initialData'

export default function ViewMedicine() {
  const { id } = useParams()
  const nav = useNavigate()

  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // ข้อมูลยา (จาก DTO)
  const [med, setMed] = useState(null)

  // อ้างอิงชื่อ
  const [forms, setForms] = useState([])            // [{id, form_name}]
  const [units, setUnits] = useState([])            // [{id, unit_name}] (ตาม form_id)
  const [instructions, setInstructions] = useState([]) // [{id, instruction_name}]

  // โหลดหลัก: medicine + forms + instructions
  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoading(true); setError('')
      try {
        const [resMed, fs, ins] = await Promise.all([
          getMedicine(id),   // -> { data: {...} }
          listForms(),       // -> array
          listInstructions() // -> array
        ])
        if (cancelled) return
        const m = resMed?.data
        if (!m) throw new Error('ไม่พบข้อมูลยา')
        setMed(m)
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

  // เมื่อมี form_id → โหลด units ของฟอร์มนั้น
  useEffect(() => {
    let cancelled = false
    const fid = Number(med?.form_id)
    if (!fid) { setUnits([]); return }
    ;(async () => {
      try {
        const res = await listUnitsByForm(fid) // -> { form_id, units: [...] }
        if (!cancelled) setUnits(Array.isArray(res?.units) ? res.units : [])
      } catch {
        if (!cancelled) setUnits([])
      }
    })()
    return () => { cancelled = true }
  }, [med?.form_id])

  // ---------- helpers ----------
  // หาชื่อจาก id หรือจากฟิลด์ตรง ๆ
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

  // ทำรายการ option โดยรับ label ปัจจุบัน + รายการ label ทั้งหมด
  const ensureOptions = (currentLabel, allLabels) => {
    const labels = (allLabels || []).filter(Boolean)
    if (currentLabel) {
      const rest = labels.filter(l => l !== currentLabel)
      return [{ value: currentLabel, label: currentLabel }, ...rest.map(l => ({ value: l, label: l }))]
    }
    // ไม่มีค่า -> ใส่ placeholder value="" ก่อน
    return [{ value: "", label: "-" }, ...labels.map(l => ({ value: l, label: l }))]
  }

  if (loading) return <div className={styles.page}>กำลังโหลด...</div>
  if (error)   return <div className={styles.page}>{error}</div>
  if (!med)    return <div className={styles.page}>ไม่พบข้อมูลยา</div>

  // options สำหรับ select (disabled)
  const formOpts         = ensureOptions(formName,        forms.map(f => f.form_name || f.name))
  const unitOpts         = ensureOptions(unitName,        units.map(u => u.unit_name || u.name))
  const instructionOpts  = ensureOptions(instructionName, instructions.map(i => i.instruction_name || i.name))
  const statusLabel      = (med.med_status?.toLowerCase() === 'inactive') ? 'Inactive' : 'Active'
  const statusOpts       = ensureOptions(statusLabel, ['Active', 'Inactive'])

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h2 className={styles.title}>View Medicine</h2>
        <button className={styles.back} onClick={() => nav(-1)}>← Back</button>
      </div>
      <hr className={styles.hr} />

      <div className={styles.card}>
        <form className={styles.form} onSubmit={e=>e.preventDefault()}>
          <div className={styles.col}>
            <label className={styles.label}>
              <span>Medicine Name</span>
              <input className={styles.input} value={med.med_name || ''} readOnly />
            </label>

            <label className={styles.label}>
              <span>Generic Name</span>
              <input className={styles.input} value={med.generic_name || ''} readOnly />
            </label>

            <label className={styles.label}>
              <span>Properties</span>
              <textarea className={styles.textarea} rows={4} value={med.properties || ''} readOnly />
            </label>

            <label className={styles.label}>
              <span>Strength</span>
              <input className={styles.input} value={med.strength || ''} readOnly />
            </label>
          </div>

          <div className={styles.col}>
            <label className={styles.label}>
              <span>Form</span>
              <select className={styles.select} value={formName || ""} disabled>
                {formOpts.map(o => <option key={`${o.value}-${o.label}`} value={o.value}>{o.label}</option>)}
              </select>
            </label>

            <label className={styles.label}>
              <span>Unit</span>
              <select className={styles.select} value={unitName || ""} disabled>
                {unitOpts.map(o => <option key={`${o.value}-${o.label}`} value={o.value}>{o.label}</option>)}
              </select>
            </label>

            <label className={styles.label}>
              <span>Instruction</span>
              <select className={styles.select} value={instructionName || ""} disabled>
                {instructionOpts.map(o => <option key={`${o.value}-${o.label}`} value={o.value}>{o.label}</option>)}
              </select>
            </label>

            <label className={styles.label}>
              <span>Status</span>
              <select className={styles.select} value={statusLabel || ""} disabled>
                {statusOpts.map(o => <option key={`${o.value}-${o.label}`} value={o.value}>{o.label}</option>)}
              </select>
            </label>
          </div>
        </form>
      </div>
    </div>
  )
}
