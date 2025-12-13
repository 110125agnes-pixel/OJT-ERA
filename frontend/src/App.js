import React, { useState } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import './App.css';
import Login from './Login';
import Dashboard from './Dashboard';
import Immunization from './Immunization';
import FemaleHistory from './FemaleHistory';

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  const handleLogin = () => {
    setIsLoggedIn(true);
  };

  const handleLogout = () => {
    setIsLoggedIn(false);
  };

  return (
    <BrowserRouter>
      <Routes>
        {/* Standalone route to Immunization panel (frontend-only) */}
        <Route path="/immunization" element={<Immunization />} />
        {/* Standalone route to Female history panel */}
        <Route path="/female" element={<FemaleHistory />} />

        {/* Default app flow: login -> dashboard */}
        <Route
          path="/"
          element={
            !isLoggedIn ? (
              <Login onLogin={handleLogin} />
            ) : (
              <Dashboard onLogout={handleLogout} />
            )
          }
        />
        {/* Fallback to root */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
