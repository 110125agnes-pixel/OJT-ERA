import React, { useState, useEffect } from 'react';
import './PhysicalExaminations.css';

const PhysicalExaminations = ({ employeeId }) => {
  const [examData, setExamData] = useState({
    systolicBP: '',
    diastolicBP: '',
    heartRate: '',
    respiratoryRate: '',
    temperature: '',
    height: '',
    weight: '',
    bmi: '',
    pzScore: '',
    leftEyeVision: '',
    rightEyeVision: '',
    pediatricLength: '',
    pediatricHeadCirc: '',
    pediatricSkinfold: '',
    pediatricWaist: '',
    pediatricHip: '',
    pediatricLimbs: '',
    pediatricArmCirc: ''
  });

  const [isEditing, setIsEditing] = useState(false);
  const [error, setError] = useState(null);
  const [successMessage, setSuccessMessage] = useState(null);

  // Calculate BMI when height and weight change
  useEffect(() => {
    if (examData.height && examData.weight) {
      const heightInMeters = examData.height / 100;
      const bmi = (examData.weight / (heightInMeters * heightInMeters)).toFixed(2);
      setExamData(prev => ({
        ...prev,
        bmi: bmi
      }));
    }
  }, [examData.height, examData.weight]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setExamData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSave = () => {
    try {
      setError(null);
      setSuccessMessage('Physical examination data saved successfully!');
      setTimeout(() => setSuccessMessage(null), 3000);
      setIsEditing(false);
    } catch (err) {
      setError('Failed to save examination data');
    }
  };

  const handleCancel = () => {
    setIsEditing(false);
  };

  return (
    <div className="physical-exam-container">
      <div className="exam-header">
        <h2>Pertinent Physical Examinations</h2>
        <button 
          className="edit-btn"
          onClick={() => setIsEditing(!isEditing)}
        >
          {isEditing ? '✕ Cancel' : '✎ Edit'}
        </button>
      </div>

      {error && <div className="error-message">{error}</div>}
      {successMessage && <div className="success-message">{successMessage}</div>}

      <div className="exam-content">
        {/* Left Column */}
        <div className="exam-column">
          <div className="exam-section">
            <h3>Vital Signs</h3>
            
            <div className="exam-field">
              <label>1. Systolic Blood Pressure of Adult Patient -mmHG</label>
              <input
                type="number"
                name="systolicBP"
                placeholder="120"
                value={examData.systolicBP}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>

            <div className="exam-field">
              <label>2. Diastolic Blood Pressure of Adult Patient</label>
              <input
                type="number"
                name="diastolicBP"
                placeholder="80"
                value={examData.diastolicBP}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>

            <div className="exam-field">
              <label>3. Heart Rate of Patient per Minute</label>
              <input
                type="number"
                name="heartRate"
                placeholder="44"
                value={examData.heartRate}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>

            <div className="exam-field">
              <label>4. Respiratory Rate of Patient per Minute</label>
              <input
                type="number"
                name="respiratoryRate"
                placeholder="44"
                value={examData.respiratoryRate}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>

            <div className="exam-field">
              <label>5. Temperature in Celsius</label>
              <input
                type="number"
                name="temperature"
                placeholder="36"
                value={examData.temperature}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>
          </div>

          <div className="exam-section">
            <h3>Body Measurements</h3>

            <div className="exam-field">
              <label>6. Height of Patient in cm</label>
              <input
                type="number"
                name="height"
                placeholder="165"
                value={examData.height}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>

            <div className="exam-field">
              <label>7. Weight of Patient in kg</label>
              <input
                type="number"
                name="weight"
                placeholder="54"
                value={examData.weight}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>

            <div className="exam-field">
              <label>8. Body Mass Index Computation: weight (kg)/ height(m)*height(m)</label>
              <input
                type="number"
                name="bmi"
                placeholder="23"
                value={examData.bmi}
                onChange={handleInputChange}
                disabled={true}
                step="0.01"
              />
            </div>

            <div className="exam-field">
              <label>9. pZScore</label>
              <input
                type="number"
                name="pzScore"
                placeholder="128"
                value={examData.pzScore}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>

            <div className="exam-field">
              <label>10. Left Eye Visual Acuity (Vision) of Patient</label>
              <input
                type="text"
                name="leftEyeVision"
                placeholder="20/20"
                value={examData.leftEyeVision}
                onChange={handleInputChange}
                disabled={!isEditing}
              />
            </div>
          </div>
        </div>

        {/* Right Column */}
        <div className="exam-column">
          <div className="exam-section">
            <h3>Vision & Pediatric (0-24 Months)</h3>

            <div className="exam-field">
              <label>11. Right Eye Visual Acuity (Vision) of Patient</label>
              <input
                type="text"
                name="rightEyeVision"
                placeholder="20/20"
                value={examData.rightEyeVision}
                onChange={handleInputChange}
                disabled={!isEditing}
              />
            </div>

            <div className="exam-field">
              <label>12. Length of Patient in cm - for Pediatric Patient only age 0-24 Months</label>
              <input
                type="number"
                name="pediatricLength"
                placeholder="1"
                value={examData.pediatricLength}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>

            <div className="exam-field">
              <label>13. Head Circumference of Patient in cm - for Pediatric Patient only age 0-24 Months</label>
              <input
                type="number"
                name="pediatricHeadCirc"
                placeholder="2"
                value={examData.pediatricHeadCirc}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>

            <div className="exam-field">
              <label>14. Skinfold Thickness of Patient in cm - for Pediatric Patient only age 0-24 Months</label>
              <input
                type="number"
                name="pediatricSkinfold"
                placeholder="3"
                value={examData.pediatricSkinfold}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>

            <div className="exam-field">
              <label>15. Waist of Patient in cm - for Pediatric Patient only age 0-24 Months</label>
              <input
                type="number"
                name="pediatricWaist"
                placeholder="4"
                value={examData.pediatricWaist}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>

            <div className="exam-field">
              <label>16. Hip of Patient in cm - for Pediatric Patient only age 0-24 Months</label>
              <input
                type="number"
                name="pediatricHip"
                placeholder="5"
                value={examData.pediatricHip}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>

            <div className="exam-field">
              <label>17. Limbs of Patient in cm - for Pediatric Patient only age 0-24 Months</label>
              <input
                type="number"
                name="pediatricLimbs"
                placeholder="6"
                value={examData.pediatricLimbs}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>

            <div className="exam-field">
              <label>18. Middle and Upper Arm Circumference - for Pediatric Patient only age 0-24 Months</label>
              <input
                type="number"
                name="pediatricArmCirc"
                placeholder="7"
                value={examData.pediatricArmCirc}
                onChange={handleInputChange}
                disabled={!isEditing}
                step="0.1"
              />
            </div>
          </div>

          {isEditing && (
            <div className="exam-section">
              <div className="action-buttons">
                <button className="save-btn" onClick={handleSave}>Save</button>
                <button className="cancel-btn" onClick={handleCancel}>Cancel</button>
              </div>
            </div>
          )}
        </div>
      </div>

      <div className="exam-footer">
        <span className="last-updated">Last updated: {new Date().toLocaleDateString('en-US', { month: '2-digit', day: '2-digit', year: 'numeric' })}</span>
      </div>
    </div>
  );
};

export default PhysicalExaminations;
