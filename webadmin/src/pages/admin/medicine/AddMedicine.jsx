import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/admin/medicine/AddMedicine.module.css'
import { createMedicine } from '../../../services/medicines'
import { listForms, listUnitsByForm, listInstructions } from '../../../services/initialData'

export default function AddMedicine() {
  const nav = useNavigate()

  const [form, setForm] = useState({
    med_name: '',
    generic_name: '',
    properties: '',
    strength: '',
    form_id: '',
    unit_id: '',          // optional
    instruction_id: '',   // optional
    med_status: 'active'
  })

  const [forms, setForms] = useState([])
  const [units, setUnits] = useState([])
  const [instructions, setInstructions] = useState([])

  const [loading, setLoading] = useState(false)
  const [unitsLoading, setUnitsLoading] = useState(false)
  const [error, setError] = useState('')
  const [ok, setOk] = useState('')

  const onChange = (k, v) => setForm(s => ({ ...s, [k]: v }))

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      try {
        const [fs, ins] = await Promise.all([ listForms(), listInstructions() ])
        if (cancelled) return
        setForms(Array.isArray(fs) ? fs : [])
        setInstructions(Array.isArray(ins) ? ins : [])
      } catch (e) {
        if (!cancelled) setError(e.message || 'โหลดข้อมูลอ้างอิงไม่สำเร็จ')
      }
    })()
    return () => { cancelled = true }
  }, [])

  // เปลี่ยน form_id → โหลด units ของฟอร์มนั้น (ถ้าไม่เลือก form ก็ไม่ต้องมี units)
  useEffect(() => {
    const fid = Number(form.form_id)
    if (!fid) { setUnits([]); onChange('unit_id', ''); return }
    let cancelled = false
    setUnitsLoading(true)
    ;(async () => {
      try {
        const res = await listUnitsByForm(fid)
        if (cancelled) return
        const arr = Array.isArray(res?.units) ? res.units : []
        setUnits(arr)
        if (!arr.some(u => String(u.id) === String(form.unit_id))) {
          onChange('unit_id', '')
        }
      } catch {
        if (!cancelled) { setUnits([]); onChange('unit_id', '') }
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
    // unit_id & instruction_id ไม่บังคับ ✅
    return ''
  }

  // สร้าง payload โดยใส่เฉพาะคีย์ที่มีค่า
  const buildPayload = () => {
    const payload = {
      med_name: form.med_name.trim(),
      generic_name: form.generic_name.trim(),
      properties: form.properties.trim(),
      strength: form.strength.trim(),
      form_id: Number(form.form_id),
      med_status: form.med_status,
    }
    if (form.unit_id)        payload.unit_id = Number(form.unit_id)
    if (form.instruction_id) payload.instruction_id = Number(form.instruction_id)
    return payload
  }

  const onSubmit = async (e) => {
    e.preventDefault()
    if (loading) return
    const v = validate()
    if (v) { setError(v); setOk(''); return }
    setError(''); setOk(''); setLoading(true)
    try {
      await createMedicine(buildPayload())
      setOk('เพิ่มยาเรียบร้อย')
      nav('/admin/medicine-info', { replace: true })
    } catch (err) {
      setError(err.message || 'บันทึกไม่สำเร็จ')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h2 className={styles.title}>Add Medicine</h2>
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
                placeholder="เช่น 500 mg หรือ 4 mg/5 ml"
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
                {instructions.map(x => (
                  <option key={x.id} value={x.id}>
                    {x.instruction_name || x.name}
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
            <button className={styles.submit} type="submit" disabled={loading}>
              {loading ? 'Saving…' : 'Add Medicine'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
