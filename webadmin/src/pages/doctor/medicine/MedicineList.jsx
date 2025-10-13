import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/doctor/medicine/MedicineList.module.css'
import { listMedicinesForDoctor } from '../../../services/medicines'

export default function MedicineList(){
  const nav = useNavigate()
  const [rows, setRows] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoading(true); setError('')
      try {
        const res = await listMedicinesForDoctor()        // ← ใช้เส้นทาง /doctor
        const list = Array.isArray(res?.data) ? res.data : []
        if (!cancelled) setRows(list)
      } catch (e) {
        if (!cancelled) setError(e.message || 'โหลดข้อมูลยาไม่สำเร็จ')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => { cancelled = true }
  }, [])

  const pick = (o, keys, fallback = '-') => {
    for (const k of keys) { const v = o?.[k]; if (v != null && v !== '') return v }
    return fallback
  }

  return (
    <div>
      <h2 className={styles.title}>Medicine Information</h2>

      {error && <div className={styles.error}>{error}</div>}
      {loading ? (
        <div className={styles.loading}>กำลังโหลด...</div>
      ) : (
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
                <tr key={r.id ?? i}>
                  <td>{i+1}</td>
                  <td>{pick(r,['med_name','medName'])}</td>
                  <td>{pick(r,['generic_name','generic'])}</td>
                  <td>{pick(r,['strength'])}</td>
                  <td>{pick(r,['form_name','form'])}</td>
                  <td className={styles.actions}>
                    <button className={styles.view} onClick={()=>nav(`/doc/medicine-info/${r.id}`)}>
                      View
                    </button>
                  </td>
                </tr>
              ))}
              {rows.length === 0 && (
                <tr><td colSpan={6} style={{textAlign:'center', padding:'12px', opacity:.7}}>ไม่พบข้อมูล</td></tr>
              )}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
