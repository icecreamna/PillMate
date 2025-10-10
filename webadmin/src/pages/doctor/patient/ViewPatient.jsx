import { useEffect, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/patient/ViewPatient.module.css'

const GENDERS = ['ชาย','หญิง','ไม่ระบุ']

function calcAge(dobStr){
  if(!dobStr) return ''
  const dob = new Date(dobStr)
  const t = new Date()
  let a = t.getFullYear() - dob.getFullYear()
  const m = t.getMonth() - dob.getMonth()
  if (m < 0 || (m === 0 && t.getDate() < dob.getDate())) a--
  return String(Math.max(a,0))
}

export default function ViewPatient(){
  const { id } = useParams()
  const nav = useNavigate()

  const [p, setP] = useState({
    firstName:'', lastName:'', birthDay:'', gender:'', idcard:'', phone:''
  })
  const [age, setAge] = useState('')

  useEffect(() => {
    // TODO: เรียก API จริง: /api/patients/:id
    const mock = { firstName:'สมชาย', lastName:'โชติ', birthDay:'1995-08-24', gender:'ชาย', idcard:'1234567890100', phone:'0811234567' }
    setP(mock)
  }, [id])

  useEffect(() => { setAge(calcAge(p.birthDay)) }, [p.birthDay])

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h2 className={styles.title}>View Patient</h2>
        <button className={styles.back} onClick={()=>nav(-1)}>← Back</button>
      </div>
      <hr className={styles.hr} />

      <div className={styles.card}>
        <div className={styles.row}>
          <label className={styles.label}>
            <span>First Name</span>
            <input className={styles.input} value={p.firstName} disabled />
          </label>
          <label className={styles.label}>
            <span>Birth Day</span>
            <input className={styles.input} type="date" value={p.birthDay} disabled />
          </label>
        </div>

        <div className={styles.row}>
          <label className={styles.label}>
            <span>Last Name</span>
            <input className={styles.input} value={p.lastName} disabled />
          </label>
          <label className={styles.label}>
            <span>Age</span>
            <div className={styles.inputGroup}>
              <input className={styles.input} value={age} disabled />
              <span className={styles.suffix}>ปี</span>
            </div>
          </label>
        </div>

        <div className={styles.row}>
          <label className={styles.label}>
            <span>ID Card Number</span>
            <input className={styles.input} value={p.idcard} disabled />
          </label>
          <label className={styles.label}>
            <span>Gender</span>
            <select className={styles.select} value={p.gender} disabled>
              {[p.gender, ...GENDERS.filter(g=>g!==p.gender)].map(g=>(
                <option key={g}>{g}</option>
              ))}
            </select>
          </label>
        </div>

        <div className={styles.row}>
          <label className={styles.label} style={{gridColumn:'1 / -1'}}>
            <span>Phone Number</span>
            <input className={styles.input} value={p.phone} disabled />
          </label>
        </div>
      </div>
    </div>
  )
}
