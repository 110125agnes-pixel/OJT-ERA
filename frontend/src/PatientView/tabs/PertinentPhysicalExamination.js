import React, { useState, useEffect } from 'react';
import './PertinentPhysicalExamination.css';

const PertinentPhysicalExamination = ({ patientId }) => {
    const [formData, setFormData] = useState({
        systolic_bp: '',
        diastolic_bp: '',
        heart_rate: '',
        respiratory_rate: '',
        temperature: '',
        height: '',
        weight: '',
        bmi: '',
        pzscore: '',
        left_eye_vision: '',
        right_eye_vision: '',
        length_pediatric: '',
        head_circumference: '',
        skinfold_thickness: '',
        waist: '',
        hip: '',
        limbs: '',
        arm_circumference: ''
    });
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState('');
    const [hasExistingData, setHasExistingData] = useState(false);

    useEffect(() => {
        if (patientId) {
            fetchPertinentPhysicalExam();
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [patientId]);

    // Auto-calculate BMI when height and weight change
    useEffect(() => {
        if (formData.height && formData.weight) {
            const heightInMeters = parseFloat(formData.height) / 100;
            const weightInKg = parseFloat(formData.weight);
            if (heightInMeters > 0 && weightInKg > 0) {
                const bmi = (weightInKg / (heightInMeters * heightInMeters)).toFixed(2);
                setFormData(prev => ({ ...prev, bmi }));
            }
        }
    }, [formData.height, formData.weight]);

    const fetchPertinentPhysicalExam = async () => {
        try {
            const response = await fetch(`http://localhost:8080/api/patients/${patientId}/pertinent-physical-exam`);
            if (response.ok) {
                const data = await response.json();
                setFormData(data);
                setHasExistingData(true);
                setMessage('Existing data loaded');
                setTimeout(() => setMessage(''), 3000);
            }
        } catch (error) {
            console.log('No existing pertinent physical examination found');
            setHasExistingData(false);
        }
    };

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        
        if (!patientId) {
            setMessage('Error: No patient ID provided');
            return;
        }

        setLoading(true);
        setMessage('');

        try {
            const response = await fetch(`http://localhost:8080/api/patients/${patientId}/pertinent-physical-exam`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    patient_id: parseInt(patientId),
                    ...formData
                }),
            });

            if (response.ok) {
                const savedData = await response.json();
                setHasExistingData(true);
                setMessage('Pertinent physical examination saved successfully!');
                setTimeout(() => setMessage(''), 3000);
            } else {
                const errorData = await response.text();
                setMessage('Error saving data: ' + errorData);
            }
        } catch (error) {
            setMessage('Error: ' + error.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="pertinent-exam-container">
            <form onSubmit={handleSubmit} className="pertinent-exam-form">
                <h3>Pertinent Physical Examinations</h3>

                <div className="form-grid">
                    <div className="form-row">
                        <div className="form-field">
                            <label>1. Systolic Blood Pressure of Adult Patient -mmHG</label>
                            <input
                                type="number"
                                name="systolic_bp"
                                value={formData.systolic_bp}
                                onChange={handleInputChange}
                                placeholder="120"
                            />
                        </div>
                        <div className="form-field">
                            <label>11. Right Eye Visual Acuity (Vision) of Patient</label>
                            <input
                                type="text"
                                name="right_eye_vision"
                                value={formData.right_eye_vision}
                                onChange={handleInputChange}
                                placeholder="20/20"
                            />
                        </div>
                    </div>

                    <div className="form-row">
                        <div className="form-field">
                            <label>2. Diastolic Blood Pressure of Adult Patient</label>
                            <input
                                type="number"
                                name="diastolic_bp"
                                value={formData.diastolic_bp}
                                onChange={handleInputChange}
                                placeholder="80"
                            />
                        </div>
                        <div className="form-field">
                            <label>12. Length of Patient in cm - for Pediatric Patient only age 0-24 Months</label>
                            <input
                                type="number"
                                name="length_pediatric"
                                value={formData.length_pediatric}
                                onChange={handleInputChange}
                                placeholder="1"
                            />
                        </div>
                    </div>

                    <div className="form-row">
                        <div className="form-field">
                            <label>3. Heart Rate of Patient per Minute</label>
                            <input
                                type="number"
                                name="heart_rate"
                                value={formData.heart_rate}
                                onChange={handleInputChange}
                                placeholder="44"
                            />
                        </div>
                        <div className="form-field">
                            <label>13. Head Circumference of Patient in cm - for Pediatric Patient only age 0-24 Months</label>
                            <input
                                type="number"
                                name="head_circumference"
                                value={formData.head_circumference}
                                onChange={handleInputChange}
                                placeholder="2"
                            />
                        </div>
                    </div>

                    <div className="form-row">
                        <div className="form-field">
                            <label>4. Respiratory Rate of Patient per Minute</label>
                            <input
                                type="number"
                                name="respiratory_rate"
                                value={formData.respiratory_rate}
                                onChange={handleInputChange}
                                placeholder="44"
                            />
                        </div>
                        <div className="form-field">
                            <label>14. Skinfold Thickness of Patient in cm - for Pediatric Patient only age 0-24 Months</label>
                            <input
                                type="number"
                                name="skinfold_thickness"
                                value={formData.skinfold_thickness}
                                onChange={handleInputChange}
                                placeholder="3"
                            />
                        </div>
                    </div>

                    <div className="form-row">
                        <div className="form-field">
                            <label>5. Temperature in Celsius</label>
                            <input
                                type="number"
                                step="0.1"
                                name="temperature"
                                value={formData.temperature}
                                onChange={handleInputChange}
                                placeholder="36"
                            />
                        </div>
                        <div className="form-field">
                            <label>15. Waist of Patient in cm - for Pediatric Patient only age 0-24 Months</label>
                            <input
                                type="number"
                                name="waist"
                                value={formData.waist}
                                onChange={handleInputChange}
                                placeholder="4"
                            />
                        </div>
                    </div>

                    <div className="form-row">
                        <div className="form-field">
                            <label>6. Height of Patient in cm</label>
                            <input
                                type="number"
                                name="height"
                                value={formData.height}
                                onChange={handleInputChange}
                                placeholder="165"
                            />
                        </div>
                        <div className="form-field">
                            <label>16. Hip of Patient in cm - for Pediatric Patient only age 0-24 Months</label>
                            <input
                                type="number"
                                name="hip"
                                value={formData.hip}
                                onChange={handleInputChange}
                                placeholder="5"
                            />
                        </div>
                    </div>

                    <div className="form-row">
                        <div className="form-field">
                            <label>7. Weight of Patient in kg</label>
                            <input
                                type="number"
                                name="weight"
                                value={formData.weight}
                                onChange={handleInputChange}
                                placeholder="54"
                            />
                        </div>
                        <div className="form-field">
                            <label>17. Limbs of Patient in cm - for Pediatric Patient only age 0-24 Months</label>
                            <input
                                type="number"
                                name="limbs"
                                value={formData.limbs}
                                onChange={handleInputChange}
                                placeholder="6"
                            />
                        </div>
                    </div>

                    <div className="form-row">
                        <div className="form-field">
                            <label>8. Body Mass Index Computation: weight (kg)/ height(m)*height(m)</label>
                            <input
                                type="number"
                                step="0.01"
                                name="bmi"
                                value={formData.bmi}
                                onChange={handleInputChange}
                                placeholder="23"
                                readOnly
                                className="readonly-field"
                            />
                        </div>
                        <div className="form-field">
                            <label>18. Middle and Upper Arm Circumference - for Pediatric Patient only age 0-24 Months</label>
                            <input
                                type="number"
                                name="arm_circumference"
                                value={formData.arm_circumference}
                                onChange={handleInputChange}
                                placeholder="7"
                            />
                        </div>
                    </div>

                    <div className="form-row">
                        <div className="form-field">
                            <label>9. pZScore</label>
                            <input
                                type="number"
                                name="pzscore"
                                value={formData.pzscore}
                                onChange={handleInputChange}
                                placeholder="128"
                            />
                        </div>
                    </div>

                    <div className="form-row">
                        <div className="form-field full-width">
                            <label>10. Left Eye Visual Acuity (Vision) of Patient</label>
                            <input
                                type="text"
                                name="left_eye_vision"
                                value={formData.left_eye_vision}
                                onChange={handleInputChange}
                                placeholder="20/20"
                            />
                        </div>
                    </div>
                </div>

                <div className="form-actions">
                    {hasExistingData && (
                        <span className="existing-data-indicator">
                            ✓ Has existing data
                        </span>
                    )}
                    <button 
                        type="submit" 
                        className="save-button"
                        disabled={loading}
                    >
                        {loading ? 'Saving...' : (hasExistingData ? 'Update Data' : 'Save Data')}
                    </button>
                </div>

                {message && (
                    <div className={`message ${message.includes('Error') ? 'error-message' : 'success-message'}`}>
                        {message}
                    </div>
                )}
            </form>
        </div>
    );
};

export default PertinentPhysicalExamination;
