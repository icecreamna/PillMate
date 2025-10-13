import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/admin/medicine/EditMedicine.module.css'
import { getMedicine, updateMedicine } from '../../../services/medicines'
import { listForms, listUnitsByForm, listInstructions } from '../../../services/initialData'

export default function EditMedicine() {
  const { id } = useParams()
  const nav = useNavigate()

  // ฟอร์ม (แม็ปตามสเปก BE)
  const [form, setForm] = useState({
    med_name:'', generic_name:'', properties:'', strength:'',
    form_id:'', unit_id:'', instruction_id:'', med_status:'active',
  })

  // ออปชันจาก API
  const [forms, setForms] = useState([])               // [{id, form_name}]
  const [units, setUnits] = useState([])               // [{id, unit_name}] (ขึ้นกับ form_id)
  const [instructions, setInstructions] = useState([]) // [{id, instruction_name}]

  const [loading, setLoading] = useState(true)         // โหลดหน้า
  const [saving, setSaving] = useState(false)          // กำลังบันทึก
  const [unitsLoading, setUnitsLoading] = useState(false)
  const [error, setError] = useState('')
  const [ok, setOk] = useState('')

  const onChange = (k, v) => setForm(s => ({ ...s, [k]: v }))

  // โหลดข้อมูลเริ่มต้น: ยา + forms + instructions
  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoading(true); setError(''); setOk('')
      try {
        const [resMed, fs, ins] = await Promise.all([
          getMedicine(id),      // GET /admin/medicine-info/:id -> { data: {...} }
          listForms(),          // GET /forms
          listInstructions(),   // GET /instructions
        ])
        if (cancelled) return
        const m = resMed?.data
        if (!m) throw new Error('ไม่พบข้อมูลยา')

        setForms(Array.isArray(fs) ? fs : [])
        setInstructions(Array.isArray(ins) ? ins : [])

        // ใส่ค่าเดิมลงฟอร์ม (รองรับทั้ง *_name และ *_id)
        setForm({
          med_name:      m.med_name || '',
          generic_name:  m.generic_name || '',
          properties:    m.properties || '',
          strength:      m.strength || '',
          form_id:       m.form_id ?? '',            // ต้องเป็น id
          unit_id:       m.unit_id ?? '',            // optional
          instruction_id:m.instruction_id ?? '',     // optional
          med_status:    (m.med_status?.toLowerCase() === 'inactive') ? 'inactive' : 'active',
        })
      } catch (e) {
        if (!cancelled) setError(e.message || 'โหลดข้อมูลไม่สำเร็จ')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => { cancelled = true }
  }, [id])

  // เมื่อเปลี่ยน form_id → โหลด units ของฟอร์มนั้น
  useEffect(() => {
    const fid = Number(form.form_id)
    if (!fid) { setUnits([]); onChange('unit_id',''); return }
    let cancelled = false
    setUnitsLoading(true)
    ;(async () => {
      try {
        const res = await listUnitsByForm(fid) // GET /forms/:id/units -> {units:[{id,unit_name}]}
        if (cancelled) return
        const arr = Array.isArray(res?.units) ? res.units : []
        setUnits(arr)
        // ถ้า unit เดิมไม่อยู่ในลิสต์ -> เคลียร์
        if (!arr.some(u => String(u.id) === String(form.unit_id))) {
          onChange('unit_id','')
        }
      } catch {
        if (!cancelled) { setUnits([]); onChange('unit_id','') }
      } finally {
        if (!cancelled) setUnitsLoading(false)
      }
    })()
    return () => { cancelled = true }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [form.form_id])

  const validate = () => {
    if (!form.med_name.trim()) return 'กรุณากรอก Medicine Name'
    if (!form.generic_name.trim()) return 'กรุณากรอก Generic Name'
    if (!form.form_id) return 'กรุณาเลือก Form'
    // unit_id / instruction_id ไม่บังคับ
    return ''
  }

  // สร้าง payload — แนบเฉพาะฟิลด์ที่มีค่า (เลี่ยงส่ง 0/null)
  const buildPayload = () => {
    const payload = {
      med_name: form.med_name.trim(),
      generic_name: form.generic_name.trim(),
      properties: form.properties.trim(),
      strength: form.strength.trim(),
      form_id: Number(form.form_id),
      med_status: form.med_status, // 'active' | 'inactive'
    }
    if (form.unit_id)        payload.unit_id = Number(form.unit_id)
    if (form.instruction_id) payload.instruction_id = Number(form.instruction_id)
    return payload
  }

  const onSubmit = async (e) => {
    e.preventDefault()
    if (saving) return
    const v = validate()
    if (v) { setError(v); setOk(''); return }
    setError(''); setOk(''); setSaving(true)
    try {
      await updateMedicine(id, buildPayload()) // PUT /admin/medicine-info/:id
      setOk('บันทึกสำเร็จ')
      nav('/admin/medicine-info', { replace:true })
    } catch (err) {
      setError(err.message || 'บันทึกไม่สำเร็จ')
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <div className={styles.page}>กำลังโหลด...</div>
  if (error)   return <div className={styles.page}>{error}</div>

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h2 className={styles.title}>Edit Medicine</h2>
        <button className={styles.back} onClick={()=>nav(-1)}>← Back</button>
      </div>
      <hr className={styles.hr} />

      <div className={styles.card}>
        {error && <div className={styles.error}>{error}</div>}
        {ok && <div className={styles.success}>{ok}</div>}

        <form onSubmit={onSubmit} className={styles.form}>
          <div className={styles.col}>
            <label className={styles.label}>
              <span>Medicine Name</span>
              <input
                className={styles.input}
                value={form.med_name}
                onChange={e=>onChange('med_name', e.target.value)}
              />
            </label>

            <label className={styles.label}>
              <span>Generic Name</span>
              <input
                className={styles.input}
                value={form.generic_name}
                onChange={e=>onChange('generic_name', e.target.value)}
              />
            </label>

            <label className={styles.label}>
              <span>Properties</span>
              <textarea
                className={styles.textarea}
                rows={4}
                value={form.properties}
                onChange={e=>onChange('properties', e.target.value)}
              />
            </label>

            <label className={styles.label}>
              <span>Strength</span>
              <input
                className={styles.input}
                value={form.strength}
                onChange={e=>onChange('strength', e.target.value)}
              />
            </label>
          </div>

          <div className={styles.col}>
            <label className={styles.label}>
              <span>Form</span>
              <select
                className={styles.select}
                value={form.form_id}
                onChange={e=>onChange('form_id', e.target.value)}
              >
                <option value="">Select Form</option>
                {forms.map(f => (
                  <option key={f.id} value={f.id}>
                    {f.form_name || f.name}
                  </option>
                ))}
              </select>
            </label>

            <label className={styles.label}>
              <span>Unit (optional)</span>
              <select
                className={styles.select}
                value={form.unit_id}
                onChange={e=>onChange('unit_id', e.target.value)}
                disabled={!form.form_id || unitsLoading}
              >
                <option value="">
                  {!form.form_id ? 'เลือกฟอร์มก่อน' : (unitsLoading ? 'กำลังโหลด...' : 'ไม่ระบุ / Select Unit')}
                </option>
                {units.map(u => (
                  <option key={u.id} value={u.id}>
                    {u.unit_name || u.name}
                  </option>
                ))}
              </select>
            </label>

            <label className={styles.label}>
              <span>Instruction (optional)</span>
              <select
                className={styles.select}
                value={form.instruction_id}
                onChange={e=>onChange('instruction_id', e.target.value)}
              >
                <option value="">ไม่ระบุ / Select Instruction</option>
                {instructions.map(i => (
                  <option key={i.id} value={i.id}>
                    {i.instruction_name || i.name}
                  </option>
                ))}
              </select>
            </label>

            <label className={styles.label}>
              <span>Status</span>
              <select
                className={styles.select}
                value={form.med_status}
                onChange={e=>onChange('med_status', e.target.value)}
              >
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
              </select>
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
