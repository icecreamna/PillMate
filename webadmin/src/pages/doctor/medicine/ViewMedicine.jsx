import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/medicine/ViewMedicine.module.css'

const MOCK = [
  { id:1, medName:'Medicine A', genericName:'Generic A', strength:'500 mg', form:'ยาเม็ด', properties:'บรรเทาอาการปวด ลดไข้', unit:'mg', instruction:'หลังอาหาร', status:'Active' },
  { id:2, medName:'Medicine B', genericName:'Generic B', strength:'400 mg', form:'แคปซูล', properties:'', unit:'mg', instruction:'ก่อนอาหาร', status:'Active' },
  { id:3, medName:'Medicine C', genericName:'Generic C', strength:'4 mg/5 ml', form:'ยาน้ำ', properties:'', unit:'ml', instruction:'วันละ 3 ครั้ง', status:'Inactive' },
]

export default function ViewMedicine(){
  const { id } = useParams()
  const nav = useNavigate()
  const data = useMemo(()=>MOCK.find(x=>String(x.id)===String(id)),[id])
  const [item, setItem] = useState(null)

  useEffect(()=>{ setItem(data || null) }, [data])

  if(!item) return <div className={styles.page}>ไม่พบข้อมูลยา</div>

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
              <input className={styles.input} value={item.medName} disabled />
            </label>
            <label className={styles.label}>
              <span>Generic Name</span>
              <input className={styles.input} value={item.genericName} disabled />
            </label>
            <label className={styles.label}>
              <span>Properties</span>
              <textarea className={styles.textarea} rows={4} value={item.properties} disabled />
            </label>
            <label className={styles.label}>
              <span>Strength</span>
              <input className={styles.input} value={item.strength} disabled />
            </label>
          </div>

          <div className={styles.col}>
            <label className={styles.label}>
              <span>Form</span>
              <input className={styles.input} value={item.form} disabled />
            </label>
            <label className={styles.label}>
              <span>Unit</span>
              <input className={styles.input} value={item.unit} disabled />
            </label>
            <label className={styles.label}>
              <span>Instruction</span>
              <input className={styles.input} value={item.instruction} disabled />
            </label>
            <label className={styles.label}>
              <span>Status</span>
              <input className={styles.input} value={item.status} disabled />
            </label>
          </div>
        </div>
      </div>
    </div>
  )
}
