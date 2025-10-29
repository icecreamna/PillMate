// src/pages/doctor/prescription/AddPrescription.jsx
import { useEffect, useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/prescription/AddPrescription.module.css'
import { getPatient } from '../../../services/patients'
import { listMedicinesForDoctor } from '../../../services/medicines' // GET /doctor/medicine-infos
import { listForms } from '../../../services/initialData'            // GET /forms?with_relations=true
import { createPrescription } from '../../../services/prescriptions'

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

// ช่วย normalize ข้อความให้ค้นหาง่ายขึ้น (ไม่สนช่องว่าง/ตัวพิมพ์ใหญ่เล็ก/สัญลักษณ์)
function norm(s) {
  return String(s || '')
    .toLocaleLowerCase('th-TH')
    .normalize('NFKD')
    .replace(/[\u0300-\u036f]/g, '') // strip combining marks
    .replace(/\s+/g, '')             // remove spaces
    .replace(/[^\p{L}\p{N}.%-]/gu, '') // keep letters/numbers/บางสัญลักษณ์พื้นฐาน
}

export default function AddPrescription() {
  const nav = useNavigate()
  const { patientId } = useParams()

  const [loading, setLoading] = useState(true)
  const [saving, setSaving]   = useState(false)
  const [error, setError]     = useState('')

  const [patient, setPatient] = useState(null)
  const [rows, setRows] = useState([])

  // --- Search state ---
  const [q, setQ] = useState('')

  // --- Modal state ---
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState(null)
  const [dosePerTime, setDosePerTime] = useState('1')
  const [timesPerDay, setTimesPerDay] = useState('3')

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      setLoading(true); setError('')
      try {
        // ผู้ป่วย
        const pRes = await getPatient(patientId)
        const p = pRes?.data
        if (!p) throw new Error('ไม่พบข้อมูลผู้ป่วย')
        if (cancelled) return
        setPatient(p)

        // forms + units
        const fs = await listForms({ with_relations: true })
        const forms = Array.isArray(fs) ? fs : []
        const formMap = new Map(forms.map(f => [String(f.id), f]))
        const unitLookupByForm = new Map(
          forms.map(f => [
            String(f.id),
            new Map((Array.isArray(f.units) ? f.units : []).map(u => [String(u.id), u.unit_name || u.name || '']))
          ])
        )

        // ยาที่หมอเห็นได้
        const mRes = await listMedicinesForDoctor()
        const meds = Array.isArray(mRes?.data) ? mRes.data : []

        // map ยา -> form + unit
        const mapped = meds.map(m => {
          const fid = String(m.form_id ?? '')
          const form = formMap.get(fid) || null
          const formName =
            m.form_name || m.form ||
            form?.form_name || form?.name || '-'

          const unitFromMed = String(m.unit_name || m.unit || '').trim()
          const unitFromLookup = unitLookupByForm.get(fid)?.get(String(m.unit_id ?? '')) || ''
          const unitLabel = unitFromMed || unitFromLookup || inferUnitFromForm(formName) || ''

          return {
            id: m.id,
            medName: m.med_name || m.MedName || '-',
            generic: m.generic_name || m.GenericName || '-',
            strength: m.strength || '-',
            form: formName,
            unitLabel,
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

  // สำเนาเริ่มต้น (กันกรณีอยากเทียบ diff ภายหลัง)
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

  // กรองรายการตามคำค้นหา (ค้นหาจากหลายช่อง)
  const visibleRows = useMemo(() => {
    const nq = norm(q)
    if (!nq) return rows
    return rows.filter(r => {
      return (
        norm(r.medName).includes(nq) ||
        norm(r.generic).includes(nq) ||
        norm(r.form).includes(nq) ||
        norm(r.strength).includes(nq)
      )
    })
  }, [rows, q])

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

    const toNumStr = (v) => {
      const s = String(v ?? '').trim()
      return s === '' ? '0' : s
    }

    const items = selected.map(s => ({
      medicine_info_id: s.id,
      amount_per_time: toNumStr(s.dosePerTime),
      times_per_day:   toNumStr(s.timesPerDay),
    }))

    try {
      setSaving(true)
      await createPrescription({ id_card_number: idCard, items })
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

      {/* Patient bar */}
      {patient && (
        <div className={styles.patientBar}>
          <div><strong>ชื่อผู้ป่วย:</strong> {[patient.first_name, patient.last_name].filter(Boolean).join(' ') || '-'}</div>
          <div><strong>Patient Code:</strong> {patient.patient_code || '-'}</div>
          <div><strong>เลขบัตร:</strong> {patient.id_card_number || '-'}</div>
        </div>
      )}

      {/* Search box */}
      <div className={styles.toolsBar}>
        <input
          type="text"
          className={styles.searchInput}
          value={q}
          onChange={e => setQ(e.target.value)}
          placeholder="ค้นหาชื่อยา / Generic / รูปแบบยา (ไทย/English)"
          aria-label="ค้นหาชื่อยา"
        />
        <div className={styles.searchInfo}>
          แสดง {visibleRows.length} จากทั้งหมด {rows.length} รายการ
        </div>
      </div>

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
                {visibleRows.map((r) => (
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
                          ครั้งละ {r.dosePerTime} {r.unitLabel} | วันละ {r.timesPerDay} ครั้ง
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
                {visibleRows.length === 0 && (
                  <tr>
                    <td colSpan={6} style={{textAlign:'center', color:'#6b7280', height:56}}>
                      ไม่พบรายการที่ตรงกับคำค้น
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
                <div className={styles.unit}>{editing.unitLabel}</div>
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
