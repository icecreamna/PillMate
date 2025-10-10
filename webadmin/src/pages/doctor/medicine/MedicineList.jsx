import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/doctor/medicine/MedicineList.module.css'

const rows = [
  { id:1, medName:'Medicine A', generic:'Generic A', strength:'500 mg', form:'ยาเม็ด' },
  { id:2, medName:'Medicine B', generic:'Generic B', strength:'400 mg', form:'แคปซูล' },
  { id:3, medName:'Medicine C', generic:'Generic C', strength:'4 mg/5 ml', form:'ยาน้ำ' },
]

export default function MedicineList(){
  const nav = useNavigate()
  return (
    <div>
      <h2 className={styles.title}>Medicine Information</h2>
      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th style={{width:'6%'}}>#</th>
              <th style={{width:'28%'}}>MedName</th>
              <th style={{width:'28%'}}>GenericName</th>
              <th style={{width:'18%'}}>Strength</th>
              <th style={{width:'12%'}}>Form</th>
              <th style={{width:'8%'}}># Action</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((r,i)=>(
              <tr key={r.id}>
                <td>{i+1}</td>
                <td>{r.medName}</td>
                <td>{r.generic}</td>
                <td>{r.strength}</td>
                <td>{r.form}</td>
                <td className={styles.actions}>
                  <button className={styles.view} onClick={()=>nav(`/doc/medicine-info/${r.id}`)}>
                    View
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
