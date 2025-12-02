import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import Navbar from './components/Navbar';
import ProtectedRoute from './components/ProtectedRoute';
import Home from './pages/Home';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Organizations from './pages/Organizations';
import Locations from './pages/Locations';
import Vehicles from './pages/Vehicles';
import VehicleProfile from './pages/VehicleProfile';
import VehicleForm from './pages/VehicleForm';
import VehicleBulkUpload from './pages/VehicleBulkUpload';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min';
import 'bootstrap-icons/font/bootstrap-icons.css';
import './bootstrap-theme.css';

function App() {
  return (
    <AuthProvider>
      <Router>
        <div className="d-flex flex-column h-100">
          <Navbar />
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login />} />
            <Route
              path="/dashboard"
              element={
                <ProtectedRoute>
                  <Dashboard />
                </ProtectedRoute>
              }
            />
            <Route
              path="/organizations"
              element={
                <ProtectedRoute>
                  <Organizations />
                </ProtectedRoute>
              }
            />
            <Route
              path="/locations"
              element={
                <ProtectedRoute>
                  <Locations />
                </ProtectedRoute>
              }
            />
            <Route
              path="/vehicles"
              element={
                <ProtectedRoute>
                  <Vehicles />
                </ProtectedRoute>
              }
            />
            <Route
              path="/vehicles/new"
              element={
                <ProtectedRoute>
                  <VehicleForm />
                </ProtectedRoute>
              }
            />
            <Route
              path="/vehicles/bulk-upload"
              element={
                <ProtectedRoute>
                  <VehicleBulkUpload />
                </ProtectedRoute>
              }
            />
            <Route
              path="/vehicles/:id"
              element={
                <ProtectedRoute>
                  <VehicleProfile />
                </ProtectedRoute>
              }
            />
            <Route
              path="/vehicles/:id/edit"
              element={
                <ProtectedRoute>
                  <VehicleForm />
                </ProtectedRoute>
              }
            />
          </Routes>
        </div>
      </Router>
    </AuthProvider>
  );
}

export default App;
