import { useMemo, useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/prescription/ViewPrescription.module.css'

// ----- MOCK DATA: ประวัติการสั่งยา (หลายครั้ง) -----
// โดยปกติจะ GET จาก API: /api/patients/:id/prescriptions
const MOCK_PRESCRIPTIONS = [
  {
    id: 'rx-1001',
    patientId: 1,
    orderedAt: '2025-03-10T09:25:00+07:00',
    doctorName: 'นพ. สมปอง ใจดี',
    items: [
      { id: 1, medName: 'Medicine A', generic: 'Generic A', strength: '500 mg', form: 'ยาเม็ด', dosePerTime: 1, timesPerDay: 3 },
      { id: 2, medName: 'Medicine B', generic: 'Generic B', strength: '400 mg', form: 'แคปซูล', dosePerTime: 1, timesPerDay: 2 },
    ],
  },
  {
    id: 'rx-1002',
    patientId: 1,
    orderedAt: '2025-03-10T15:40:00+07:00',
    doctorName: 'นพ. สมปอง ใจดี',
    items: [
      { id: 1, medName: 'Medicine A', generic: 'Generic A', strength: '500 mg', form: 'ยาเม็ด', dosePerTime: 2, timesPerDay: 3 },
    ],
  },
  {
    id: 'rx-1003',
    patientId: 1,
    orderedAt: '2025-02-28T10:05:00+07:00',
    doctorName: 'พญ. สมหญิง สุขดี',
    items: [
      { id: 3, medName: 'Medicine C', generic: 'Generic C', strength: '4 mg/5 ml', form: 'ยาน้ำ', dosePerTime: 10, timesPerDay: 2 },
    ],
  },
];

// helper: เวลาโซนไทยในรูปแบบ ISO+07:00
function nowTHISO() {
  const tzOffsetMin = -420; // Asia/Bangkok = UTC+7 => offset -420 นาทีสำหรับ toISOString แบบกำหนดเอง
  const now = new Date();
  const local = new Date(now.getTime() - now.getTimezoneOffset()*60000 + tzOffsetMin*60000);
  const yyyy = local.getUTCFullYear();
  const mm = String(local.getUTCMonth()+1).padStart(2,'0');
  const dd = String(local.getUTCDate()).padStart(2,'0');
  const hh = String(local.getUTCHours()).padStart(2,'0');
  const mi = String(local.getUTCMinutes()).padStart(2,'0');
  const ss = String(local.getUTCSeconds()).padStart(2,'0');
  return `${yyyy}-${mm}-${dd}T${hh}:${mi}:${ss}+07:00`;
}

function formatDateTH(d) {
  const dt = new Date(d);
  return dt.toLocaleDateString('th-TH', { day:'numeric', month:'short', year:'numeric' });
}
function formatTimeTH(d) {
  const dt = new Date(d);
  return dt.toLocaleTimeString('th-TH', { hour:'2-digit', minute:'2-digit' });
}

export default function ViewPrescription() {
  const nav = useNavigate();
  const { id: patientId } = useParams();

  // เริ่มต้นกรอง mock ตาม patientId
  const initial = useMemo(
    () => MOCK_PRESCRIPTIONS.filter(p => String(p.patientId) === String(patientId ?? 1)),
    [patientId]
  );

  // เก็บเป็น state เพื่อให้ลบ/ทำซ้ำแล้วอัปเดต UI ได้
  const [data, setData] = useState(initial);
  useEffect(() => { setData(initial); }, [initial]);

  // กลุ่มตาม "วันที่" เดียวกัน + เรียงใหม่สุดก่อนทั้งระดับกลุ่มและในกลุ่ม
  const groups = useMemo(() => {
    const map = new Map();
    for (const rx of data) {
      const key = formatDateTH(rx.orderedAt);
      if (!map.has(key)) map.set(key, []);
      map.get(key).push(rx);
    }
    const sorted = [...map.entries()].sort(([a],[b]) => {
      const ad = new Date(a); const bd = new Date(b);
      return bd - ad;
    });
    for (const [, arr] of sorted) arr.sort((a,b) => new Date(b.orderedAt) - new Date(a.orderedAt));
    return sorted;
  }, [data]);

  // ลบใบสั่งยา
  const onDelete = (rxId) => {
    if (!confirm('ยืนยันลบใบสั่งยานี้?')) return;
    // TODO: DELETE API แล้วค่อย setData ตามผลลัพธ์
    setData(prev => prev.filter(rx => rx.id !== rxId));
  };

  // ทำซ้ำใบสั่งยา -> เพิ่มรายการใหม่เป็นอันล่าสุดใน state
  const onRepeat = (rx) => {
    // TODO: POST API เพื่อสร้างใบสั่งใหม่ แล้วใช้ผลลัพธ์จริงแทน mock ที่สร้างเอง
    const newRx = {
      ...rx,
      id: `rx-${Date.now()}`,            // id ใหม่ง่าย ๆ จาก timestamp
      orderedAt: nowTHISO(),             // เวลา “ตอนนี้” โซนไทย
      // ถ้าต้องการ set ชื่อหมอจากผู้ล็อกอินจริง ๆ ค่อยดึงจาก auth/localStorage ได้
    };
    setData(prev => [newRx, ...prev]);   // push เข้าไปแล้วให้เป็นอันแรกสุด
  };

  return (
    <div>
      <div className={styles.header}>
        <h2 className={styles.title}>View Prescription</h2>
        <button className={styles.back} onClick={() => nav(-1)}>← Back</button>
      </div>

      {groups.length === 0 ? (
        <div className={styles.empty}>ยังไม่มีประวัติการสั่งยา</div>
      ) : (
        groups.map(([dateLabel, rxs]) => (
          <section key={dateLabel} className={styles.section}>
            <div className={styles.dateBadge}>{dateLabel}</div>

            {rxs.map((rx) => (
              <div key={rx.id} className={styles.card}>
                <div className={styles.cardHeader}>
                  <div className={styles.rxMeta}>
                    <span className={styles.time}>{formatTimeTH(rx.orderedAt)}</span>
                    <span className={styles.dot}>•</span>
                    <span className={styles.doctor}>แพทย์ผู้สั่ง: {rx.doctorName}</span>
                    <span className={styles.dot}>•</span>
                    <span className={styles.rxId}>เลขที่ใบสั่ง: {rx.id}</span>
                  </div>

                  <div className={styles.actionsRow}>
                    <button
                      className={styles.actionPrimary}
                      onClick={() => onRepeat(rx)}
                    >
                      ทำซ้ำ
                    </button>
                    <button
                      className={styles.actionDanger}
                      onClick={() => onDelete(rx.id)}
                    >
                      ลบ
                    </button>
                  </div>
                </div>

                <div className={styles.tableWrap}>
                  <table className={styles.table}>
                    <thead>
                      <tr>
                        <th style={{width:'6%'}}>#</th>
                        <th style={{width:'26%'}}>MedName</th>
                        <th style={{width:'26%'}}>GenericName</th>
                        <th style={{width:'14%'}}>Strength</th>
                        <th style={{width:'12%'}}>Form</th>
                        <th style={{width:'16%'}}>Dosage</th>
                      </tr>
                    </thead>
                    <tbody>
                      {rx.items.map((it, idx) => (
                        <tr key={idx}>
                          <td>
                            <input type="checkbox" checked readOnly className={styles.checkbox}/>
                          </td>
                          <td>{it.medName}</td>
                          <td>{it.generic}</td>
                          <td>{it.strength}</td>
                          <td>{it.form}</td>
                          <td className={styles.dose}>
                            ครั้งละ <strong>{it.dosePerTime}</strong> | วันละ <strong>{it.timesPerDay}</strong>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            ))}
          </section>
        ))
      )}
    </div>
  );
}
