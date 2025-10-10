import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/admin/medicine/MedicineList.module.css'

const mock = [
  { id: 1, medName: 'Medicine A', genericName: 'Generic A', strength: '500 mg', form: 'ยาเม็ด' },
  { id: 2, medName: 'Medicine B', genericName: 'Generic B', strength: '400 mg', form: 'แคปซูล' },
  { id: 3, medName: 'Medicine C', genericName: 'Generic C', strength: '4 mg/5 ml', form: 'ยาน้ำ' },
]

export default function MedicineList() {
  const navigate = useNavigate()

  const onDelete = (id) => {
    if (!confirm('ยืนยันการลบรายการนี้?')) return
    // TODO: เรียก API ลบจริง แล้วรีเฟรชข้อมูล
    alert(`(mock) ลบ med id = ${id}`)
  }

  return (
    <div>
      <div className={styles.headerRow}>
        <h2 className={styles.title}>Medicine Information</h2>
        <button className={styles.addBtn} onClick={()=>navigate('/admin/medicine-info/add')}>
          + Add Medicine
        </button>
      </div>

      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th style={{width:'5%'}}>#</th>
              <th style={{width:'20%'}}>MedName</th>
              <th style={{width:'20%'}}>GenericName</th>
              <th style={{width:'20%'}}>Strength</th>
              <th style={{width:'15%'}}>Form</th>
              <th style={{width:'20%'}}># Action</th>
            </tr>
          </thead>
          <tbody>
            {mock.map((m, i)=>(
              <tr key={m.id}>
                <td>{i+1}</td>
                <td>{m.medName}</td>
                <td>{m.genericName}</td>
                <td>{m.strength}</td>
                <td>{m.form}</td>
                <td className={styles.actions}>
                  <button className={styles.view} onClick={()=>navigate(`/admin/medicine-info/${m.id}`)}>View</button>
                  <button className={styles.edit} onClick={()=>navigate(`/admin/medicine-info/${m.id}/edit`)}>Edit</button>
                  <button className={styles.delete} onClick={()=>onDelete(m.id)}>Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
