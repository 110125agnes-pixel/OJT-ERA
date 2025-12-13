import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { itemService } from './services/api';
import './PatientView.css';

function PatientView() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [patient, setPatient] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [activeTab, setActiveTab] = useState('Profiling');
  const [activeSubTab, setActiveSubTab] = useState('Medical');

  // Medical history state
  const [medicalHistory, setMedicalHistory] = useState({
    allergy: false,
    asthma: false,
    cancer: false,
    cerebrovascularDisease: false,
    coronaryArteryDisease: false,
    diabetesMellitus: false,
    emphysema: false,
    epilepsySeizureDisorder: false,
    hepatitis: false,
    hyperlipidemia: false,
    hypertension: false,
    pepticUlcer: false,
    pneumonia: false,
    thyroidDisease: false,
    pulmonaryTuberculosis: false,
    extrapulmonaryTuberculosis: false,
    urinaryTractInfection: false,
    mentalIllness: false,
    others: false
  });

  useEffect(() => {
    fetchPatient();
  }, [id]);

  const fetchPatient = async () => {
    try {
      setLoading(true);
      const data = await itemService.getItem(id);
      setPatient(data);
      setError('');
    } catch (err) {
      setError('Failed to fetch patient: ' + (err.response?.data?.error || err.message));
      console.error('Error fetching patient:', err);
    } finally {
      setLoading(false);
    }
  };

  const calculateAge = (birthdate) => {
    if (!birthdate) return '';
    const birth = new Date(birthdate);
    const now = new Date();
    const years = now.getFullYear() - birth.getFullYear();
    const months = now.getMonth() - birth.getMonth();
    const days = now.getDate() - birth.getDate();
    
    let ageYears = years;
    let ageMonths = months;
    let ageDays = days;
    
    if (ageDays < 0) {
      ageMonths--;
      ageDays += 30;
    }
    if (ageMonths < 0) {
      ageYears--;
      ageMonths += 12;
    }
    
    return `${ageYears}Y${ageMonths}M${ageDays}D`;
  };

  const formatDateTime = (dateTime) => {
    if (!dateTime) return '';
    return new Date(dateTime).toLocaleString('en-US', {
      month: '2-digit',
      day: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: true
    });
  };

  const formatDate = (date) => {
    if (!date) return '';
    return new Date(date).toLocaleDateString('en-US', {
      month: '2-digit',
      day: '2-digit',
      year: 'numeric'
    });
  };

  if (loading) return <div className="loading-container">Loading patient data...</div>;
  if (error) return <div className="error-container">{error}</div>;
  if (!patient) return <div className="error-container">Patient not found</div>;

  return (
    <div className="patient-view-container">
      <div className="patient-header">
        <div className="header-icon">
          <div className="patient-icon">👤</div>
        </div>
        <div className="header-info">
          <div className="header-row">
            <div className="header-field">
              <label>Case No.:</label>
              <span>{patient.caseNo}</span>
            </div>
            <div className="header-field">
              <label>Hospital No.:</label>
              <span>{patient.hospitalNo}</span>
            </div>
            <div className="header-field">
              <label>Sex:</label>
              <span>{patient.sex}</span>
            </div>
          </div>
          <div className="header-row">
            <div className="header-field name-field">
              <label>Name:</label>
              <span>{`${patient.lastname?.toUpperCase() || ''} ${patient.firstname?.toUpperCase() || ''} ${patient.middlename?.toUpperCase() || ''} ${patient.suffix?.toUpperCase() || ''}`.trim()}</span>
            </div>
            <div className="header-field">
              <label>Room:</label>
              <span>{patient.room}</span>
            </div>
            <div className="header-field">
              <label>Height:</label>
              <span>{patient.height} cm</span>
            </div>
          </div>
          <div className="header-row">
            <div className="header-field">
              <label>Birthdate:</label>
              <span>{formatDate(patient.birthdate)}</span>
            </div>
            <div className="header-field">
              <label>Admission:</label>
              <span>{formatDateTime(patient.admissionDate)}</span>
            </div>
            <div className="header-field">
              <label>Weight:</label>
              <span>{patient.weight} kg</span>
            </div>
          </div>
          <div className="header-row">
            <div className="header-field">
              <label>Age:</label>
              <span>{calculateAge(patient.birthdate)}</span>
            </div>
            <div className="header-field">
              <label>Discharged:</label>
              <span>{formatDateTime(patient.dischargeDate)}</span>
            </div>
            <div className="header-field">
              <label>Complaint:</label>
              <span>{patient.complaint}</span>
            </div>
          </div>
        </div>
        <div className="header-close">
          <button onClick={() => navigate('/dashboard')} className="close-button">✕ Close</button>
        </div>
      </div>

      <div className="tabs-container">
        <div className="main-tabs">
          <button 
            className={activeTab === 'Profiling' ? 'tab active' : 'tab'}
            onClick={() => setActiveTab('Profiling')}
          >
            Profiling
          </button>
          <button 
            className={activeTab === 'SOAP' ? 'tab active' : 'tab'}
            onClick={() => setActiveTab('SOAP')}
          >
            SOAP
          </button>
          <button 
            className={activeTab === 'Medicine' ? 'tab active' : 'tab'}
            onClick={() => setActiveTab('Medicine')}
          >
            Medicine
          </button>
        </div>

        {activeTab === 'Profiling' && (
          <div className="sub-tabs">
            {['Medical', 'Surgery', 'Family', 'Immunization', 'Social History', 'Female', '*Pertinent Physical Examinations', '*Physical Examination', 'NCDQANS'].map(tab => (
              <button
                key={tab}
                className={activeSubTab === tab ? 'sub-tab active' : 'sub-tab'}
                onClick={() => setActiveSubTab(tab)}
              >
                {tab}
              </button>
            ))}
          </div>
        )}
      </div>

      <div className="content-area">
        {activeTab === 'Profiling' && activeSubTab === 'Medical' && (
          <div className="medical-content">
            <div className="medical-history-section">
              <h3>Medical History Specifics</h3>
              <div className="medical-checkboxes">
                <div className="checkbox-group">
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.allergy}
                      onChange={(e) => setMedicalHistory({...medicalHistory, allergy: e.target.checked})}
                    />
                    <span>Allergy</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.asthma}
                      onChange={(e) => setMedicalHistory({...medicalHistory, asthma: e.target.checked})}
                    />
                    <span>Asthma</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.cancer}
                      onChange={(e) => setMedicalHistory({...medicalHistory, cancer: e.target.checked})}
                    />
                    <span>Cancer</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.cerebrovascularDisease}
                      onChange={(e) => setMedicalHistory({...medicalHistory, cerebrovascularDisease: e.target.checked})}
                    />
                    <span>Cerebrovascular Disease</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.coronaryArteryDisease}
                      onChange={(e) => setMedicalHistory({...medicalHistory, coronaryArteryDisease: e.target.checked})}
                    />
                    <span>Coronary Artery Disease</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.diabetesMellitus}
                      onChange={(e) => setMedicalHistory({...medicalHistory, diabetesMellitus: e.target.checked})}
                    />
                    <span>Diabetes Mellitus</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.emphysema}
                      onChange={(e) => setMedicalHistory({...medicalHistory, emphysema: e.target.checked})}
                    />
                    <span>Emphysema</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.epilepsySeizureDisorder}
                      onChange={(e) => setMedicalHistory({...medicalHistory, epilepsySeizureDisorder: e.target.checked})}
                    />
                    <span>Epilepsy/Seizure Disorder</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.hepatitis}
                      onChange={(e) => setMedicalHistory({...medicalHistory, hepatitis: e.target.checked})}
                    />
                    <span>Hepatitis</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.hyperlipidemia}
                      onChange={(e) => setMedicalHistory({...medicalHistory, hyperlipidemia: e.target.checked})}
                    />
                    <span>Hyperlipidemia</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.hypertension}
                      onChange={(e) => setMedicalHistory({...medicalHistory, hypertension: e.target.checked})}
                    />
                    <span>Hypertension</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.pepticUlcer}
                      onChange={(e) => setMedicalHistory({...medicalHistory, pepticUlcer: e.target.checked})}
                    />
                    <span>Peptic Ulcer</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.pneumonia}
                      onChange={(e) => setMedicalHistory({...medicalHistory, pneumonia: e.target.checked})}
                    />
                    <span>Pneumonia</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.thyroidDisease}
                      onChange={(e) => setMedicalHistory({...medicalHistory, thyroidDisease: e.target.checked})}
                    />
                    <span>Thyroid Disease</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.pulmonaryTuberculosis}
                      onChange={(e) => setMedicalHistory({...medicalHistory, pulmonaryTuberculosis: e.target.checked})}
                    />
                    <span>Pulmonary Tuberculosis</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.extrapulmonaryTuberculosis}
                      onChange={(e) => setMedicalHistory({...medicalHistory, extrapulmonaryTuberculosis: e.target.checked})}
                    />
                    <span>Extrapulmonary Tuberculosis</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.urinaryTractInfection}
                      onChange={(e) => setMedicalHistory({...medicalHistory, urinaryTractInfection: e.target.checked})}
                    />
                    <span>Urinary Tract Infection</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.mentalIllness}
                      onChange={(e) => setMedicalHistory({...medicalHistory, mentalIllness: e.target.checked})}
                    />
                    <span>Mental Illness</span>
                  </label>
                  <label>
                    <input
                      type="checkbox"
                      checked={medicalHistory.others}
                      onChange={(e) => setMedicalHistory({...medicalHistory, others: e.target.checked})}
                    />
                    <span>Others</span>
                  </label>
                </div>
                <div className="medical-details">
                  <textarea
                    placeholder="Enter medical history details here..."
                    rows="10"
                  ></textarea>
                </div>
              </div>

              <div className="medical-codes-table">
                <table>
                  <thead>
                    <tr>
                      <th>Code</th>
                      <th>Description</th>
                    </tr>
                  </thead>
                  <tbody>
                    {medicalHistory.emphysema && (
                      <tr>
                        <td>007</td>
                        <td>Emphysema</td>
                      </tr>
                    )}
                    {medicalHistory.pepticUlcer && (
                      <tr>
                        <td>012</td>
                        <td>Peptic Ulcer</td>
                      </tr>
                    )}
                    {medicalHistory.thyroidDisease && (
                      <tr>
                        <td>014</td>
                        <td>Thyroid Disease</td>
                      </tr>
                    )}
                    {medicalHistory.allergy && (
                      <tr>
                        <td>001</td>
                        <td>Allergy</td>
                      </tr>
                    )}
                    {medicalHistory.cancer && (
                      <tr>
                        <td>003</td>
                        <td>Cancer</td>
                      </tr>
                    )}
                    {medicalHistory.asthma && (
                      <tr>
                        <td>002</td>
                        <td>Asthma</td>
                      </tr>
                    )}
                    {medicalHistory.diabetesMellitus && (
                      <tr>
                        <td>006</td>
                        <td>Diabetes Mellitus</td>
                      </tr>
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'SOAP' && (
          <div className="soap-content">
            <h3>SOAP Notes</h3>
            <p>SOAP notes content will be displayed here.</p>
          </div>
        )}

        {activeTab === 'Medicine' && (
          <div className="medicine-content">
            <h3>Medicine Records</h3>
            <p>Medicine records will be displayed here.</p>
          </div>
        )}
      </div>

      <div className="footer-info">
        <span>{new Date().toLocaleDateString('en-US')}</span>
        <span>{new Date().toLocaleTimeString('en-US')}</span>
      </div>
    </div>
  );
}

export default PatientView;
