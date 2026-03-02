import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import './SocialHistory.css';

// SocialHistory tab — 6 questions about smoking, alcohol, drug use, and sexual activity.
// Purpose:
//   - Provide a compact UI for collecting a patient's social history (smoking, alcohol, drugs, sexual activity).
//   - Persist every change immediately to the backend so the UI and DB stay in sync.
// Uses:
//   - GET /api/patients/{patientId}/social-history to load saved data (reads from `patient_socialhistory`).
//   - POST /api/patients/{patientId}/social-history to save changes (handler upserts into `patient_socialhistory`).
// Design notes:
//   - The component is data-driven: `QUESTIONS` defines fields visible in the UI. This avoids duplicated JSX
//     and makes adding/removing fields a single change in the config.
//   - Immediate-save behavior matches `FemaleHistory` and `MedicalHistory` components (useCallback + setState then POST).

// QUESTIONS — defines all social history fields in a single place.
// Purpose and usage:
//   - Purpose: single source of truth for field metadata (key, label, type, options). Avoids hardcoding inputs.
//   - Usage: the render loop maps over this array to produce radio groups or selects. Keys must match backend columns.
//   - To add a field: add a new object with `key`, `label`, `type` and `options` (ensure `key` exists in DB/handler).
const QUESTIONS = [
  {
    key: 'is_patient_smoker',
    label: '1. Is Patient a Smoker',
    type: 'radio',
    options: [
      { value: 'Yes', label: 'Yes' },
      { value: 'No',  label: 'No'  },
      { value: 'Quit', label: 'Quit' },
    ],
  },
  {
    key: 'cigarette_packs_per_year',
    label: '2. Number of Cigarette Pack Consumed per Year',
    type: 'select',
    options: [0, 1, 2, 3, 4, 5, 10, 20, 30].map(n => ({ value: n, label: String(n) }))
      .concat([{ value: 50, label: '50+' }]),
  },
  {
    key: 'is_alcohol_drinker',
    label: '3. Is Patient an Alcohol Drinker',
    type: 'radio',
    options: [
      { value: 'Yes',  label: 'Yes'  },
      { value: 'No',   label: 'No'   },
      { value: 'Quit', label: 'Quit' },
    ],
  },
  {
    key: 'bottles_per_day',
    label: '4. Number of bottles consumed per day',
    type: 'select',
    options: [0, 1, 2, 3, 4, 5].map(n => ({ value: n, label: String(n) }))
      .concat([{ value: 10, label: '10+' }]),
  },
  {
    key: 'is_illicit_drug_user',
    label: '5. Is patient an illicit Drug User',
    type: 'radio',
    options: [
      { value: 'Yes', label: 'Yes' },
      { value: 'No',  label: 'No'  },
    ],
  },
  {
    key: 'is_sexually_active',
    label: '6. Is patient Sexually Active',
    type: 'radio',
    options: [
      { value: 'Yes', label: 'Yes' },
      { value: 'No',  label: 'No'  },
    ],
  },
];

// DEFAULT_STATE — derived from `QUESTIONS` so initial values remain consistent when fields change.
// Purpose:
//   - Ensures component state always contains the keys expected by the backend when loading/saving.
//   - For `select` fields we set a numeric default (0). For radio fields we default to 'No'.
// Note: if backend returns different default semantics, adjust this generation accordingly.
const DEFAULT_STATE = Object.fromEntries(
  QUESTIONS.map((q) => [q.key, q.type === 'select' ? 0 : 'No'])
);

function SocialHistory({ patientId }) {
  // socialHistory: mirrors patient_socialhistory columns via backend JSON tags
  const [socialHistory, setSocialHistory] = useState(DEFAULT_STATE);

  // fetchSocialHistory — loads saved values from backend on mount / patientId change.
  // Purpose:
  //   - Query server for previously saved `patient_socialhistory` row for the patient and merge into local state.
  //   - Keep UI authoritative values coming from the database (server is source of truth).
  // Endpoint: GET /api/patients/{patientId}/social-history
  const fetchSocialHistory = useCallback(async () => {
    if (!patientId) return;
    try {
      const res = await axios.get(`/api/patients/${patientId}/social-history`);
      if (res.data && Object.keys(res.data).length > 0) {
        setSocialHistory({ ...DEFAULT_STATE, ...res.data });
      }
    } catch (err) {
      // 404 means no record yet — keep defaults
      if (err.response?.status !== 404) {
        console.error('Failed to load social history', err);
      }
    }
  }, [patientId]);

  useEffect(() => {
    fetchSocialHistory();
  }, [fetchSocialHistory]);

  // saveSocialHistory — POSTs current state to backend immediately after every change.
  // Purpose:
  //   - Persist the full socialHistory object so the backend can upsert a single row for the patient.
  //   - Backend uses `patno` (case_no) as primary key and `ON DUPLICATE KEY UPDATE` to replace values.
  // Endpoint: POST /api/patients/{patientId}/social-history
  const saveSocialHistory = useCallback(async (next) => {
    if (!patientId) return;
    try {
      const res = await axios.post(`/api/patients/${patientId}/social-history`, {
        patient_id: parseInt(patientId),
        ...next,
      });
      if (res.data && Object.keys(res.data).length > 0) {
        setSocialHistory({ ...DEFAULT_STATE, ...res.data });
      }
    } catch (err) {
      console.error('Failed to save social history', err);
    }
  }, [patientId]);

  // handleChange — updates local state and immediately persists to DB.
  // Purpose:
  //   - Provide a single handler used by all inputs; it normalizes values and triggers `saveSocialHistory`.
  // Usage notes:
  //   - `type === 'select'` values are converted to integers before saving because DB columns are numeric.
  //   - The `saveSocialHistory` call sends the complete object so the backend has full context.
  const handleChange = useCallback((key, value, type) => {
    const parsed = type === 'select' ? parseInt(value) : value;
    setSocialHistory((prev) => {
      const next = { ...prev, [key]: parsed };
      saveSocialHistory(next);
      return next;
    });
  }, [saveSocialHistory]);

  return (
    <div className="medical-content">
      <div className="social-history-section">
        <h3>Social History</h3>

        <div className="social-history-questions">
          {/* Rendered from QUESTIONS config — no hardcoded JSX per question */}
          {/* Render notes:
              - Each question maps to a backend column named by `key`.
              - Inputs bind to `socialHistory[key]` so loaded values persist through re-renders.
              - This section intentionally contains no business logic; all logic lives in handlers above.
          */}
          {QUESTIONS.map((q) => (
            <div key={q.key} className="question-group">
              <label className="question-label">{q.label}</label>

              {q.type === 'radio' && (
                <div className="radio-group">
                  {q.options.map((opt) => (
                    <label key={opt.value} className="radio-option">
                      <input
                        type="radio"
                        name={q.key}
                        value={opt.value}
                        checked={socialHistory[q.key] === opt.value}
                        onChange={(e) => handleChange(q.key, e.target.value, 'radio')}
                      />
                      {opt.label}
                    </label>
                  ))}
                </div>
              )}

              {q.type === 'select' && (
                <select
                  name={q.key}
                  value={socialHistory[q.key]}
                  onChange={(e) => handleChange(q.key, e.target.value, 'select')}
                  className="select-input"
                >
                  {q.options.map((opt) => (
                    <option key={opt.value} value={opt.value}>{opt.label}</option>
                  ))}
                </select>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

export default SocialHistory;
