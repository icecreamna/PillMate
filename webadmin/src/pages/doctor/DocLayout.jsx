import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import styles from '../../styles/doctor/DocLayout.module.css'
import LogoutButton from '../../components/LogoutButton'

export default function DocLayout(){
  const nav = useNavigate()
  const email = localStorage.getItem('email') || 'hospital_a@pillmate.com'

  return (
    <div className={styles.app}>
      <aside className={styles.sidebar}>
        <div className={styles.logo}>Doctor</div>
        <nav className={styles.menu}>
          <NavLink to="/doc/patients" end className={({isActive})=>isActive?styles.active:styles.link}>👥 Patients</NavLink>
          <NavLink to="/doc/medicine-info" className={({isActive})=>isActive?styles.active:styles.link}>💊 MedicineInfo</NavLink>
          <NavLink to="/doc/prescription" className={({isActive})=>isActive?styles.active:styles.link}>🧾 Prescription</NavLink>
          <NavLink to="/doc/appointment" className={({isActive})=>isActive?styles.active:styles.link}>📅 Appointment</NavLink>
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
