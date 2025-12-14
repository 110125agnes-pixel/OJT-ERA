import React, { useState } from 'react';
import './FemaleHistory.css';

function FemaleHistory() {
  const [femaleHistory, setFemaleHistory] = useState({
    // Menstrual History
    ageOfFirstMenstruation: '12',
    dateOfLastMenstrualPeriod: '',
    durationOfMenstrualPeriod: '7',
    intervalCycleOfMenstruation: '30',
    numberOfPadsPerDay: '4',
    onsetOfSexualIntercourse: '18',
    birthControlMethod: '',
    isMenopause: false,
    ageOfMenopause: '40',
    isMenstrualHistoryApplicable: false,
    // Pregnancy History
    numberOfPregnancyToDate: '4',
    numberOfDeliveryToDate: '2',
    typeOfDelivery: 'Normal',
    numberOfFullTermPregnancy: '',
    numberOfPrematurePregnancy: '',
    numberOfAbortion: '',
    numberOfLivingChildren: '3',
    pregnancyInducedHypertension: true,
    accessToFamilyPlanningCounselling: true,
    isPregnancyHistoryApplicable: false
  });

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
                onChange={(e) => setFemaleHistory({...femaleHistory, ageOfFirstMenstruation: e.target.value})}
              >
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
                onChange={(e) => setFemaleHistory({...femaleHistory, dateOfLastMenstrualPeriod: e.target.value})}
              />
            </div>

            <div className="form-field">
              <label>3. Duration of Menstrual Period in Number of Days</label>
              <select 
                value={femaleHistory.durationOfMenstrualPeriod}
                onChange={(e) => setFemaleHistory({...femaleHistory, durationOfMenstrualPeriod: e.target.value})}
              >
                {Array.from({length: 15}, (_, i) => i + 1).map(days => (
                  <option key={days} value={days}>{days}</option>
                ))}
              </select>
            </div>

            <div className="form-field">
              <label>4. Interval/Cycle of Menstruation in Number of Days</label>
              <select 
                value={femaleHistory.intervalCycleOfMenstruation}
                onChange={(e) => setFemaleHistory({...femaleHistory, intervalCycleOfMenstruation: e.target.value})}
              >
                {Array.from({length: 50}, (_, i) => i + 20).map(days => (
                  <option key={days} value={days}>{days}</option>
                ))}
              </select>
            </div>

            <div className="form-field">
              <label>5. Number of Pads/Napkins Used per Day during Menstruation</label>
              <select 
                value={femaleHistory.numberOfPadsPerDay}
                onChange={(e) => setFemaleHistory({...femaleHistory, numberOfPadsPerDay: e.target.value})}
              >
                {Array.from({length: 15}, (_, i) => i + 1).map(num => (
                  <option key={num} value={num}>{num}</option>
                ))}
              </select>
            </div>

            <div className="form-field">
              <label>6. Onset of Sexual Intercourse (Age of First Sexual Intercourse)</label>
              <select 
                value={femaleHistory.onsetOfSexualIntercourse}
                onChange={(e) => setFemaleHistory({...femaleHistory, onsetOfSexualIntercourse: e.target.value})}
              >
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
                onChange={(e) => setFemaleHistory({...femaleHistory, birthControlMethod: e.target.value})}
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
                    onChange={() => setFemaleHistory({...femaleHistory, isMenopause: true})}
                  />
                  <span>Yes</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.isMenopause === false}
                    onChange={() => setFemaleHistory({...femaleHistory, isMenopause: false})}
                  />
                  <span>No</span>
                </label>
              </div>
            </div>

            <div className="form-field">
              <label>9. If Menopause, Age of Menopause</label>
              <select 
                value={femaleHistory.ageOfMenopause}
                onChange={(e) => setFemaleHistory({...femaleHistory, ageOfMenopause: e.target.value})}
                disabled={!femaleHistory.isMenopause}
              >
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
                    onChange={() => setFemaleHistory({...femaleHistory, isMenstrualHistoryApplicable: true})}
                  />
                  <span>Yes</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.isMenstrualHistoryApplicable === false}
                    onChange={() => setFemaleHistory({...femaleHistory, isMenstrualHistoryApplicable: false})}
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
                onChange={(e) => setFemaleHistory({...femaleHistory, numberOfPregnancyToDate: e.target.value})}
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
                onChange={(e) => setFemaleHistory({...femaleHistory, numberOfDeliveryToDate: e.target.value})}
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
                    onChange={() => setFemaleHistory({...femaleHistory, typeOfDelivery: 'Normal'})}
                  />
                  <span>Normal</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.typeOfDelivery === 'Operative'}
                    onChange={() => setFemaleHistory({...femaleHistory, typeOfDelivery: 'Operative'})}
                  />
                  <span>Operative</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.typeOfDelivery === 'Both'}
                    onChange={() => setFemaleHistory({...femaleHistory, typeOfDelivery: 'Both'})}
                  />
                  <span>Both Normal and Operative</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.typeOfDelivery === 'NotApplicable'}
                    onChange={() => setFemaleHistory({...femaleHistory, typeOfDelivery: 'NotApplicable'})}
                  />
                  <span>Not Applicable</span>
                </label>
              </div>
            </div>

            <div className="form-field">
              <label>4. Number of Full Term Pregnancy</label>
              <select 
                value={femaleHistory.numberOfFullTermPregnancy}
                onChange={(e) => setFemaleHistory({...femaleHistory, numberOfFullTermPregnancy: e.target.value})}
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
                onChange={(e) => setFemaleHistory({...femaleHistory, numberOfPrematurePregnancy: e.target.value})}
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
                onChange={(e) => setFemaleHistory({...femaleHistory, numberOfAbortion: e.target.value})}
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
                onChange={(e) => setFemaleHistory({...femaleHistory, numberOfLivingChildren: e.target.value})}
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
                    onChange={() => setFemaleHistory({...femaleHistory, pregnancyInducedHypertension: true})}
                  />
                  <span>Yes</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.pregnancyInducedHypertension === false}
                    onChange={() => setFemaleHistory({...femaleHistory, pregnancyInducedHypertension: false})}
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
                    onChange={() => setFemaleHistory({...femaleHistory, accessToFamilyPlanningCounselling: true})}
                  />
                  <span>Yes</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.accessToFamilyPlanningCounselling === false}
                    onChange={() => setFemaleHistory({...femaleHistory, accessToFamilyPlanningCounselling: false})}
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
                    onChange={() => setFemaleHistory({...femaleHistory, isPregnancyHistoryApplicable: true})}
                  />
                  <span>Yes</span>
                </label>
                <label>
                  <input 
                    type="radio"
                    checked={femaleHistory.isPregnancyHistoryApplicable === false}
                    onChange={() => setFemaleHistory({...femaleHistory, isPregnancyHistoryApplicable: false})}
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