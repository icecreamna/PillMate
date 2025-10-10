import { useMemo, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/doctor/patient/PatientList.module.css'

const mock = [
  { id:1, name:'สมชาย ใจดี', idcard:'1234567890100', gender:'ชาย', age:30 },
  { id:2, name:'สมหญิง ใจร้าย', idcard:'1234567890101', gender:'หญิง', age:40 },
  { id:3, name:'สมหมาย ใจบุญ', idcard:'1234567890102', gender:'ชาย', age:55 },
]

export default function PatientList(){
  const nav = useNavigate()

  // state สำหรับช่องค้นหา และผลลัพธ์ที่เลือก
  const [query, setQuery] = useState('')
  const [selected, setSelected] = useState(null) // เก็บผู้ป่วยที่ค้นหาเจอ (หรือ null = แสดงทั้งหมด)

  const onDelete = (id)=>{ if(confirm('ยืนยันลบ?')) alert('(mock) ลบ id='+id) }

  // เมื่อกดค้นหา: ถ้าเว้นว่าง => กลับมาแสดงทั้งหมด, ถ้ามีค่า => ค้นหาตรงกับ idcard
  const handleSearch = () => {
    const q = query.trim()
    if (!q) {
      setSelected(null) // แสดงทั้งหมด
      return
    }
    const found = mock.find(p => p.idcard === q)
    if (found) {
      setSelected(found) // แสดงเฉพาะคนนี้
    } else {
      alert('ไม่พบหมายเลขบัตรประชาชนนี้')
      // ไม่เปลี่ยน selected เพื่อคงรายการเดิม (หรือจะ setSelected(null) เพื่อโชว์ทั้งหมดก็ได้)
    }
  }

  // รองรับกด Enter ในช่องค้นหา
  const onKeyDown = (e) => {
    if (e.key === 'Enter') handleSearch()
  }

  // รายการที่จะแสดง: ถ้ามี selected ให้แสดงเฉพาะคนนั้น, ไม่งั้นแสดงทั้งหมด
  const rows = useMemo(() => {
    return selected ? [selected] : mock
  }, [selected])

  return (
    <div>
      <div className={styles.headerRow}>
        <h2 className={styles.title}>Patients</h2>
        <div className={styles.headerActions}>
          <input
            className={styles.search}
            placeholder="ค้นหาหมายเลขบัตรประชาชน"
            value={query}
            onChange={e=>setQuery(e.target.value)}
            onKeyDown={onKeyDown}
          />
          <button className={styles.searchBtn} onClick={handleSearch}>ค้นหา</button>
          <button className={styles.addBtn} onClick={()=>nav('/doc/patients/add')}>+ Add Patient</button>
        </div>
      </div>

      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th style={{width:'6%'}}>#</th>
              <th style={{width:'20%'}}>Name</th>
              <th style={{width:'20%'}}>IDCardNumber</th>
              <th style={{width:'12%'}}>Gender</th>
              <th style={{width:'10%'}}>Age</th>
              <th style={{width:'25%'}}># Action</th>
            </tr>
          </thead>
          <tbody>
            {rows.map((p,i)=>(
              <tr key={p.id}>
                <td>{i+1}</td>
                <td>{p.name}</td>
                <td>{p.idcard}</td>
                <td>{p.gender}</td>
                <td>{p.age}</td>
                <td className={styles.actions}>
                  <button className={styles.view} onClick={()=>nav(`/doc/patients/${p.id}`)}>View</button>
                  <button className={styles.edit} onClick={()=>nav(`/doc/patients/${p.id}/edit`)}>Edit</button>
                  <button className={styles.delete} onClick={()=>onDelete(p.id)}>Delete</button>
                </td>
              </tr>
            ))}
            {rows.length === 0 && (
              <tr><td colSpan={6} style={{textAlign:'center', opacity:.7, padding:'12px'}}>ไม่พบข้อมูล</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  )
}
