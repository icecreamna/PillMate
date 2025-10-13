import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import styles from "../../styles/admin/AdminLayout.module.css"
import LogoutButton from "../../components/LogoutButton"

export default function AdminLayout() {
  const nav = useNavigate()
  const email = localStorage.getItem('email') || 'super_admin@pillmate.com'

  return (
    <div className={styles.app}>
      <aside className={styles.sidebar}>
        <div className={styles.logo}>Admin</div>
        <nav className={styles.menu}>
          <NavLink
            to="/admin"
            end
            className={({ isActive }) => (isActive ? styles.active : styles.link)}
          >
            <span className={styles.icon}>👤</span> Doctors
          </NavLink>
          <NavLink
            to="/admin/medicine-info"
            className={({ isActive }) => (isActive ? styles.active : styles.link)}
          >
            <span className={styles.icon}>💊</span> MedicineInfo
          </NavLink>
        </nav>
      </aside>

      <main className={styles.main}>
        <header className={styles.topbar}>
          <input className={styles.email} value={email} disabled />
          <LogoutButton
            className={styles.logout}
            confirm
            onAfter={() => nav('/login', { replace: true })}
          >
            ออกจากระบบ
          </LogoutButton>
        </header>

        <div className={styles.content}>
          <Outlet />
        </div>
      </main>
    </div>
  )
}
