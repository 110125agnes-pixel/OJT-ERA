import React, { useState } from 'react';
import './Dashboard.css';
import Inventory from './Inventory';
import EmployeeProfiling from './EmployeeProfiling';
import DiseaseManagement from './DiseaseManagement';
import PhysicalExaminations from './PhysicalExaminations';

function Dashboard({ onLogout }) {
  const [activeModule, setActiveModule] = useState('home');

  const renderContent = () => {
    switch (activeModule) {
      case 'inventory':
        return <Inventory />;
      case 'profiling':
        return <EmployeeProfiling />;
      case 'diseases':
        return <DiseaseManagement />;
      case 'examinations':
        return <PhysicalExaminations />;
      case 'home':
      default:
        return (
          <div className="dashboard-home">
            <h2>Welcome to Employee Portal</h2>
            <div className="dashboard-cards">
              <div className="dashboard-card" onClick={() => setActiveModule('profiling')}>
                <div className="card-icon">👥</div>
                <h3>Employee Profiling</h3>
                <p>View and manage employee information</p>
              </div>
              <div className="dashboard-card" onClick={() => setActiveModule('diseases')}>
                <div className="card-icon">🏥</div>
                <h3>Disease Management</h3>
                <p>Manage medical conditions and diseases</p>
              </div>
              <div className="dashboard-card" onClick={() => setActiveModule('examinations')}>
                <div className="card-icon">🩺</div>
                <h3>Physical Examinations</h3>
                <p>Record and track patient vital signs and measurements</p>
              </div>
              <div className="dashboard-card" onClick={() => setActiveModule('inventory')}>
                <div className="card-icon">📦</div>
                <h3>Inventory Module</h3>
                <p>Manage items, stock counts, and inventory operations</p>
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
          <h1>Employee Portal</h1>
        </div>
        <div className="nav-menu">
          <button 
            className={activeModule === 'home' ? 'nav-item active' : 'nav-item'}
            onClick={() => setActiveModule('home')}
          >
            🏠 Home
          </button>
          <button 
            className={activeModule === 'profiling' ? 'nav-item active' : 'nav-item'}
            onClick={() => setActiveModule('profiling')}
          >
            👥 Employees
          </button>
          <button 
            className={activeModule === 'diseases' ? 'nav-item active' : 'nav-item'}
            onClick={() => setActiveModule('diseases')}
          >
            🏥 Diseases
          </button>
          <button 
            className={activeModule === 'examinations' ? 'nav-item active' : 'nav-item'}
            onClick={() => setActiveModule('examinations')}
          >
            🩺 Examinations
          </button>
          <button 
            className={activeModule === 'inventory' ? 'nav-item active' : 'nav-item'}
            onClick={() => setActiveModule('inventory')}
          >
            📦 Inventory
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
