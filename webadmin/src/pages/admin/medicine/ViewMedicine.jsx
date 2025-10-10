import { useMemo } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/admin/medicine/ViewMedicine.module.css'

// mock data (ภายหลังดึงจาก API ตาม id)
const MOCK = [
  { id: 1, medName: 'Medicine A', genericName: 'Generic A', properties: 'บรรเทาอาการปวด ลดไข้', strength: '500 mg', form: 'ยาเม็ด', unit: 'mg', instruction: 'หลังอาหาร', status: 'Active' },
  { id: 2, medName: 'Medicine B', genericName: 'Generic B', properties: '', strength: '400 mg', form: 'แคปซูล', unit: 'mg', instruction: 'ก่อนอาหาร', status: 'Active' },
  { id: 3, medName: 'Medicine C', genericName: 'Generic C', properties: '', strength: '4 mg/5 ml', form: 'ยาน้ำ', unit: 'ml', instruction: 'วันละ 3 ครั้ง', status: 'Inactive' },
]

const FORM_OPTIONS = ['ยาเม็ด','แคปซูล','ยาน้ำ','ขี้ผึ้ง','สเปรย์']
const UNIT_OPTIONS = ['mg','g','ml','IU','mcg']
const INSTRUCTION_OPTIONS = ['หลังอาหาร','ก่อนอาหาร','พร้อมอาหาร','วันละ 1 ครั้ง','วันละ 2 ครั้ง','วันละ 3 ครั้ง']
const STATUS_OPTIONS = ['Active','Inactive']

export default function ViewMedicine() {
  const { id } = useParams()
  const nav = useNavigate()
  // ของจริง: ใช้ useEffect ไป fetch `/api/medicines/:id`
  const data = useMemo(() => MOCK.find(x => String(x.id) === String(id)), [id])

  if (!data) return <div className={styles.page}>ไม่พบข้อมูลยา</div>

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
              <input className={styles.input} value={data.medName} readOnly />
            </label>

            <label className={styles.label}>
              <span>Generic Name</span>
              <input className={styles.input} value={data.genericName} readOnly />
            </label>

            <label className={styles.label}>
              <span>Properties</span>
              <textarea className={styles.textarea} rows={4} value={data.properties} readOnly />
            </label>

            <label className={styles.label}>
              <span>Strength</span>
              <input className={styles.input} value={data.strength} readOnly />
            </label>
          </div>

          <div className={styles.col}>
            <label className={styles.label}>
              <span>Form</span>
              <select className={styles.select} value={data.form} disabled>
                {[data.form, ...FORM_OPTIONS.filter(o=>o!==data.form)].map(o => <option key={o}>{o}</option>)}
              </select>
            </label>

            <label className={styles.label}>
              <span>Unit</span>
              <select className={styles.select} value={data.unit} disabled>
                {[data.unit, ...UNIT_OPTIONS.filter(o=>o!==data.unit)].map(o => <option key={o}>{o}</option>)}
              </select>
            </label>

            <label className={styles.label}>
              <span>Instruction</span>
              <select className={styles.select} value={data.instruction} disabled>
                {[data.instruction, ...INSTRUCTION_OPTIONS.filter(o=>o!==data.instruction)].map(o => <option key={o}>{o}</option>)}
              </select>
            </label>

            <label className={styles.label}>
              <span>Status</span>
              <select className={styles.select} value={data.status} disabled>
                {[data.status, ...STATUS_OPTIONS.filter(o=>o!==data.status)].map(o => <option key={o}>{o}</option>)}
              </select>
            </label>
          </div>
        </form>
      </div>
    </div>
  )
}
