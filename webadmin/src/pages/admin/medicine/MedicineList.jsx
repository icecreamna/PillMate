import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/admin/medicine/MedicineList.module.css'
import { listMedicines, deleteMedicine } from '../../../services/medicines'
import { listForms } from '../../../services/initialData'  // ⬅️ โหลดฟอร์มมา map ชื่อ

export default function MedicineList() {
  const navigate = useNavigate()
  const [rows, setRows] = useState([])
  const [formMap, setFormMap] = useState({})   // { [form_id]: form_name }
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  useEffect(() => {
    const run = async () => {
      setLoading(true); setError("")
      try {
        // โหลดพร้อมกัน: รายการยา + รายการฟอร์ม
        const [resMeds, forms] = await Promise.all([
          listMedicines(),           // GET /admin/medicine-infos -> { data: [...] }
          listForms(),               // GET /forms -> array of { id, form_name }
        ])

        const list = Array.isArray(resMeds?.data) ? resMeds.data : []
        setRows(list)

        const map = {}
        ;(Array.isArray(forms) ? forms : []).forEach(f => {
          map[f.id] = f.form_name || f.name || `Form #${f.id}`
        })
        setFormMap(map)
      } catch (e) {
        setError(e.message || "โหลดข้อมูลไม่สำเร็จ")
      } finally {
        setLoading(false)
      }
    }
    run()
  }, [])

  const onDelete = async (id) => {
    if (!confirm('ยืนยันการลบรายการนี้?')) return
    try {
      await deleteMedicine(id)                // DELETE /admin/medicine-info/:id
      setRows(prev => prev.filter(x => x.id !== id))
    } catch (e) {
      alert(e.message || "ลบไม่สำเร็จ")
    }
  }

  // ดึงชื่อฟอร์ม: ถ้า DTO มี form_name ใช้อันนั้นก่อน, ไม่มีก็ map จาก form_id
  const renderForm = (m) => {
    if (m?.form_name) return m.form_name
    if (m?.form) return m.form
    if (m?.form_id && formMap[m.form_id]) return formMap[m.form_id]
    return "-"
  }

  const pick = (o, keys, fallback = "-") => {
    for (const k of keys) {
      const v = o?.[k]
      if (v != null && v !== "") return v
    }
    return fallback
  }

  return (
    <div>
      <div className={styles.headerRow}>
        <h2 className={styles.title}>Medicine Information</h2>
        <button className={styles.addBtn} onClick={()=>navigate('/admin/medicine-info/add')}>
          + Add Medicine
        </button>
      </div>

      {error && <div className={styles.error}>{error}</div>}
      {loading ? (
        <div className={styles.loading}>กำลังโหลด...</div>
      ) : (
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
              {rows.map((m, i)=>(
                <tr key={m.id ?? i}>
                  <td>{i+1}</td>
                  <td>{pick(m, ["med_name","medName"])}</td>
                  <td>{pick(m, ["generic_name","genericName"])}</td>
                  <td>{pick(m, ["strength"])}</td>
                  <td>{renderForm(m)}</td>
                  <td className={styles.actions}>
                    <button className={styles.view}  onClick={()=>navigate(`/admin/medicine-info/${m.id}`)}>View</button>
                    <button className={styles.edit}  onClick={()=>navigate(`/admin/medicine-info/${m.id}/edit`)}>Edit</button>
                    <button className={styles.delete} onClick={()=>onDelete(m.id)}>Delete</button>
                  </td>
                </tr>
              ))}
              {rows.length === 0 && (
                <tr><td colSpan={6} style={{textAlign:'center', padding:12}}>ไม่พบข้อมูล</td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
