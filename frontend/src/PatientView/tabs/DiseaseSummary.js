import React, { useEffect, useState } from 'react';
import './DiseaseSummary.css';

const CODE_MAP = {
  allergy: { code: '001', desc: 'Allergy' },
  asthma: { code: '002', desc: 'Asthma' },
  cancer: { code: '003', desc: 'Cancer' },
  cerebrovascularDisease: { code: '004', desc: 'Cerebrovascular Disease' },
  coronaryArteryDisease: { code: '005', desc: 'Coronary Artery Disease' },
  diabetesMellitus: { code: '006', desc: 'Diabetes Mellitus' },
  emphysema: { code: '007', desc: 'Emphysema' },
  epilepsySeizureDisorder: { code: '008', desc: 'Epilepsy/Seizure Disorder' },
  hepatitis: { code: '009', desc: 'Hepatitis' },
  hyperlipidemia: { code: '010', desc: 'Hyperlipidemia' },
  hypertension: { code: '011', desc: 'Hypertension' },
  pepticUlcer: { code: '012', desc: 'Peptic Ulcer' },
  pneumonia: { code: '013', desc: 'Pneumonia' },
  thyroidDisease: { code: '014', desc: 'Thyroid Disease' },
  pulmonaryTuberculosis: { code: '015', desc: 'Pulmonary Tuberculosis' },
  extrapulmonaryTuberculosis: { code: '016', desc: 'Extrapulmonary Tuberculosis' },
  urinaryTractInfection: { code: '017', desc: 'Urinary Tract Infection' },
  mentalIllness: { code: '018', desc: 'Mental Illness' },
  others: { code: '019', desc: 'Others' }
};

function prettyKey(key) {
  return key
    .replace(/([A-Z])/g, ' $1')
    .replace(/^./, (s) => s.toUpperCase());
}

const DiseaseSummary = () => {
  const [selected, setSelected] = useState({});

  useEffect(() => {
    try {
      const raw = localStorage.getItem('medicalHistorySelections');
      if (raw) setSelected(JSON.parse(raw));
    } catch (e) {
      setSelected({});
    }
  }, []);

  const entries = Object.keys(selected || {}).filter(k => selected[k]);

  return (
    <div className="disease-summary">
      <h3>Selected Diseases Summary</h3>
      {entries.length === 0 ? (
        <p>No diseases selected.</p>
      ) : (
        <div className="summary-list">
          <table>
            <thead>
              <tr>
                <th>Code</th>
                <th>Disease</th>
                <th>Description</th>
              </tr>
            </thead>
            <tbody>
              {entries.map(key => {
                const map = CODE_MAP[key] || { code: '---', desc: prettyKey(key) };
                return (
                  <tr key={key}>
                    <td>{map.code}</td>
                    <td>{prettyKey(key)}</td>
                    <td>{map.desc}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default DiseaseSummary;
