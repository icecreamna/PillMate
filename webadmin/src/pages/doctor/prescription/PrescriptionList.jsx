import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/doctor/prescription/PrescriptionList.module.css'

/** mock แทน API จริง */
const MOCK_PATIENTS = [
  { id: 1, name: 'สมชาย ใจดี', idcard: '1234567890100', gender: 'ชาย', age: 30 },
  { id: 2, name: 'สมหญิง ใจร้าย', idcard: '1234567890101', gender: 'หญิง', age: 40 },
  { id: 3, name: 'สมหมาย ใจบุญ', idcard: '1234567890102', gender: 'ชาย', age: 55 },
]

export default function PrescriptionList() {
  const nav = useNavigate()

  const [q, setQ] = useState('')
  // เริ่มต้นให้เห็น "ผู้ป่วยทั้งหมด"
  const [results, setResults] = useState(MOCK_PATIENTS)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const onSearch = async () => {
    setError('')
    const qtrim = q.trim()

    try {
      setLoading(true)
      // TODO: เรียก API จริงหากต้องการ
      await new Promise(r => setTimeout(r, 200)) // mock latency เล็กน้อย

      if (!qtrim) {
        // ถ้าไม่กรอก -> แสดงทั้งหมด
        setResults(MOCK_PATIENTS)
        return
      }

      // ค้นหาแบบ "ตรงตัว" ตามเลขบัตรประชาชน → ให้ขึ้นมา 1 รายการ (หรือว่างถ้าไม่พบ)
      const data = MOCK_PATIENTS.filter(p => p.idcard === qtrim)
      if (data.length === 0) setError('ไม่พบผู้ป่วย')
      setResults(data)
    } catch (err) {
      setError('ค้นหาไม่สำเร็จ')
    } finally {
      setLoading(false)
    }
  }

  const onKeyDown = (e) => {
    if (e.key === 'Enter') { e.preventDefault(); onSearch() }
  }

  return (
    <div>
      <h2 className={styles.title}>Prescription</h2>

      {/* แถวค้นหา */}
      <div className={styles.searchRow}>
        <div className={styles.searchWrap}>
          <input
            className={styles.searchInput}
            placeholder="ค้นหาเลขบัตรประชาชน"
            value={q}
            onChange={e => setQ(e.target.value)}
            onKeyDown={onKeyDown}
          />
          <span className={styles.searchIcon}>🔍</span>
        </div>
        <button className={styles.searchBtn} onClick={onSearch} disabled={loading}>
          {loading ? 'กำลังค้นหา...' : 'ค้นหา'}
        </button>
      </div>

      {error && <div className={styles.error}>{error}</div>}

      {/* ตารางผลลัพธ์ */}
      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th style={{width:'6%'}}>#</th>
              <th style={{width:'32%'}}>Name</th>
              <th style={{width:'26%'}}>IDCardNumber</th>
              <th style={{width:'12%'}}>Gender</th>
              <th style={{width:'12%'}}>Age</th>
              <th style={{width:'15%'}}># Action</th>
            </tr>
          </thead>
          <tbody>
            {results.length === 0 ? (
              <tr>
                <td colSpan={6} style={{textAlign:'center', color:'#6b7280', height:56}}>
                  ไม่มีข้อมูล
                </td>
              </tr>
            ) : (
              results.map((r, i) => (
                <tr key={r.id}>
                  <td>{i + 1}</td>
                  <td>{r.name}</td>
                  <td>{r.idcard}</td>
                  <td>{r.gender}</td>
                  <td>{r.age}</td>
                  <td className={styles.actions}>
                    <button className={styles.viewBtn} onClick={()=>nav(`/doc/prescription/view/${r.id}`)}>view</button>
                    <button className={styles.addBtn} onClick={()=>nav(`/doc/prescription/add/${r.id}`)}>Add</button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
