// src/pages/doctor/prescription/AddPrescription.jsx
import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/prescription/AddPrescription.module.css'
import { getPatient } from '../../../services/patients'
import { listMedicinesForDoctor } from '../../../services/medicines' // GET /doctor/medicine-infos
import { listForms } from '../../../services/initialData'            // GET /forms?with_relations=true
import { createPrescription } from '../../../services/prescriptions'

/** (เลิกใช้) เดาหน่วยจากฟอร์ม — คงไว้เพื่ออ้างอิง แต่จะไม่ถูกนำมาใช้แล้ว */
function inferUnitFromForm(formName) {
  const s = String(formName || '').trim()
  if (!s) return ''
  if (s.includes('เม็ด')) return 'เม็ด'
  if (s.includes('แคปซูล')) return 'แคปซูล'
  if (s.includes('ยาน้ำ')) return 'ml'
  if (s.includes('ขี้ผึ้ง')) return 'กรัม'
  if (s.includes('สเปรย์')) return 'สเปรย์'
  return ''
}

export default function AddPrescription() {
  const nav = useNavigate()
  const { patientId } = useParams()

  const [loading, setLoading] = useState(true)
  const [saving, setSaving]   = useState(false)
  const [error, setError]     = useState('')

  const [patient, setPatient] = useState(null)
  const [rows, setRows] = useState([])

  // --- Modal state ---
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState(null)
  const [dosePerTime, setDosePerTime] = useState('1') // ครั้งละ (ตัวเลข)
  const [timesPerDay, setTimesPerDay] = useState('3') // วันละ (ตัวเลข)

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoading(true); setError('')
      try {
        // 1) ผู้ป่วย -> ใช้ id_card_number ตอนยิง create
        const pRes = await getPatient(patientId)
        const p = pRes?.data
        if (!p) throw new Error('ไม่พบข้อมูลผู้ป่วย')
        if (cancelled) return
        setPatient(p)

        // 2) โหลด forms (พร้อม relations) เพื่อทำ lookup form_name (ไม่ใช้เพื่อหาหน่วย)
        //    คาดหวังโครง: [{id, form_name, units:[{id, unit_name}, ...]}]
        const fRes = await listForms({ with_relations: true })
        const forms = Array.isArray(fRes) ? fRes : []
        const formMap = new Map(forms.map(f => [String(f.id), f]))

        // 3) โหลดยาที่หมอเห็นได้
        const mRes = await listMedicinesForDoctor() // -> { data:[...] }
        const meds = Array.isArray(mRes?.data) ? mRes.data : []

        // 4) map ยา -> เติม form_name จาก form_id/ยา
        //    หน่วย: ใช้เฉพาะ "ที่ติดมากับยาตัวนั้น" เท่านั้น ถ้าไม่มีให้เป็น '' (แล้วไปแสดง '-' ที่ UI)
        const mapped = meds.map(m => {
          const fm = formMap.get(String(m.form_id || '')) || null
          const formName =
            m.form_name || m.form ||
            fm?.form_name || fm?.name || '-'

          // หน่วย: ไม่ดึงจาก form.units / ไม่เดาจากชื่อฟอร์มอีกต่อไป
          const unitLabel = String(m.unit_name || m.unit || '').trim()

          return {
            id: m.id,
            medName: m.med_name || m.MedName || '-',
            generic: m.generic_name || m.GenericName || '-',
            strength: m.strength || '-',
            form: formName,
            unitLabel, // ถ้าไม่มี = ''
            checked: false,
            dosePerTime: undefined,
            timesPerDay: undefined,
          }
        })

        if (!cancelled) setRows(mapped)
      } catch (e) {
        if (!cancelled) setError(e.message || 'โหลดข้อมูลไม่สำเร็จ')
      } finally {
        if (!cancelled) setLoading(false)
      }
    })()
    return () => { cancelled = true }
  }, [patientId])

  const initial = useMemo(() => rows.map(d => ({ ...d })), [rows])

  const toggle = (id) => {
    setRows(rs => rs.map(r => (r.id === id ? { ...r, checked: !r.checked } : r)))
  }

  const specify = (row) => {
    setEditing(row)
    setDosePerTime(String(row.dosePerTime ?? '1'))
    setTimesPerDay(String(row.timesPerDay ?? '3'))
    setOpen(true)
  }
  const closeModal = () => { setOpen(false); setEditing(null) }
  const confirmSpecify = () => {
    setRows(rs =>
      rs.map(r =>
        r.id === editing.id
          ? { ...r, dosePerTime: Number(dosePerTime || 0), timesPerDay: Number(timesPerDay || 0) }
          : r
      )
    )
    closeModal()
  }

  const submit = async () => {
    setError('')
    const selected = rows.filter(r => r.checked)
    if (selected.length === 0) {
      setError('กรุณาเลือกรายการยาอย่างน้อย 1 รายการ')
      return
    }
    const idCard = patient?.id_card_number
    if (!idCard) {
      setError('ไม่พบเลขบัตรประชาชนของผู้ป่วย')
      return
    }

    // helper: ส่งเป็น "เลขล้วน (สตริง)" ถ้าไม่มีให้เป็น "0"
    const toNumStr = (v) => {
      const s = String(v ?? '').trim()
      return s === '' ? '0' : s
    }

    // items -> ส่งเฉพาะ "เลขล้วน" ไม่พ่วงหน่วย/คำว่า "ครั้ง"
    const items = selected.map(s => ({
      medicine_info_id: s.id,
      amount_per_time: toNumStr(s.dosePerTime), // e.g. "1"
      times_per_day:   toNumStr(s.timesPerDay),  // e.g. "3"
    }))

    try {
      setSaving(true)
      await createPrescription({
        id_card_number: idCard,
        items,
        // doctor_id / sync_until / app_sync_status ไม่ต้องส่ง ถ้า backend เติมเองจาก token/Default
      })
      nav('/doc/prescription', { replace: true })
    } catch (e) {
      setError(e.message || 'บันทึกไม่สำเร็จ')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div>
      <div className={styles.header}>
        <h2 className={styles.title}>Add Prescription</h2>
        <button className={styles.back} onClick={() => nav(-1)}>← Back</button>
      </div>

      {patient && (
        <div className={styles.patientBar}>
          <div><strong>ชื่อผู้ป่วย:</strong> {[patient.first_name, patient.last_name].filter(Boolean).join(' ') || '-'}</div>
          <div><strong>เลขบัตร:</strong> {patient.id_card_number || '-'}</div>
        </div>
      )}

      {error && <div className={styles.error}>{error}</div>}

      {loading ? (
        <div className={styles.loading}>กำลังโหลด...</div>
      ) : (
        <>
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
                {rows.map((r) => (
                  <tr key={r.id} className={r.checked ? '' : styles.dim}>
                    <td>
                      <input
                        type="checkbox"
                        checked={r.checked}
                        onChange={() => toggle(r.id)}
                        className={styles.checkbox}
                        aria-label={`เลือก ${r.medName}`}
                      />
                    </td>
                    <td>
                      {r.medName}
                      {(r.dosePerTime != null && r.timesPerDay != null) ? (
                        <div className={styles.doseNote}>
                          {/* แสดงผล: ถ้าไม่มีหน่วยให้ขึ้น '-' */}
                          ครั้งละ {r.dosePerTime} {r.unitLabel?.trim() ? r.unitLabel : '-'} | วันละ {r.timesPerDay} ครั้ง
                        </div>
                      ) : null}
                    </td>
                    <td>{r.generic}</td>
                    <td>{r.strength}</td>
                    <td>{r.form}</td>
                    <td className={styles.actions}>
                      <button className={styles.specify} onClick={() => specify(r)}>specify</button>
                    </td>
                  </tr>
                ))}
                {rows.length === 0 && (
                  <tr>
                    <td colSpan={6} style={{textAlign:'center', color:'#6b7280', height:56}}>
                      ไม่มียาให้เลือก
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          <div className={styles.footer}>
            <button className={styles.submit} onClick={submit} disabled={saving || rows.length === 0}>
              {saving ? 'กำลังบันทึก…' : 'จ่ายยา'}
            </button>
          </div>
        </>
      )}

      {/* Modal */}
      {open && editing && (
        <div className={styles.modalOverlay} onClick={closeModal}>
          <div className={styles.modal} onClick={e => e.stopPropagation()}>
            <button className={styles.modalClose} onClick={closeModal} aria-label="Close">×</button>

            <div className={styles.modalHead}>
              <div>MedName : <strong>{editing.medName}</strong></div>
              <div>GenericName : <strong>{editing.generic}</strong></div>
            </div>

            <div className={styles.box}>
              <div className={styles.boxTitle}>การใช้ยา</div>

              <div className={styles.row}>
                <label className={styles.label}>
                  <span>ครั้งละ</span>
                  <input
                    type="number"
                    min="0"
                    className={styles.input}
                    value={dosePerTime}
                    onChange={e => setDosePerTime(e.target.value)}
                  />
                </label>
                <div className={styles.unit}>
                  {editing.unitLabel && editing.unitLabel.trim() ? editing.unitLabel : '-'}
                </div>
              </div>

              <div className={styles.row}>
                <label className={styles.label}>
                  <span>วันละ</span>
                  <input
                    type="number"
                    min="0"
                    className={styles.input}
                    value={timesPerDay}
                    onChange={e => setTimesPerDay(e.target.value)}
                  />
                </label>
                <div className={styles.unit}>ครั้ง</div>
              </div>
            </div>

            <div className={styles.modalActions}>
              <button className={styles.confirm} onClick={confirmSpecify}>ยืนยัน</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
