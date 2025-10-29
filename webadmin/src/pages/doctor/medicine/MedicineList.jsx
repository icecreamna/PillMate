import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import styles from '../../../styles/doctor/medicine/MedicineList.module.css'
import { listMedicinesForDoctor } from '../../../services/medicines'
import { listForms } from '../../../services/initialData'

// --- Helpers: normalize various response shapes ---
const asArray = (x) => {
  if (Array.isArray(x)) return x
  if (Array.isArray(x?.data)) return x.data
  if (Array.isArray(x?.items)) return x.items
  if (Array.isArray(x?.list)) return x.list
  if (Array.isArray(x?.rows)) return x.rows
  if (Array.isArray(x?.result)) return x.result
  if (Array.isArray(x?.forms)) return x.forms      // เผื่อ /forms คืน {forms:[...]}
  return []
}

const pick = (o, keys, fallback = '-') => {
  for (const k of keys) {
    const v = o?.[k]
    if (v != null && String(v).trim() !== '') return String(v).trim()
  }
  return fallback
}

export default function MedicineList() {
  const nav = useNavigate()
  const [rows, setRows] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoading(true); setError('')
      try {
        // ดึงยา + ฟอร์มพร้อมกัน
        const [medRes, formRes] = await Promise.all([
          listMedicinesForDoctor(),   // ควรคืน array หรือ {data:[]}
          listForms({ with_relations: false }), // ควรคืน array หรือ {data:[]}/{forms:[]}
        ])

        const meds = asArray(medRes)
        const forms = asArray(formRes)

        if (!meds.length) {
          if (!cancelled) {
            setRows([])
            setError('ไม่พบข้อมูลยา (GET /doctor/medicine-infos)')
          }
          return
        }

        // map: form_id -> form_name (รองรับหลายคีย์ของชื่อฟอร์ม)
        const formNameById = new Map()
        for (const f of forms) {
          const id = f?.id ?? f?.form_id
          if (!id) continue
          const name = f?.form_name ?? f?.name ?? f?.formName ?? f?.FormName ?? f?.title ?? f?.label
          if (name) formNameById.set(Number(id), String(name))
        }

        // enrich ชื่อฟอร์มให้แต่ละแถว
        const enriched = meds.map(m => {
          const fid = Number(m?.form_id || 0)
          const formName = m?.form_name || (fid ? formNameById.get(fid) : '')
          return { ...m, form_name: formName || '-' }
        })

        if (!cancelled) setRows(enriched)
      } catch (e) {
        if (!cancelled) setError(e?.message || 'โหลดข้อมูลไม่สำเร็จ')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => { cancelled = true }
  }, [])

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
                  <td>{pick(r,['strength','strength_text','strengthText'])}</td>
                  <td title={r.form_id ? `form_id=${r.form_id}` : ''}>
                    {pick(r,['form_name','form'])}
                  </td>
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
