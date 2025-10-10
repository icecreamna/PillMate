import { Routes, Route, Navigate } from 'react-router-dom'
import Login from './pages/Login.jsx'

/** Admin (superadmin) */
import AdminLayout from './pages/admin/AdminLayout.jsx'

import DoctorList from './pages/admin/doctor/DoctorList.jsx'
import AddDoctor from './pages/admin/doctor/AddDoctor.jsx'
import EditDoctor from './pages/admin/doctor/EditDoctor.jsx'
import ResetDoctorPassword from './pages/admin/doctor/ResetDoctorPassword.jsx'

import MedicineList from './pages/admin/medicine/MedicineList.jsx'
import AddMedicine from './pages/admin/medicine/AddMedicine.jsx'
import ViewMedicine from './pages/admin/medicine/ViewMedicine.jsx'
import EditMedicine from './pages/admin/medicine/EditMedicine.jsx'

/** Doctor */
import DocLayout from './pages/doctor/DocLayout.jsx'

import PatientList from './pages/doctor/patient/PatientList.jsx'
import AddPatient from './pages/doctor/patient/AddPatient.jsx'
import ViewPatient from './pages/doctor/patient/ViewPatient.jsx';
import EditPatient from './pages/doctor/patient/EditPatient.jsx';

import DocMedList from './pages/doctor/medicine/MedicineList.jsx'
import DocMedView from './pages/doctor/medicine/ViewMedicine.jsx'

import PrescriptionList from './pages/doctor/prescription/PrescriptionList.jsx'
import AddPrescription from './pages/doctor/prescription/AddPrescription.jsx';
import ViewPrescription from './pages/doctor/prescription/ViewPrescription.jsx';

import AppointmentList from './pages/doctor/appointment/AppointmentList.jsx'
import ViewAppointment from './pages/doctor/appointment/ViewAppointment.jsx';

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
  if (role === 'doctor') return <Navigate to="/doc" replace />
  return <Navigate to="/login" replace />
}

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<HomeRedirect />} />
      <Route path="/login" element={<Login />} />

      {/* -------- Admin -------- */}
      <Route path="/admin" element={
        <PrivateRoute needRole="superadmin"><AdminLayout /></PrivateRoute>
      }>
        <Route index element={<DoctorList />} />
        <Route path="add" element={<AddDoctor />} />
        <Route path=":id/edit" element={<EditDoctor />} />
        <Route path=":id/reset-password" element={<ResetDoctorPassword />} />

        <Route path="medicine-info" element={<MedicineList />} />
        <Route path="medicine-info/add" element={<AddMedicine />} />
        <Route path="medicine-info/:id" element={<ViewMedicine />} />
        <Route path="medicine-info/:id/edit" element={<EditMedicine />} />
      </Route>

      {/* -------- Doctor -------- */}
      <Route path="/doc" element={
        <PrivateRoute needRole="doctor"><DocLayout /></PrivateRoute>
      }>
        <Route index element={<PatientList />} />
        <Route path="patients" element={<PatientList />} />
        <Route path="patients/add" element={<AddPatient/>} />
        <Route path="patients/:id" element={<ViewPatient/>} />
        <Route path="patients/:id/edit" element={<EditPatient/>} />

        <Route path="medicine-info" element={<DocMedList />} />
        <Route path="medicine-info/:id" element={<DocMedView />} />
        
        <Route path="prescription" element={<PrescriptionList />} />
        <Route path="prescription/add/:patientId" element={<AddPrescription />} />
        <Route path="prescription/view/:id" element={<ViewPrescription />} />
        

        <Route path="appointment" element={<AppointmentList />} />
        <Route path="appointment/view/:id" element={<ViewAppointment />} />
        <Route path="appointment/add/:id" element={<div style={{padding:16}}>Add Appointment (stub)</div>} />
      </Route>

      <Route path="*" element={<Navigate to="/login" replace />} />
    </Routes>
  )
}
