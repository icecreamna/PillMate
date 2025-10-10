import { useNavigate } from 'react-router-dom'
import styles from './Dashboard.module.css'

export default function Dashboard() {
  const nav = useNavigate()
  const logout = () => {
    localStorage.removeItem('auth_token')
    nav('/login', { replace: true })
  }

  return (
    <div className={styles.wrap}>
      <div className={styles.card}>
        <h1 className={styles.title}>Dashboard</h1>
        <p>ยินดีต้อนรับ 🎉 (หน้านี้กันด้วย PrivateRoute)</p>
        <button onClick={logout} className={styles.logoutBtn}>ออกจากระบบ</button>
      </div>
    </div>
  )
}
