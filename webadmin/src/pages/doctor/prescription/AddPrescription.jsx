import { useEffect, useMemo, useState, useRef } from "react";
import { useNavigate, useParams } from "react-router-dom";
import styles from "../../../styles/doctor/prescription/AddPrescription.module.css";
import { getPatient } from "../../../services/patients";
import { listMedicinesForDoctor } from "../../../services/medicines";
import { listForms } from "../../../services/initialData";
import { createPrescription } from "../../../services/prescriptions";

function inferUnitFromForm(formName) {
  const s = String(formName || "").trim();
  if (!s) return "";
  if (s.includes("เม็ด")) return "เม็ด";
  if (s.includes("แคปซูล")) return "แคปซูล";
  if (s.includes("ยาน้ำ")) return "ml";
  if (s.includes("ขี้ผึ้ง")) return "กรัม";
  if (s.includes("สเปรย์")) return "สเปรย์";
  return "";
}

function norm(s) {
  return String(s || "")
    .toLocaleLowerCase("th-TH")
    .normalize("NFKD")
    .replace(/[\u0300-\u036f]/g, "")
    .replace(/\s+/g, "")
    .replace(/[^\p{L}\p{N}.%-]/gu, "");
}

/** YYYY-MM-DD -> DD/MM/YYYY */
function isoToDDMMYYYY(iso) {
  const s = String(iso || "");
  if (!/^\d{4}-\d{2}-\d{2}$/.test(s)) return "";
  const [y, m, d] = s.split("-");
  return `${d}/${m}/${y}`;
}

export default function AddPrescription() {
  const nav = useNavigate();
  const { patientId } = useParams();

  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  const [patient, setPatient] = useState(null);
  const [rows, setRows] = useState([]);

  const [q, setQ] = useState("");

  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState(null);
  const [dosePerTime, setDosePerTime] = useState("1");
  const [timesPerDay, setTimesPerDay] = useState("3");

  const [startDateISO, setStartDateISO] = useState("");
  const [endDateISO, setEndDateISO] = useState("");
  const [note, setNote] = useState("");

  const startRef = useRef(null);
  const endRef = useRef(null);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      setLoading(true);
      setError("");
      try {
        const pRes = await getPatient(patientId);
        const p = pRes?.data;
        if (!p) throw new Error("ไม่พบข้อมูลผู้ป่วย");
        if (cancelled) return;
        setPatient(p);

        const fs = await listForms({ with_relations: true });
        const forms = Array.isArray(fs) ? fs : [];
        const formMap = new Map(forms.map((f) => [String(f.id), f]));
        const unitLookupByForm = new Map(
          forms.map((f) => [
            String(f.id),
            new Map(
              (Array.isArray(f.units) ? f.units : []).map((u) => [
                String(u.id),
                u.unit_name || u.name || "",
              ])
            ),
          ])
        );

        const mRes = await listMedicinesForDoctor();
        const meds = Array.isArray(mRes?.data) ? mRes.data : [];

        const mapped = meds.map((m) => {
          const fid = String(m.form_id ?? "");
          const form = formMap.get(fid) || null;
          const formName =
            m.form_name || m.form || form?.form_name || form?.name || "-";

          const unitFromMed = String(m.unit_name || m.unit || "").trim();
          const unitFromLookup =
            unitLookupByForm.get(fid)?.get(String(m.unit_id ?? "")) || "";
          const unitLabel =
            unitFromMed || unitFromLookup || inferUnitFromForm(formName) || "";

          return {
            id: m.id,
            medName: m.med_name || m.MedName || "-",
            generic: m.generic_name || m.GenericName || "-",
            strength: m.strength || "-",
            form: formName,
            unitLabel,
            checked: false,
            dosePerTime: undefined,
            timesPerDay: undefined,
            startDate: "",
            endDate: "",
            note: "",
          };
        });

        if (!cancelled) setRows(mapped);
      } catch (e) {
        if (!cancelled) setError(e.message || "โหลดข้อมูลไม่สำเร็จ");
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [patientId]);

  const initial = useMemo(() => rows.map((d) => ({ ...d })), [rows]);

  const toggle = (id) => {
    setRows((rs) =>
      rs.map((r) => (r.id === id ? { ...r, checked: !r.checked } : r))
    );
  };

  const specify = (row) => {
    setEditing(row);
    setDosePerTime(String(row.dosePerTime ?? "1"));
    setTimesPerDay(String(row.timesPerDay ?? "3"));
    setStartDateISO(row.startDate || "");
    setEndDateISO(row.endDate || "");
    setNote(String(row.note ?? ""));
    setOpen(true);
  };
  const closeModal = () => {
    setOpen(false);
    setEditing(null);
  };

  const confirmSpecify = () => {
    if (startDateISO && endDateISO) {
      const a = new Date(startDateISO);
      const b = new Date(endDateISO);
      if (!Number.isNaN(a.getTime()) && !Number.isNaN(b.getTime()) && a > b) {
        setError("วันสิ้นสุดต้องไม่น้อยกว่าวันเริ่มต้น");
        return;
      }
    }

    setRows((rs) =>
      rs.map((r) =>
        r.id === editing.id
          ? {
              ...r,
              dosePerTime: Number(String(dosePerTime || "0")),
              timesPerDay: Number(String(timesPerDay || "0")),
              startDate: startDateISO || "",
              endDate: endDateISO || "",
              note: String(note || "").trim(),
            }
          : r
      )
    );
    closeModal();
  };

  const visibleRows = useMemo(() => {
    const nq = norm(q);
    if (!nq) return rows;
    return rows.filter(
      (r) =>
        norm(r.medName).includes(nq) ||
        norm(r.generic).includes(nq) ||
        norm(r.form).includes(nq) ||
        norm(r.strength).includes(nq)
    );
  }, [rows, q]);

  const submit = async () => {
    setError("");
    const selected = rows.filter((r) => r.checked);
    if (selected.length === 0) {
      setError("กรุณาเลือกรายการยาอย่างน้อย 1 รายการ");
      return;
    }
    const idCard = patient?.id_card_number;
    if (!idCard) {
      setError("ไม่พบเลขบัตรประชาชนของผู้ป่วย");
      return;
    }

    const toNumStr = (v) => {
      const s = String(v ?? "").trim();
      return s === "" ? "0" : s;
    };

    const items = selected.map((s) => {
      const payload = {
        medicine_info_id: s.id,
        amount_per_time: toNumStr(s.dosePerTime),
        times_per_day: toNumStr(s.timesPerDay),
      };
      if (s.startDate) payload.start_date = s.startDate;
      if (s.endDate) payload.end_date = s.endDate;
      if (s.note) payload.note = s.note;
      return payload;
    });

    try {
      setSaving(true);
      await createPrescription({ id_card_number: idCard, items });
      nav("/doc/prescription", { replace: true });
    } catch (e) {
      setError(e.message || "บันทึกไม่สำเร็จ");
    } finally {
      setSaving(false);
    }
  };

  const openStartPicker = () => startRef.current?.showPicker?.();
  const openEndPicker = () => endRef.current?.showPicker?.();

  return (
    <div>
      <div className={styles.header}>
        <h2 className={styles.title}>Add Prescription</h2>
        <button className={styles.back} onClick={() => nav(-1)}>
          ← Back
        </button>
      </div>

      {patient && (
        <div className={styles.patientBar}>
          <div>
            <strong>ชื่อผู้ป่วย:</strong>{" "}
            {[patient.first_name, patient.last_name]
              .filter(Boolean)
              .join(" ") || "-"}
          </div>
          <div>
            <strong>Patient Code:</strong> {patient.patient_code || "-"}
          </div>
          <div>
            <strong>เลขบัตร:</strong> {patient.id_card_number || "-"}
          </div>
        </div>
      )}

      <div className={styles.toolsBar}>
        <input
          type="text"
          className={styles.searchInput}
          value={q}
          onChange={(e) => setQ(e.target.value)}
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
                  <th style={{ width: "6%" }}>#</th>
                  <th style={{ width: "26%" }}>MedName</th>
                  <th style={{ width: "22%" }}>GenericName</th>
                  <th style={{ width: "14%" }}>Strength</th>
                  <th style={{ width: "12%" }}>Form</th>
                  <th style={{ width: "20%" }}>รายละเอียด</th>
                </tr>
              </thead>
              <tbody>
                {visibleRows.map((r) => (
                  <tr key={r.id} className={r.checked ? "" : styles.dim}>
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
                      <div className={styles.medName}>{r.medName}</div>
                      {r.dosePerTime != null && r.timesPerDay != null ? (
                        <div className={styles.doseNote}>
                          ครั้งละ {r.dosePerTime} {r.unitLabel} | วันละ{" "}
                          {r.timesPerDay} ครั้ง
                        </div>
                      ) : null}
                    </td>
                    <td>{r.generic}</td>
                    <td>{r.strength}</td>
                    <td>{r.form}</td>
                    <td className={styles.actions}>
                      {(r.startDate || r.endDate || r.note) && (
                        <div className={styles.brief}>
                          {r.startDate && (
                            <span>เริ่ม {isoToDDMMYYYY(r.startDate)}</span>
                          )}
                          {r.endDate && (
                            <span style={{ marginLeft: 8 }}>
                              สิ้นสุด {isoToDDMMYYYY(r.endDate)}
                            </span>
                          )}
                          {r.note && (
                            <div className={styles.noteBrief}>
                              note: {r.note}
                            </div>
                          )}
                        </div>
                      )}
                      <button
                        className={styles.specify}
                        onClick={() => specify(r)}
                      >
                        specify
                      </button>
                    </td>
                  </tr>
                ))}
                {visibleRows.length === 0 && (
                  <tr>
                    <td
                      colSpan={6}
                      style={{
                        textAlign: "center",
                        color: "#6b7280",
                        height: 56,
                      }}
                    >
                      ไม่พบรายการที่ตรงกับคำค้น
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>

          <div className={styles.footer}>
            <button
              className={styles.submit}
              onClick={submit}
              disabled={saving || rows.length === 0}
            >
              {saving ? "กำลังบันทึก…" : "จ่ายยา"}
            </button>
          </div>
        </>
      )}

      {open && editing && (
        <div className={styles.modalOverlay} onClick={closeModal}>
          <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
            <button
              className={styles.modalClose}
              onClick={closeModal}
              aria-label="Close"
            >
              ×
            </button>

            <div className={styles.modalHead}>
              <div className={styles.modalTitle}>
                <strong>{editing.medName}</strong>
              </div>
              <div className={styles.modalSub}>{editing.generic}</div>
            </div>

            <div className={styles.modalContent}>
              <div className={styles.hRow}>
                <label className={styles.hField}>
                  <span className={styles.hLabel}>ครั้งละ</span>
                  <input
                    type="number"
                    min="0"
                    className={styles.hInput}
                    value={dosePerTime}
                    onChange={(e) => setDosePerTime(e.target.value)}
                  />
                  <span className={styles.hUnit}>{editing.unitLabel}</span>
                </label>

                <label className={styles.hField}>
                  <span className={styles.hLabel}>วันละ</span>
                  <input
                    type="number"
                    min="0"
                    className={styles.hInput}
                    value={timesPerDay}
                    onChange={(e) => setTimesPerDay(e.target.value)}
                  />
                  <span className={styles.hUnit}>ครั้ง</span>
                </label>
              </div>

              {/* Date fields */}
              <div className={styles.hRow}>
                <label className={styles.hField}>
                  <span className={styles.hLabel}>เริ่ม</span>
                  <div className={styles.dateWrap}>
                    <span className={styles.dateDisplay} aria-hidden="true">
                      {startDateISO
                        ? isoToDDMMYYYY(startDateISO)
                        : "dd/mm/yyyy"}
                    </span>
                    <input
                      ref={startRef}
                      type="date"
                      className={styles.dateNative}
                      value={startDateISO}
                      onChange={(e) => setStartDateISO(e.target.value)}
                    />
                    <button
                      type="button"
                      className={styles.calBtn}
                      aria-label="เปิดปฏิทินเลือกวันที่เริ่ม"
                      onClick={() => startRef.current?.showPicker?.()}
                    >
                      {/* calendar icon (inline SVG) */}
                      <svg
                        viewBox="0 0 24 24"
                        width="16"
                        height="16"
                        aria-hidden="true"
                      >
                        <path
                          d="M7 2v2M17 2v2M3 8h18M4 6h16a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V7a1 1 0 0 1 1-1Z"
                          fill="none"
                          stroke="currentColor"
                          strokeWidth="1.8"
                          strokeLinecap="round"
                          strokeLinejoin="round"
                        />
                      </svg>
                    </button>
                  </div>
                </label>

                <label className={styles.hField}>
                  <span className={styles.hLabel}>สิ้นสุด</span>
                  <div className={styles.dateWrap}>
                    <span className={styles.dateDisplay} aria-hidden="true">
                      {endDateISO ? isoToDDMMYYYY(endDateISO) : "dd/mm/yyyy"}
                    </span>
                    <input
                      ref={endRef}
                      type="date"
                      className={styles.dateNative}
                      value={endDateISO}
                      onChange={(e) => setEndDateISO(e.target.value)}
                    />
                    <button
                      type="button"
                      className={styles.calBtn}
                      aria-label="เปิดปฏิทินเลือกวันที่สิ้นสุด"
                      onClick={() => endRef.current?.showPicker?.()}
                    >
                      <svg
                        viewBox="0 0 24 24"
                        width="16"
                        height="16"
                        aria-hidden="true"
                      >
                        <path
                          d="M7 2v2M17 2v2M3 8h18M4 6h16a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V7a1 1 0 0 1 1-1Z"
                          fill="none"
                          stroke="currentColor"
                          strokeWidth="1.8"
                          strokeLinecap="round"
                          strokeLinejoin="round"
                        />
                      </svg>
                    </button>
                  </div>
                </label>

                <div className={styles.hHint}>
                  * ถ้าระบุวันสิ้นสุด ระบบจะตั้งวันหมดอายุ = สิ้นสุด + 1 วัน
                </div>
              </div>

              <div className={styles.hRowNote}>
                <textarea
                  className={styles.hTextarea}
                  rows={2}
                  placeholder="วิธีใช้/ข้อควรระวัง (ไม่บังคับ)"
                  value={note}
                  onChange={(e) => setNote(e.target.value)}
                />
              </div>
            </div>

            <div className={styles.modalActions}>
              <button className={styles.confirm} onClick={confirmSpecify}>
                ยืนยัน
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
