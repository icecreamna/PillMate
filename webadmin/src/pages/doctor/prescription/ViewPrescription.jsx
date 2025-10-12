import { useMemo, useState, useEffect } from 'react'
import { useNavigate, useParams } from 'react-router-dom'
import styles from '../../../styles/doctor/prescription/ViewPrescription.module.css'

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

  const initial = useMemo(
    () => MOCK_PRESCRIPTIONS.filter(p => String(p.patientId) === String(patientId ?? 1)),
    [patientId]
  );

  const [data, setData] = useState(initial);
  useEffect(() => { setData(initial); }, [initial]);

  const latestAt = useMemo(() => {
    if (!data.length) return null;
    return Math.max(...data.map(rx => new Date(rx.orderedAt).getTime()));
  }, [data]);

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

            {rxs.map((rx) => {
              const isLatest = latestAt && (new Date(rx.orderedAt).getTime() === latestAt);
              return (
                <div
                  key={rx.id}
                  className={`${styles.card} ${isLatest ? styles.cardLatest : ''}`}
                >
                  <div className={styles.cardHeader}>
                    <div className={styles.rxMeta}>
                      <span className={styles.time}>{formatTimeTH(rx.orderedAt)}</span>
                      <span className={styles.dot}>•</span>
                      <span className={styles.doctor}>แพทย์ผู้สั่ง: {rx.doctorName}</span>
                      {isLatest && <span className={styles.latestTag}>ล่าสุด</span>}
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
              );
            })}
          </section>
        ))
      )}
    </div>
  );
}
