import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/admin/medicine/EditMedicine.module.css'

const FORM_OPTIONS = ['ยาเม็ด','แคปซูล','ยาน้ำ','ขี้ผึ้ง','สเปรย์']
const UNIT_OPTIONS = ['เม็ด','mg','g','ml','IU','mcg']
const INSTRUCTION_OPTIONS = ['หลังอาหาร','ก่อนอาหาร','พร้อมอาหาร','วันละ 1 ครั้ง','วันละ 2 ครั้ง','วันละ 3 ครั้ง']
const STATUS_OPTIONS = ['Active','Inactive']

export default function EditMedicine() {
  const { id } = useParams()
  const nav = useNavigate()

  const [form, setForm] = useState({
    medName:'', genericName:'', properties:'', strength:'',
    form:'', unit:'', instruction:'', status:'Active',
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [ok, setOk] = useState('')

  // โหลดข้อมูลจริงด้วย id (ตอนนี้ mock)
  useEffect(() => {
    // TODO: fetch(`/api/medicines/${id}`)
    const mock = { medName:'Medicine A', genericName:'Generic A', properties:'บรรเทาอาการปวด ลดไข้', strength:'500 mg', form:'ยาเม็ด', unit:'เม็ด', instruction:'หลังอาหาร', status:'Active' }
    setForm(mock)
  }, [id])

  const onChange = (k, v) => setForm(s => ({ ...s, [k]: v }))

  const validate = () => {
    if (!form.medName.trim()) return 'กรุณากรอก Medicine Name'
    if (!form.genericName.trim()) return 'กรุณากรอก Generic Name'
    if (!form.form) return 'กรุณาเลือก Form'
    if (!form.unit) return 'กรุณาเลือก Unit'
    return ''
  }

  const onSubmit = async (e) => {
    e.preventDefault()
    const v = validate()
    if (v) { setError(v); setOk(''); return }
    setError(''); setOk(''); setLoading(true)
    try {
      // TODO: PUT /api/medicines/:id
      // await fetch(`/api/medicines/${id}`, { method:'PUT', headers:{'Content-Type':'application/json'}, body: JSON.stringify(form), credentials:'include' })
      await new Promise(r => setTimeout(r, 600)) // mock
      setOk('บันทึกสำเร็จ')
      nav('/admin/medicine-info', { replace:true })
    } catch (err) {
      setError(err.message || 'บันทึกไม่สำเร็จ')
    } finally { setLoading(false) }
  }

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
              <input className={styles.input} value={form.medName} onChange={e=>onChange('medName', e.target.value)} />
            </label>

            <label className={styles.label}>
              <span>Generic Name</span>
              <input className={styles.input} value={form.genericName} onChange={e=>onChange('genericName', e.target.value)} />
            </label>

            <label className={styles.label}>
              <span>Properties</span>
              <textarea className={styles.textarea} rows={4} value={form.properties} onChange={e=>onChange('properties', e.target.value)} />
            </label>

            <label className={styles.label}>
              <span>Strength</span>
              <input className={styles.input} value={form.strength} onChange={e=>onChange('strength', e.target.value)} />
            </label>
          </div>

          <div className={styles.col}>
            <label className={styles.label}>
              <span>Form</span>
              <select className={styles.select} value={form.form} onChange={e=>onChange('form', e.target.value)}>
                <option value="">Select Form</option>
                {FORM_OPTIONS.map(o => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>

            <label className={styles.label}>
              <span>Unit</span>
              <select className={styles.select} value={form.unit} onChange={e=>onChange('unit', e.target.value)}>
                <option value="">Select Unit</option>
                {UNIT_OPTIONS.map(o => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>

            <label className={styles.label}>
              <span>Instruction</span>
              <select className={styles.select} value={form.instruction} onChange={e=>onChange('instruction', e.target.value)}>
                <option value="">Select Instruction</option>
                {INSTRUCTION_OPTIONS.map(o => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>

            <label className={styles.label}>
              <span>Status</span>
              <select className={styles.select} value={form.status} onChange={e=>onChange('status', e.target.value)}>
                {STATUS_OPTIONS.map(o => <option key={o} value={o}>{o}</option>)}
              </select>
            </label>
          </div>

          <div className={styles.actions}>
            <button className={styles.submit} type="submit" disabled={loading}>
              {loading ? 'Saving…' : 'Save Changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
