import React, { useState } from 'react';
import './FamilyHistory.css';

const FamilyHistory = () => {
  const [selectedDiseases, setSelectedDiseases] = useState({
    none: false,
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

  const [notes, setNotes] = useState('');
  const [tableData, setTableData] = useState([]);

  const diseaseMapping = {
    none: { code: '000', name: 'None' },
    allergy: { code: '001', name: 'Allergy' },
    asthma: { code: '002', name: 'Asthma' },
    cancer: { code: '003', name: 'Cancer' },
    cerebrovascularDisease: { code: '004', name: 'Cerebrovascular Disease' },
    coronaryArteryDisease: { code: '005', name: 'Coronary Artery Disease' },
    diabetesMellitus: { code: '006', name: 'Diabetes Mellitus' },
    emphysema: { code: '007', name: 'Emphysema' },
    epilepsySeizureDisorder: { code: '008', name: 'Epilepsy/Seizure Disorder' },
    hepatitis: { code: '009', name: 'Hepatitis' },
    hyperlipidemia: { code: '010', name: 'Hyperlipidemia' },
    hypertension: { code: '011', name: 'Hypertension' },
    pepticUlcer: { code: '012', name: 'Peptic Ulcer' },
    pneumonia: { code: '013', name: 'Pneumonia' },
    thyroidDisease: { code: '014', name: 'Thyroid Disease' },
    pulmonaryTuberculosis: { code: '015', name: 'Pulmonary Tuberculosis' },
    extrapulmonaryTuberculosis: { code: '016', name: 'Extrapulmonary Tuberculosis' },
    urinaryTractInfection: { code: '017', name: 'Urinary Tract Infection' },
    mentalIllness: { code: '018', name: 'Mental Illness' },
    others: { code: '998', name: 'Others' }
  };

  const handleCheckboxChange = (disease) => {
    setSelectedDiseases(prev => ({
      ...prev,
      [disease]: !prev[disease]
    }));
  };

  const handleAdd = () => {
    const selected = Object.entries(selectedDiseases)
      .filter(([_, isSelected]) => isSelected)
      .map(([disease]) => ({
        code: diseaseMapping[disease].code,
        description: diseaseMapping[disease].name
      }));

    if (selected.length > 0) {
      setTableData(selected);
    }
  };

  const handleSave = () => {
    console.log('Saving data:', { selectedDiseases, notes, tableData });
    alert('Data saved successfully!');
  };

  const handleClear = () => {
    setSelectedDiseases({
      none: false,
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
    setNotes('');
    setTableData([]);
  };

  return (
    <div className="family-history-container">
      <div className="family-history-content">
        {/* Left Panel - Checkboxes */}
        <div className="family-history-left">
          <div className="family-header">
            <h3>Family History Specifics</h3>
          </div>
          <div className="checkbox-list-family">
            <label><input type="checkbox" checked={selectedDiseases.none} onChange={() => handleCheckboxChange('none')} /> None</label>
            <label><input type="checkbox" checked={selectedDiseases.allergy} onChange={() => handleCheckboxChange('allergy')} /> Allergy</label>
            <label><input type="checkbox" checked={selectedDiseases.asthma} onChange={() => handleCheckboxChange('asthma')} /> Asthma</label>
            <label><input type="checkbox" checked={selectedDiseases.cancer} onChange={() => handleCheckboxChange('cancer')} /> Cancer</label>
            <label><input type="checkbox" checked={selectedDiseases.cerebrovascularDisease} onChange={() => handleCheckboxChange('cerebrovascularDisease')} /> Cerebrovascular Disease</label>
            <label><input type="checkbox" checked={selectedDiseases.coronaryArteryDisease} onChange={() => handleCheckboxChange('coronaryArteryDisease')} /> Coronary Artery Disease</label>
            <label><input type="checkbox" checked={selectedDiseases.diabetesMellitus} onChange={() => handleCheckboxChange('diabetesMellitus')} /> Diabetes Mellitus</label>
            <label><input type="checkbox" checked={selectedDiseases.emphysema} onChange={() => handleCheckboxChange('emphysema')} /> Emphysema</label>
            <label><input type="checkbox" checked={selectedDiseases.epilepsySeizureDisorder} onChange={() => handleCheckboxChange('epilepsySeizureDisorder')} /> Epilepsy/Seizure Disorder</label>
            <label><input type="checkbox" checked={selectedDiseases.hepatitis} onChange={() => handleCheckboxChange('hepatitis')} /> Hepatitis</label>
            <label><input type="checkbox" checked={selectedDiseases.hyperlipidemia} onChange={() => handleCheckboxChange('hyperlipidemia')} /> Hyperlipidemia</label>
            <label><input type="checkbox" checked={selectedDiseases.hypertension} onChange={() => handleCheckboxChange('hypertension')} /> Hypertension</label>
            <label><input type="checkbox" checked={selectedDiseases.pepticUlcer} onChange={() => handleCheckboxChange('pepticUlcer')} /> Peptic Ulcer</label>
            <label><input type="checkbox" checked={selectedDiseases.pneumonia} onChange={() => handleCheckboxChange('pneumonia')} /> Pneumonia</label>
            <label><input type="checkbox" checked={selectedDiseases.thyroidDisease} onChange={() => handleCheckboxChange('thyroidDisease')} /> Thyroid Disease</label>
            <label><input type="checkbox" checked={selectedDiseases.pulmonaryTuberculosis} onChange={() => handleCheckboxChange('pulmonaryTuberculosis')} /> Pulmonary Tuberculosis</label>
            <label><input type="checkbox" checked={selectedDiseases.extrapulmonaryTuberculosis} onChange={() => handleCheckboxChange('extrapulmonaryTuberculosis')} /> Extrapulmonary Tuberculosis</label>
            <label><input type="checkbox" checked={selectedDiseases.urinaryTractInfection} onChange={() => handleCheckboxChange('urinaryTractInfection')} /> Urinary Tract Infection</label>
            <label><input type="checkbox" checked={selectedDiseases.mentalIllness} onChange={() => handleCheckboxChange('mentalIllness')} /> Mental Illness</label>
            <label><input type="checkbox" checked={selectedDiseases.others} onChange={() => handleCheckboxChange('others')} /> Others</label>
          </div>
        </div>

        {/* Right Panel - Notes and Table */}
        <div className="family-history-right">
          <div className="action-buttons-top">
            <button onClick={handleAdd}>Add</button>
            <button onClick={handleSave}>Save</button>
            <button onClick={handleClear}>Clear</button>
          </div>
          
          <textarea
            className="notes-textarea"
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
            placeholder="Limit your characters to 2000"
          />

          <div className="table-container-family">
            <table className="family-table">
              <thead>
                <tr>
                  <th>Code</th>
                  <th>Description</th>
                </tr>
              </thead>
              <tbody>
                {tableData.map((item, index) => (
                  <tr key={index}>
                    <td>{item.code}</td>
                    <td>{item.description}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  );
};

export default FamilyHistory;