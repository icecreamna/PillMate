// src/App.jsx
import { Routes, Route, Navigate } from 'react-router-dom'
import Login from './pages/Login.jsx'
import AdminLayout from './pages/admin/AdminLayout.jsx'
import DoctorList from './pages/admin/doctor/DoctorList.jsx'
import AddDoctor from './pages/admin/doctor/AddDoctor.jsx'
import EditDoctor from './pages/admin/doctor/EditDoctor.jsx'
import ResetDoctorPassword from './pages/admin/doctor/ResetDoctorPassword.jsx'
import MedicineList from './pages/admin/medicine/MedicineList.jsx'
import AddMedicine from './pages/admin/medicine/AddMedicine.jsx'
import ViewMedicine from './pages/admin/medicine/ViewMedicine.jsx'
import EditMedicine from './pages/admin/medicine/EditMedicine.jsx'

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
        {/* หน้าเริ่มต้นภายใต้ /admin แสดงรายการหมอ */}
        <Route index element={<DoctorList />} />
        <Route path="add" element={<AddDoctor />} />
        <Route path=":id/edit" element={<EditDoctor />} />
        <Route path=":id/reset-password" element={<ResetDoctorPassword />} />

        {/* Medicine */}
        <Route path="medicine-info" element={<MedicineList />} />
        <Route path="medicine-info/add" element={<AddMedicine />} />
        <Route path="medicine-info/:id" element={<ViewMedicine />} />
        <Route path="medicine-info/:id/edit" element={<EditMedicine />} /> 
        {/* ตัวอย่างเส้นทางที่จะเพิ่มต่อไป:
        <Route path="medicine-info/add" element={<AddMedicine />} />
        <Route path="medicine-info/:id" element={<ViewMedicine />} />
        <Route path="medicine-info/:id/edit" element={<EditMedicine />} />
        */}
      </Route>

      <Route path="*" element={<Navigate to="/login" replace />} />
    </Routes>
  )
}
