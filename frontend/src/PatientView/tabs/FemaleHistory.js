import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import './FemaleHistory.css';

// Default state — all empty/false so DB values always win on load.
const DEFAULT_STATE = {
  ageOfFirstMenstruation: '',
  dateOfLastMenstrualPeriod: '',
  durationOfMenstrualPeriod: '',
  intervalCycleOfMenstruation: '',
  numberOfPadsPerDay: '',
  onsetOfSexualIntercourse: '',
  birthControlMethod: '',
  isMenopause: false,
  ageOfMenopause: '',
  isMenstrualHistoryApplicable: false,
  numberOfPregnancyToDate: '',
  numberOfDeliveryToDate: '',
  typeOfDelivery: '',
  numberOfFullTermPregnancy: '',
  numberOfPrematurePregnancy: '',
  numberOfAbortion: '',
  numberOfLivingChildren: '',
  pregnancyInducedHypertension: false,
  accessToFamilyPlanningCounselling: false,
  isPregnancyHistoryApplicable: false,
};

function FemaleHistory({ patientId }) {
  const [femaleHistory, setFemaleHistory] = useState(DEFAULT_STATE);

  // fetchFemaleHistory — loads saved values from backend (source of truth).
  // Called on mount and whenever patientId changes, same pattern as MedicalHistory.
  const fetchFemaleHistory = useCallback(async () => {
    if (!patientId) return;
    try {
      const res = await axios.get(`/api/patients/${patientId}/female-history`);
      if (res.data && Object.keys(res.data).length > 0) {
        // Replace state fully with server values; server always returns all keys.
        setFemaleHistory({ ...DEFAULT_STATE, ...res.data });
      }
    } catch (err) {
      console.error('Failed to load female history', err);
    }
  }, [patientId]);

  // On mount / patientId change — load saved data from DB, same as MedicalHistory.
  useEffect(() => {
    fetchFemaleHistory();
  }, [fetchFemaleHistory]);

  // saveFemaleHistory — POSTs current state to backend and merges server response
  // back into local state so the UI always reflects what was actually stored in DB.
  const saveFemaleHistory = useCallback(async (next) => {
    if (!patientId) return;
    try {
      const res = await axios.post(`/api/patients/${patientId}/female-history`, next);
      if (res.data && Object.keys(res.data).length > 0) {
        setFemaleHistory({ ...DEFAULT_STATE, ...res.data });
      }
    } catch (err) {
      console.error('Failed to save female history', err);
    }
  }, [patientId]);

  // handleChange — updates local state and immediately persists to DB,
  // same immediate-save pattern as MedicalHistory's handleToggle.
  const handleChange = useCallback((key, value) => {
    setFemaleHistory((prev) => {
      const next = { ...prev, [key]: value };
      saveFemaleHistory(next);
      return next;
    });
  }, [saveFemaleHistory]);

  return (
    <div className="female-content">
      <div className="female-history-sections">
        <div className="menstrual-history-section">
          <h3>Menstrual History</h3>
          <div className="form-grid">
            <div className="form-field">
              <label>1. Age of First Menstruation (Menarche)</label>
              <select 
                value={femaleHistory.ageOfFirstMenstruation}
                onChange={(e) => handleChange('ageOfFirstMenstruation', e.target.value)}
              >
                <option value="">Select</option>
                {Array.from({length: 20}, (_, i) => i + 8).map(age => (
                  <option key={age} value={age}>{age}</option>
                ))}
              </select>
            </div>
            
            <div className="form-field">
              <label>2. Date of Last Menstrual Period</label>
              <input 
                type="date"
                value={femaleHistory.dateOfLastMenstrualPeriod}
                onChange={(e) => handleChange('dateOfLastMenstrualPeriod', e.target.value)}
              />
            </div>

            <div className="form-field">
              <label>3. Duration of Menstrual Period in Number of Days</label>
              <select 
                value={femaleHistory.durationOfMenstrualPeriod}
                onChange={(e) => handleChange('durationOfMenstrualPeriod', e.target.value)}
              >
                <option value="">Select</option>
                {Array.from({length: 15}, (_, i) => i + 1).map(days => (
                  <option key={days} value={days}>{days}</option>
                ))}
              </select>
            </div>

            <div className="form-field">
              <label>4. Interval/Cycle of Menstruation in Number of Days</label>
              <select 
                value={femaleHistory.intervalCycleOfMenstruation}
                onChange={(e) => handleChange('intervalCycleOfMenstruation', e.target.value)}
              >
                <option value="">Select</option>
                {Array.from({length: 50}, (_, i) => i + 20).map(days => (
                  <option key={days} value={days}>{days}</option>
                ))}
              </select>
            </div>

            <div className="form-field">
              <label>5. Number of Pads/Napkins Used per Day during Menstruation</label>
              <select 
                value={femaleHistory.numberOfPadsPerDay}
                onChange={(e) => handleChange('numberOfPadsPerDay', e.target.value)}
              >
                <option value="">Select</option>
                {Array.from({length: 15}, (_, i) => i + 1).map(num => (
                  <option key={num} value={num}>{num}</option>
                ))}
              </select>
            </div>

            <div className="form-field">
              <label>6. Onset of Sexual Intercourse (Age of First Sexual Intercourse)</label>
              <select 
                value={femaleHistory.onsetOfSexualIntercourse}
                onChange={(e) => handleChange('onsetOfSexualIntercourse', e.target.value)}
              >
                <option value="">Select</option>
                {Array.from({length: 40}, (_, i) => i + 10).map(age => (
                  <option key={age} value={age}>{age}</option>
                ))}
              </select>
            </div>

            <div className="form-field">
              <label>7. Birth Control Method Used</label>
              <input 
                type="text"
                value={femaleHistory.birthControlMethod}
                onChange={(e) => handleChange('birthControlMethod', e.target.value)}
                placeholder="ORA"
              />
            </div>

            <div className="form-field">
              <label>8. Is Menopause?</label>
              <div className="radio-group">
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.isMenopause === true}
                    onChange={() => handleChange('isMenopause', true)}
                  />
                  <span>Yes</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.isMenopause === false}
                    onChange={() => handleChange('isMenopause', false)}
                  />
                  <span>No</span>
                </label>
              </div>
            </div>

            <div className="form-field">
              <label>9. If Menopause, Age of Menopause</label>
              <select 
                value={femaleHistory.ageOfMenopause}
                onChange={(e) => handleChange('ageOfMenopause', e.target.value)}
                disabled={!femaleHistory.isMenopause}
              >
                <option value="">Select</option>
                {Array.from({length: 30}, (_, i) => i + 35).map(age => (
                  <option key={age} value={age}>{age}</option>
                ))}
              </select>
            </div>

            <div className="form-field">
              <label>10. Is menstrual history applicable?</label>
              <div className="radio-group">
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.isMenstrualHistoryApplicable === true}
                    onChange={() => handleChange('isMenstrualHistoryApplicable', true)}
                  />
                  <span>Yes</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.isMenstrualHistoryApplicable === false}
                    onChange={() => handleChange('isMenstrualHistoryApplicable', false)}
                  />
                  <span>No</span>
                </label>
              </div>
            </div>
          </div>
        </div>

        <div className="pregnancy-history-section">
          <h3>Pregnancy History</h3>
          <div className="form-grid">
            <div className="form-field">
              <label>1. Number of Pregnancy to Date - Gravity Chief</label>
              <select 
                value={femaleHistory.numberOfPregnancyToDate}
                onChange={(e) => handleChange('numberOfPregnancyToDate', e.target.value)}
              >
                <option value="">Select</option>
                {Array.from({length: 20}, (_, i) => i).map(num => (
                  <option key={num} value={num}>{num}</option>
                ))}
              </select>
            </div>

            <div className="form-field">
              <label>2. Number of Delivery to Date - Parity</label>
              <select 
                value={femaleHistory.numberOfDeliveryToDate}
                onChange={(e) => handleChange('numberOfDeliveryToDate', e.target.value)}
              >
                <option value="">Select</option>
                {Array.from({length: 20}, (_, i) => i).map(num => (
                  <option key={num} value={num}>{num}</option>
                ))}
              </select>
            </div>

            <div className="form-field">
              <label>3. Type of Delivery</label>
              <div className="radio-group">
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.typeOfDelivery === 'Normal'}
                    onChange={() => handleChange('typeOfDelivery', 'Normal')}
                  />
                  <span>Normal</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.typeOfDelivery === 'Operative'}
                    onChange={() => handleChange('typeOfDelivery', 'Operative')}
                  />
                  <span>Operative</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.typeOfDelivery === 'Both'}
                    onChange={() => handleChange('typeOfDelivery', 'Both')}
                  />
                  <span>Both Normal and Operative</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.typeOfDelivery === 'NotApplicable'}
                    onChange={() => handleChange('typeOfDelivery', 'NotApplicable')}
                  />
                  <span>Not Applicable</span>
                </label>
              </div>
            </div>

            <div className="form-field">
              <label>4. Number of Full Term Pregnancy</label>
              <select 
                value={femaleHistory.numberOfFullTermPregnancy}
                onChange={(e) => handleChange('numberOfFullTermPregnancy', e.target.value)}
              >
                <option value="">Select</option>
                {Array.from({length: 20}, (_, i) => i).map(num => (
                  <option key={num} value={num}>{num}</option>
                ))}
              </select>
            </div>

            <div className="form-field">
              <label>5. Number of Premature Pregnancy</label>
              <select 
                value={femaleHistory.numberOfPrematurePregnancy}
                onChange={(e) => handleChange('numberOfPrematurePregnancy', e.target.value)}
              >
                <option value="">Select</option>
                {Array.from({length: 20}, (_, i) => i).map(num => (
                  <option key={num} value={num}>{num}</option>
                ))}
              </select>
            </div>

            <div className="form-field">
              <label>6. Number of Abortion</label>
              <select 
                value={femaleHistory.numberOfAbortion}
                onChange={(e) => handleChange('numberOfAbortion', e.target.value)}
              >
                <option value="">Select</option>
                {Array.from({length: 20}, (_, i) => i).map(num => (
                  <option key={num} value={num}>{num}</option>
                ))}
              </select>
            </div>

            <div className="form-field">
              <label>7. Number of Living Children</label>
              <select 
                value={femaleHistory.numberOfLivingChildren}
                onChange={(e) => handleChange('numberOfLivingChildren', e.target.value)}
              >
                <option value="">Select</option>
                {Array.from({length: 20}, (_, i) => i).map(num => (
                  <option key={num} value={num}>{num}</option>
                ))}
              </select>
            </div>

            <div className="form-field">
              <label>8. If Pregnancy - Induced Hypertension (Pre - Eclampsia)</label>
              <div className="radio-group">
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.pregnancyInducedHypertension === true}
                    onChange={() => handleChange('pregnancyInducedHypertension', true)}
                  />
                  <span>Yes</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.pregnancyInducedHypertension === false}
                    onChange={() => handleChange('pregnancyInducedHypertension', false)}
                  />
                  <span>No</span>
                </label>
              </div>
            </div>

            <div className="form-field">
              <label>9. If with access to Family Planning Counselling</label>
              <div className="radio-group">
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.accessToFamilyPlanningCounselling === true}
                    onChange={() => handleChange('accessToFamilyPlanningCounselling', true)}
                  />
                  <span>Yes</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.accessToFamilyPlanningCounselling === false}
                    onChange={() => handleChange('accessToFamilyPlanningCounselling', false)}
                  />
                  <span>No</span>
                </label>
              </div>
            </div>

            <div className="form-field">
              <label>10. Is pregnancy history applicable?</label>
              <div className="radio-group">
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.isPregnancyHistoryApplicable === true}
                    onChange={() => handleChange('isPregnancyHistoryApplicable', true)}
                  />
                  <span>Yes</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.isPregnancyHistoryApplicable === false}
                    onChange={() => handleChange('isPregnancyHistoryApplicable', false)}
                  />
                  <span>No</span>
                </label>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default FemaleHistory;