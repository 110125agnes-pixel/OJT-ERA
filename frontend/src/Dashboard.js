import React, { useState } from 'react';
import './Dashboard.css';
import EmployeeProfiling from './EmployeeProfiling';

function Dashboard({ onLogout }) {
  const [activeModule, setActiveModule] = useState('home');

  const renderContent = () => {
    switch (activeModule) {
      case 'inventory':
        return null;
      case 'profiling':
        return <EmployeeProfiling />;
      case 'home':
      default:
        return (
          <div className="dashboard-home">
            <h2>Welcome to Patient Management Portal</h2>
            <div className="dashboard-cards">
              {/* Inventory card removed */}
              <div className="dashboard-card" onClick={() => setActiveModule('profiling')}>
                <div className="card-icon">👥</div>
                <h3>Patient Management</h3>
                <p>View and manage patient information</p>
              </div>
            </div>
          </div>
        );
    }
  };

  return (
    <div className="dashboard">
      <nav className="dashboard-nav">
        <div className="nav-brand">
          <h1>Patient Portal</h1>
        </div>
        <div className="nav-menu">
          <button 
            className={activeModule === 'home' ? 'nav-item active' : 'nav-item'}
            onClick={() => setActiveModule('home')}
          >
            🏠 Home
          </button>
          {/* Inventory nav removed */}
          <button 
            className={activeModule === 'profiling' ? 'nav-item active' : 'nav-item'}
            onClick={() => setActiveModule('profiling')}
          >
            👥 Patients
          </button>
        </div>
        <button className="logout-btn" onClick={onLogout}>
          Logout
        </button>
      </nav>
      <main className="dashboard-content">
        {renderContent()}
      </main>
    </div>
  );
}

export default Dashboard;
