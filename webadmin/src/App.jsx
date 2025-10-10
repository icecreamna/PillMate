// src/App.jsx
import { Routes, Route, Navigate } from 'react-router-dom'
import Login from './pages/Login.jsx'
import AdminLayout from './pages/admin/AdminLayout.jsx'
import DoctorList from './pages/admin/DoctorList.jsx'
import AddDoctor from './pages/admin/AddDoctor.jsx'
import EditDoctor from './pages/admin/EditDoctor.jsx'
import ResetDoctorPassword from './pages/admin/ResetDoctorPassword.jsx'

function PrivateRoute({ children, needRole }) {
  const token = localStorage.getItem('auth_token')
  const role = localStorage.getItem('role')
  if (!token) return <Navigate to="/login" replace />
  if (needRole && role !== needRole) return <Navigate to="/login" replace />
  return children
}

function HomeRedirect() {
  const role = localStorage.getItem('role')
  if (role === 'superadmin') return <Navigate to="/admin" replace />
  return <Navigate to="/login" replace />
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<HomeRedirect />} />
      <Route path="/login" element={<Login />} />

      {/* โซนผู้ดูแลระบบ (ต้องเป็น superadmin เท่านั้น) */}
      <Route
        path="/admin"
        element={
          <PrivateRoute needRole="superadmin">
            <AdminLayout />
          </PrivateRoute>
        }
      >
        <Route index element={<DoctorList />} />
        <Route path="add" element={<AddDoctor />} />
        <Route path=":id/edit" element={<EditDoctor />} />
        <Route path=":id/reset-password" element={<ResetDoctorPassword />} />
        <Route path="medicine-info" element={<div style={{ padding: 16 }}>MedicineInfo</div>} />
      </Route>

      <Route path="*" element={<Navigate to="/login" replace />} />
    </Routes>
  )
}
