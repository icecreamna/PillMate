// src/pages/admin/DoctorList.jsx
import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/admin/doctor/DoctorList.module.css'

const mock = [
  { id: 1, name: 'Doctor A', username: 'doctor_a@pillmate.com', password: '**********' },
  { id: 2, name: 'Doctor B', username: 'doctor_b@pillmate.com', password: '**********' },
  { id: 3, name: 'Doctor C', username: 'doctor_c@pillmate.com', password: '**********' },
]

export default function DoctorList() {
  const navigate = useNavigate()

  const onDelete = (id) => {
    if (!confirm('ยืนยันการลบรายการนี้?')) return
    // TODO: เรียก API ลบจริง แล้วรีเฟรชข้อมูล
    alert(`(mock) ลบ doctor id = ${id}`)
  }

  return (
    <div>
      {/* หัวข้อ + ปุ่มอยู่บรรทัดเดียวกัน */}
      <div className={styles.headerRow}>
        <h2 className={styles.title}>Doctors</h2>
        <button
          className={styles.addBtn}
          onClick={() => navigate('/admin/add')}
        >
          + Add Doctor
        </button>
      </div>

      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th style={{ width: '5%' }}>#</th>
              <th style={{ width: '25%' }}>Name</th>
              <th style={{ width: '40%' }}>Username</th>
              <th style={{ width: '20%' }}>Password</th>
              <th style={{ width: '10%' }}># Action</th>
            </tr>
          </thead>
          <tbody>
            {mock.map((d) => (
              <tr key={d.id}>
                <td>{d.id}</td>
                <td>{d.name}</td>
                <td>{d.username}</td>
                <td>{d.password}</td>
                <td className={styles.actions}>
                  <button
                    className={styles.edit}
                    onClick={() => navigate(`/admin/${d.id}/edit`)}
                  >
                    Edit
                  </button>
                  <button
                    className={styles.delete}
                    onClick={() => onDelete(d.id)}
                  >
                    Delete
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
