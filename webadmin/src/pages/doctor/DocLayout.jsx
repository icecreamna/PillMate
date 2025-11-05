// src/pages/doctor/DocLayout.jsx
import { useEffect, useRef, useState } from "react";
import { NavLink, Outlet, useNavigate } from "react-router-dom";
import styles from "../../styles/doctor/DocLayout.module.css";
import LogoutButton from "../../components/LogoutButton";

export default function DocLayout() {
  const nav = useNavigate();
  const email = localStorage.getItem("email") || "hospital_a@pillmate.com";

  // dropdown state
  const [open, setOpen] = useState(false);
  const btnRef = useRef(null);
  const menuRef = useRef(null);

  // ปิด dropdown เมื่อคลิกนอก/กด Esc
  useEffect(() => {
    const onClick = (e) => {
      if (!open) return;
      if (
        btnRef.current &&
        !btnRef.current.contains(e.target) &&
        menuRef.current &&
        !menuRef.current.contains(e.target)
      )
        setOpen(false);
    };
    const onKey = (e) => {
      if (e.key === "Escape") setOpen(false);
    };
    document.addEventListener("mousedown", onClick);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onClick);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  return (
    <div className={styles.app}>
      <aside className={styles.sidebar}>
        <div className={styles.logo}>Doctor</div>
        <nav className={styles.menu}>
          <NavLink
            to="/doc/patients"
            end
            className={({ isActive }) =>
              isActive ? styles.active : styles.link
            }
          >
            👥 Patients
          </NavLink>
          <NavLink
            to="/doc/medicine-info"
            className={({ isActive }) =>
              isActive ? styles.active : styles.link
            }
          >
            💊 MedicineInfo
          </NavLink>
          <NavLink
            to="/doc/prescription"
            className={({ isActive }) =>
              isActive ? styles.active : styles.link
            }
          >
            🧾 Prescription
          </NavLink>
          <NavLink
            to="/doc/appointment"
            className={({ isActive }) =>
              isActive ? styles.active : styles.link
            }
          >
            📅 Appointment
          </NavLink>
        </nav>
      </aside>

      <main className={styles.main}>
        <header className={styles.topbar}>
          <div /> {/* ดันให้เมนูผู้ใช้ไปชิดขวา */}
          {/* เมนูผู้ใช้ (dropdown) */}
          <div className={styles.userMenu}>
            <button
              ref={btnRef}
              type="button"
              className={styles.userBtn}
              onClick={() => setOpen((o) => !o)}
              aria-haspopup="menu"
              aria-expanded={open ? "true" : "false"}
              title={email}
            >
              <span className={styles.userAvatar} aria-hidden>
                👩‍⚕️
              </span>
              <span className={styles.userEmail}>{email}</span>
              <span className={styles.caret} aria-hidden>
                ▾
              </span>
            </button>

            {open && (
              <div ref={menuRef} role="menu" className={styles.dropdown}>
                <button
                  role="menuitem"
                  type="button"
                  className={styles.dropdownItem}
                  onClick={() => {
                    setOpen(false);
                    nav("/doc/profile");
                  }} // ✅ เปลี่ยน path
                >
                  ⚙️ จัดการบัญชี
                </button>

                <div className={styles.separator} />

                <LogoutButton
                  role="menuitem"
                  className={`${styles.dropdownItem} ${styles.logoutItem}`}
                  confirm
                  onAfter={() => nav("/login", { replace: true })}
                >
                  ออกจากระบบ
                </LogoutButton>
              </div>
            )}
          </div>
        </header>

        <div className={styles.content}>
          <Outlet />
        </div>
      </main>
    </div>
  );
}
