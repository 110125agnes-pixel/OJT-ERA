import React, { useState, useEffect } from 'react';
import './SocialHistory.css';

const SocialHistory = ({ patientId, onClose }) => {
    const [socialHistory, setSocialHistory] = useState({
        is_patient_smoker: 'No',
        cigarette_packs_per_year: 0,
        is_alcohol_drinker: 'No',
        bottles_per_day: 0,
        is_illicit_drug_user: 'No',
        is_sexually_active: 'No'
    });
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState('');
    const [hasExistingData, setHasExistingData] = useState(false);

    useEffect(() => {
        console.log('SocialHistory component - patientId:', patientId);
        if (patientId) {
            fetchSocialHistory();
        } else {
            console.warn('SocialHistory: No patientId provided');
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [patientId]);

    const fetchSocialHistory = async () => {
        try {
            const response = await fetch(`http://localhost:8080/api/patients/${patientId}/social-history`);
            if (response.ok) {
                const data = await response.json();
                setSocialHistory(data);
                setHasExistingData(true);
                setMessage('Existing social history loaded');
                setTimeout(() => setMessage(''), 3000);
            }
        } catch (error) {
            console.log('No existing social history found');
            setHasExistingData(false);
        }
    };

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setSocialHistory(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        
        if (!patientId) {
            setMessage('Error: No patient ID provided');
            console.error('No patientId available');
            return;
        }

        setLoading(true);
        setMessage('');

        console.log('Saving social history for patient:', patientId);
        console.log('Data to save:', socialHistory);

        try {
            const response = await fetch(`http://localhost:8080/api/patients/${patientId}/social-history`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    patient_id: parseInt(patientId),
                    ...socialHistory
                }),
            });

            if (response.ok) {
                const savedData = await response.json();
                console.log('Social history saved successfully:', savedData);
                setHasExistingData(true);
                setMessage('Social history saved successfully!');
                setTimeout(() => setMessage(''), 3000);
            } else {
                const errorData = await response.text();
                console.error('Error response:', errorData);
                setMessage('Error saving social history: ' + errorData);
            }
        } catch (error) {
            console.error('Error saving social history:', error);
            setMessage('Error: ' + error.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="medical-content">
            <form onSubmit={handleSubmit} className="social-history-section">
                <h3>Social History</h3>

                {message && (
                    <div className={`message ${message.includes('Error') ? 'error-message' : 'success-message'}`}>
                        {message}
                    </div>
                )}

                <div className="social-history-questions">
                    {/* Question 1: Is Patient a Smoker */}
                    <div className="question-group">
                        <label className="question-label">1. Is Patient a Smoker</label>
                        <div className="radio-group">
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="is_patient_smoker"
                                    value="Yes"
                                    checked={socialHistory.is_patient_smoker === 'Yes'}
                                    onChange={handleInputChange}
                                />
                                Yes
                            </label>
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="is_patient_smoker"
                                    value="No"
                                    checked={socialHistory.is_patient_smoker === 'No'}
                                    onChange={handleInputChange}
                                />
                                No
                            </label>
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="is_patient_smoker"
                                    value="Quit"
                                    checked={socialHistory.is_patient_smoker === 'Quit'}
                                    onChange={handleInputChange}
                                />
                                Quit
                            </label>
                        </div>
                    </div>

                    {/* Question 2: Number of Cigarette Pack Consumed per Year */}
                    <div className="question-group">
                        <label className="question-label">2. Number of Cigarette Pack Consumed per Year</label>
                        <select
                            name="cigarette_packs_per_year"
                            value={socialHistory.cigarette_packs_per_year}
                            onChange={handleInputChange}
                            className="select-input"
                        >
                            <option value="0">0</option>
                            <option value="1">1</option>
                            <option value="2">2</option>
                            <option value="3">3</option>
                            <option value="4">4</option>
                            <option value="5">5</option>
                            <option value="10">10</option>
                            <option value="20">20</option>
                            <option value="30">30</option>
                            <option value="50">50+</option>
                        </select>
                    </div>

                    {/* Question 3: Is Patient an Alcohol Drinker */}
                    <div className="question-group">
                        <label className="question-label">3. Is Patient an Alcohol Drinker</label>
                        <div className="radio-group">
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="is_alcohol_drinker"
                                    value="Yes"
                                    checked={socialHistory.is_alcohol_drinker === 'Yes'}
                                    onChange={handleInputChange}
                                />
                                Yes
                            </label>
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="is_alcohol_drinker"
                                    value="No"
                                    checked={socialHistory.is_alcohol_drinker === 'No'}
                                    onChange={handleInputChange}
                                />
                                No
                            </label>
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="is_alcohol_drinker"
                                    value="Quit"
                                    checked={socialHistory.is_alcohol_drinker === 'Quit'}
                                    onChange={handleInputChange}
                                />
                                Quit
                            </label>
                        </div>
                    </div>

                    {/* Question 4: Number of bottles consumed per day */}
                    <div className="question-group">
                        <label className="question-label">4. Number of bottles consumed per day</label>
                        <select
                            name="bottles_per_day"
                            value={socialHistory.bottles_per_day}
                            onChange={handleInputChange}
                            className="select-input"
                        >
                            <option value="0">0</option>
                            <option value="1">1</option>
                            <option value="2">2</option>
                            <option value="3">3</option>
                            <option value="4">4</option>
                            <option value="5">5</option>
                            <option value="10">10+</option>
                        </select>
                    </div>

                    {/* Question 5: Is patient an illicit Drug User */}
                    <div className="question-group">
                        <label className="question-label">5. Is patient an illicit Drug User</label>
                        <div className="radio-group">
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="is_illicit_drug_user"
                                    value="Yes"
                                    checked={socialHistory.is_illicit_drug_user === 'Yes'}
                                    onChange={handleInputChange}
                                />
                                Yes
                            </label>
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="is_illicit_drug_user"
                                    value="No"
                                    checked={socialHistory.is_illicit_drug_user === 'No'}
                                    onChange={handleInputChange}
                                />
                                No
                            </label>
                        </div>
                    </div>

                    {/* Question 6: Is patient Sexually Active */}
                    <div className="question-group">
                        <label className="question-label">6. Is patient Sexually Active</label>
                        <div className="radio-group">
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="is_sexually_active"
                                    value="Yes"
                                    checked={socialHistory.is_sexually_active === 'Yes'}
                                    onChange={handleInputChange}
                                />
                                Yes
                            </label>
                            <label className="radio-option">
                                <input
                                    type="radio"
                                    name="is_sexually_active"
                                    value="No"
                                    checked={socialHistory.is_sexually_active === 'No'}
                                    onChange={handleInputChange}
                                />
                                No
                            </label>
                        </div>
                    </div>
                </div>

                {/* Form Actions */}
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
                        {loading ? 'Saving...' : (hasExistingData ? 'Update Social History' : 'Save Social History')}
                    </button>
                </div>
            </form>
        </div>
    );
};

export default SocialHistory;