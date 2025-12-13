import React, { useState, useEffect } from 'react';
import { diseaseService } from './services/api';
import './EmployeeDiseases.css';

function EmployeeDiseases({ employeeID, employeeName }) {
  const [diseases, setDiseases] = useState([]);
  const [employeeDiseases, setEmployeeDiseases] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (employeeID) {
      fetchData();
    }
  }, [employeeID]);

  const fetchData = async () => {
    try {
      setLoading(true);
      const allDiseases = await diseaseService.getAllDiseases();
      const empDiseases = await diseaseService.getEmployeeDiseases(employeeID);
      
      setDiseases(allDiseases || []);
      setEmployeeDiseases(empDiseases || []);
      setError('');
    } catch (err) {
      setError('Failed to fetch diseases: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  const toggleDisease = async (diseaseID) => {
    try {
      setLoading(true);
      const isDiseaseSelected = employeeDiseases.some(d => d.disease_id === diseaseID);
      
      if (isDiseaseSelected) {
        await diseaseService.removeDiseaseFromEmployee(employeeID, diseaseID);
        setEmployeeDiseases(employeeDiseases.filter(d => d.disease_id !== diseaseID));
      } else {
        await diseaseService.addDiseaseToEmployee(employeeID, diseaseID);
        await fetchData();
      }
      setError('');
    } catch (err) {
      setError('Failed to update disease: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  const getSelectedDiseaseIDs = () => {
    return employeeDiseases.map(d => d.disease_id);
  };

  return (
    <div className="employee-diseases-container">
      <h3>Medical Conditions for {employeeName}</h3>
      
      {error && <div className="error">{error}</div>}

      {loading && diseases.length === 0 ? (
        <p>Loading...</p>
      ) : diseases.length === 0 ? (
        <p>No diseases available. Please create diseases first.</p>
      ) : (
        <div className="diseases-checkboxes">
          <div className="checkbox-grid">
            {diseases.map(disease => (
              <label key={disease.id} className="disease-checkbox">
                <input
                  type="checkbox"
                  checked={getSelectedDiseaseIDs().includes(disease.id)}
                  onChange={() => toggleDisease(disease.id)}
                  disabled={loading}
                />
                <span className="checkbox-label">
                  <strong>{disease.name}</strong>
                  <br />
                  <small>Code: {disease.code} | Barcode: {disease.barcode}</small>
                </span>
              </label>
            ))}
          </div>
        </div>
      )}

      {employeeDiseases.length > 0 && (
        <div className="selected-diseases">
          <h4>Selected Medical Conditions:</h4>
          <ul className="disease-list">
            {employeeDiseases.map(disease => (
              <li key={disease.id}>
                <span>✓ {disease.disease_name} (Code: {disease.disease_code})</span>
                <button
                  onClick={() => toggleDisease(disease.disease_id)}
                  className="remove-btn"
                  disabled={loading}
                >
                  Remove
                </button>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}

export default EmployeeDiseases;
