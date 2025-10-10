import { useMemo, useState } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/prescription/AddPrescription.module.css'

// mock ยาทั้งหมด (แทนการเรียก API)
const ALL_DRUGS = [
  { id: 1, medName: 'Medicine A', generic: 'Generic A', strength: '500 mg', form: 'ยาเม็ด' },
  { id: 2, medName: 'Medicine B', generic: 'Generic B', strength: '400 mg', form: 'แคปซูล' },
  { id: 3, medName: 'Medicine C', generic: 'Generic C', strength: '4 mg/5 ml', form: 'ยาน้ำ' },
]

export default function AddPrescription() {
  const nav = useNavigate()
  const { patientId } = useParams()

  // เริ่มต้นยังไม่เลือก
  const initial = useMemo(
    () => ALL_DRUGS.map(d => ({ ...d, checked: false })),
    []
  )
  const [rows, setRows] = useState(initial)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  // --- Modal state ---
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState(null) // {id, medName, generic, ...}
  const [dosePerTime, setDosePerTime] = useState('1') // ครั้งละ
  const [timesPerDay, setTimesPerDay] = useState('3') // วันละ

  const toggle = (id) => {
    setRows(rs => rs.map(r => (r.id === id ? { ...r, checked: !r.checked } : r)))
  }

  const specify = (row) => {
    // เปิด modal โดยโหลดค่าที่เคยใส่ไว้ (ถ้ามี)
    setEditing(row)
    setDosePerTime(String(row.dosePerTime ?? '1'))
    setTimesPerDay(String(row.timesPerDay ?? '3'))
    setOpen(true)
  }

  const closeModal = () => {
    setOpen(false)
    setEditing(null)
  }

  const confirmSpecify = () => {
    // อัปเดตค่าที่ใส่กับยาตัวนั้นในตาราง
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
    try {
      setSaving(true)
      // TODO: POST /api/prescriptions  (patientId, items: selected)
      await new Promise(r => setTimeout(r, 600))
      alert(
        `(mock) สร้างใบสั่งยาให้ patientId=${patientId}\n` +
        selected.map(s => {
          const d = s.dosePerTime ? ` | ครั้งละ ${s.dosePerTime} | วันละ ${s.timesPerDay}` : ''
          return `- ${s.medName}${d}`
        }).join('\n')
      )
      nav('/doc/prescription', { replace: true })
    } catch (e) {
      setError('บันทึกไม่สำเร็จ')
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
                  {r.dosePerTime ? (
                    <div className={styles.doseNote}>
                      ครั้งละ {r.dosePerTime} | วันละ {r.timesPerDay}
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
          </tbody>
        </table>
      </div>

      {error && <div className={styles.error}>{error}</div>}

      <div className={styles.footer}>
        <button className={styles.submit} onClick={submit} disabled={saving}>
          {saving ? 'กำลังบันทึก…' : 'จ่ายยา'}
        </button>
      </div>

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
                <div className={styles.unit}>เม็ด</div>
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
