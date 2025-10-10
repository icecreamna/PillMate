import { useMemo } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/appointment/ViewAppointment.module.css'

/** === MOCK DATA (แทน API จริง) === */
const MOCK_APPOINTMENTS = [
  // patient 1
  { id:'ap-1005', patientId:1, patientName:'สมชาย ใจดี', at:'2025-07-09T08:00:00+07:00', note:'งดอาหาร 8 ชั่วโมง ก่อนเจาะเลือด' },
  { id:'ap-1004', patientId:1, patientName:'สมชาย ใจดี', at:'2025-06-25T10:30:00+07:00', note:'ติดตามผลตรวจ' },
  { id:'ap-1001', patientId:1, patientName:'สมชาย ใจ ดี', at:'2025-05-18T09:00:00+07:00', note:'ตรวจสุขภาพประจำปี' },

  // patient 2
  { id:'ap-2001', patientId:2, patientName:'สมหญิง ใจร้าย', at:'2025-06-18T14:00:00+07:00', note:'อัลตราซาวด์ช่องท้อง' },
];

function formatDateTH(d) {
  const dt = new Date(d);
  return dt.toLocaleDateString('th-TH', { day:'numeric', month:'short', year:'numeric' });
}
function formatTimeTH(d) {
  const dt = new Date(d);
  return dt.toLocaleTimeString('th-TH', { hour:'2-digit', minute:'2-digit' });
}

export default function ViewAppointment() {
  const nav = useNavigate();
  const { id: patientId } = useParams(); // มาจาก /doc/appointment/view/:id

  // TODO: เปลี่ยนไปเรียก API จริงด้วย patientId
  const rows = useMemo(() => {
    const list = MOCK_APPOINTMENTS.filter(a => String(a.patientId) === String(patientId));
    // เรียงล่าสุดก่อน
    return list.sort((a,b) => new Date(b.at) - new Date(a.at));
  }, [patientId]);

  const patientName = rows[0]?.patientName || `Patient #${patientId}`;

  return (
    <div>
      <div className={styles.header}>
        <h2 className={styles.title}>Appointment</h2>
        <button className={styles.back} onClick={()=>nav(-1)}>← Back</button>
      </div>

      <div className={styles.meta}>
        <div className={styles.badge}>ผู้ป่วย</div>
        <div className={styles.name}>{patientName}</div>
      </div>

      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th style={{width:'8%'}}>#</th>
              <th style={{width:'22%'}}>วันที่</th>
              <th style={{width:'20%'}}>เวลา</th>
              <th>หมายเหตุ</th>
              <th style={{width:'14%'}}></th>
            </tr>
          </thead>
          <tbody>
            {rows.length === 0 ? (
              <tr>
                <td colSpan={5} style={{textAlign:'center', color:'#6b7280', height:56}}>
                  ยังไม่มีการนัดหมาย
                </td>
              </tr>
            ) : (
              rows.map((r, i) => {
                const isLatest = i === 0; // แถวล่าสุด
                return (
                  <tr key={r.id} className={isLatest ? styles.latestRow : undefined}>
                    <td>{i+1}</td>
                    <td>{formatDateTH(r.at)}</td>
                    <td>{formatTimeTH(r.at)}</td>
                    <td className={styles.note}>{r.note || '-'}</td>
                    <td className={styles.statusCell}>
                      {isLatest && <span className={styles.latestBadge}>ล่าสุด</span>}
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
