import React, { useState, useEffect } from "react";
import "./MedicalHistory.css";

function MedicalHistory() {
  const saved = (() => {
    try {
      return JSON.parse(localStorage.getItem('medicalHistorySelections') || '{}');
    } catch (e) {
      return {};
    }
  })();

  const [history, setHistory] = useState({
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
    others: false,
  });

  // On mount, merge saved selections
  useEffect(() => {
    if (saved && Object.keys(saved).length) {
      setHistory(prev => ({ ...prev, ...saved }));
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleChange = (field, value) => {
    setHistory((prev) => {
      const next = { ...prev, [field]: value };
      try {
        localStorage.setItem('medicalHistorySelections', JSON.stringify(next));
      } catch (e) {}
      return next;
    });
  };

  return (
    <div className="medical-content">
      <div className="medical-history-section">
        <h3>Medical History Specifics</h3>
        <div className="medical-checkboxes">
          {/* Checkboxes Group */}
          <div className="checkbox-group">
            {Object.keys(history).map((key) => (
              <label key={key}>
                <input
                  type="checkbox"
                  checked={history[key]}
                  onChange={(e) => handleChange(key, e.target.checked)}
                />
                {/* Converts camelCase to Title Case (e.g. diabetesMellitus -> Diabetes Mellitus) */}
                <span>
                  {key
                    .replace(/([A-Z])/g, " $1")
                    .replace(/^./, (str) => str.toUpperCase())}
                </span>
              </label>
            ))}
          </div>

          <div className="medical-details">
            <textarea
              placeholder="Enter medical history details here..."
              rows="10"
            ></textarea>
          </div>
        </div>

        {/* Dynamic Code Table */}
        <div className="medical-codes-table">
          <table>
            <thead>
              <tr>
                <th>Code</th>
                <th>Description</th>
              </tr>
            </thead>
            <tbody>
              {history.allergy && (
                <tr>
                  <td>001</td>
                  <td>Allergy</td>
                </tr>
              )}
              {history.asthma && (
                <tr>
                  <td>002</td>
                  <td>Asthma</td>
                </tr>
              )}
              {history.cancer && (
                <tr>
                  <td>003</td>
                  <td>Cancer</td>
                </tr>
              )}
              {history.diabetesMellitus && (
                <tr>
                  <td>006</td>
                  <td>Diabetes Mellitus</td>
                </tr>
              )}
              {history.emphysema && (
                <tr>
                  <td>007</td>
                  <td>Emphysema</td>
                </tr>
              )}
              {history.pepticUlcer && (
                <tr>
                  <td>012</td>
                  <td>Peptic Ulcer</td>
                </tr>
              )}
              {history.thyroidDisease && (
                <tr>
                  <td>014</td>
                  <td>Thyroid Disease</td>
                </tr>
              )}
              {/* Add other mappings here as needed */}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

export default MedicalHistory;
