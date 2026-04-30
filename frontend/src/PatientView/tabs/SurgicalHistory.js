import React, { useState, useEffect, useCallback } from 'react';
import { surgicalService } from '../../services/api';
import './SurgicalHistory.css';

// SurgicalHistory tab — displays all surgical procedures from tsekap_lib_surgical as checkboxes.
// Patient selections are stored in patient_surgery (one row per patient, pipe-separated bit string).
//
// Key changes from original:
//   - Removed hardcoded surgery list; data now fetched from tsekap_lib_surgical via GET /api/lib/surgery
//   - Removed Save/Clear/Delete buttons; selections auto-save on every checkbox toggle
//   - Fixed JSON key mismatch: backend decodes SurgeryCode/IsChecked (PascalCase) matching frontend payload
//   - "None" mutual exclusivity: checking None unchecks all others; checking any item unchecks None
function SurgicalHistory({ patientId }) {
  // lib: ordered list of { code, desc } from tsekap_lib_surgical (LIB_STAT=1 only)
  const [lib, setLib] = useState([]);
  // selected: map of { [SURG_CODE]: true/false } for the current patient
  const [selected, setSelected] = useState({});
  const [loadingLib, setLoadingLib] = useState(false);
  const [errorLib, setErrorLib] = useState('');
  const [saving, setSaving] = useState(false);

  // fetchLibrary — loads all active surgery types from tsekap_lib_surgical (LIB_STAT=1)
  // Endpoint: GET /api/lib/surgery → [{ code, desc }] ordered by SORT_NO, SURG_CODE
  const fetchLibrary = useCallback(async () => {
    setLoadingLib(true);
    setErrorLib('');
    try {
      const data = await surgicalService.getSurgicalLibrary();
      setLib(data || []);
    } catch (err) {
      console.error('Failed to fetch surgical library', err);
      setErrorLib('Failed to load surgical library');
    } finally {
      setLoadingLib(false);
    }
  }, []);

  // fetchPatientSelections — loads this patient's saved choices from patient_surgery
  // Backend decodes the pipe-separated bit string (e.g. "1|0|0|1...") aligned to lib order
  // Endpoint: GET /api/patients/{patientId}/surgical-history → [{ SurgeryCode, SurgeryName, IsChecked }]
  const fetchPatientSelections = useCallback(async () => {
    if (!patientId) return;
    try {
      const res = await surgicalService.getPatientSurgicalHistory(patientId);
      const map = {};
      (res || []).forEach((it) => {
        map[it.SurgeryCode || it.surgery_code] = !!(it.IsChecked || it.is_checked);
      });
      setSelected(map);
    } catch (err) {
      console.error('Failed to fetch patient surgical selections', err);
    }
  }, [patientId]);

  useEffect(() => {
    fetchLibrary();
  }, [fetchLibrary]);

  useEffect(() => {
    fetchPatientSelections();
  }, [fetchPatientSelections]);

  // saveSelections — posts the full selection state to the backend after every toggle
  // Sends all lib items with IsChecked true/false so the backend can rebuild the bit string
  // Endpoint: POST /api/patients/{patientId}/surgical-history
  // Payload: [{ SurgeryCode, SurgeryName, IsChecked }] — PascalCase keys match Go struct tags
  // Backend upserts patient_surgery using ON DUPLICATE KEY UPDATE
  const saveSelections = useCallback(async (nextSelected, currentLib) => {
    if (!patientId) return;
    const libToUse = currentLib || lib;
    if (!libToUse || libToUse.length === 0) return;
    try {
      setSaving(true);
      const items = libToUse.map((s) => ({
        SurgeryCode: s.code,
        SurgeryName: s.desc,
        IsChecked: !!nextSelected[s.code],
      }));
      await surgicalService.savePatientSurgicalHistory(patientId, items);
      // update localStorage and notify other components (SurgerySummary)
      try {
        const selectedSurgeries = libToUse
          .filter((s) => !!nextSelected[s.code])
          .map((s) => ({ code: s.code, name: s.desc }));
        const payload = { selectedSurgeries, notes: '' };
        localStorage.setItem('surgicalHistorySelections', JSON.stringify(payload));
        window.dispatchEvent(new CustomEvent('surgicalHistoryUpdated', { detail: payload }));
      } catch (e) {
        // ignore storage errors
      }
    } catch (err) {
      console.error('Failed to save surgical selections', err);
    } finally {
      setSaving(false);
    }
  }, [patientId, lib]);

  // handleToggle — called when any checkbox changes
  // Mutual exclusivity rules for "None":
  //   - If "None" is checked: clear all other selections (only None stays checked)
  //   - If any other item is checked: automatically uncheck "None"
  // After applying the rules, auto-saves the updated selection to the backend
  const handleToggle = (code, checked) => {
    const snapshot = lib; // capture current lib order for the save call
    const noneCode = lib.find((s) => s.desc && s.desc.trim().toLowerCase() === 'none')?.code;
    setSelected((prev) => {
      let next = { ...prev, [code]: checked };
      if (checked) {
        if (noneCode && code === noneCode) {
          // "None" checked — uncheck everything else, keep only None
          next = {};
          next[noneCode] = true;
        } else if (noneCode) {
          // Any other item checked — uncheck "None"
          next[noneCode] = false;
        }
      }
      saveSelections(next, snapshot);
      return next;
    });
  };

  // Keep shared storage in sync so summary components in the same window can read it
  useEffect(() => {
    if (!lib || lib.length === 0) return;
    try {
      const selectedSurgeries = lib.filter((s) => !!selected[s.code]).map((s) => ({ code: s.code, name: s.desc }));
      const payload = { selectedSurgeries, notes: '' };
      localStorage.setItem('surgicalHistorySelections', JSON.stringify(payload));
      window.dispatchEvent(new CustomEvent('surgicalHistoryUpdated', { detail: payload }));
    } catch (e) {
      // ignore storage errors
    }
  }, [lib, selected]);

  return (
    <div className="surgical-history-container">
      <div className="surgical-header">
        <h3>Surgical History Specifics</h3>
        {saving && <div className="saving-indicator">Saving...</div>}
      </div>

      <div className="surgical-history-section">
        {loadingLib ? (
          <div className="surgical-loading">Loading surgical procedures...</div>
        ) : errorLib ? (
          <div className="surgical-error">{errorLib}</div>
        ) : lib.length === 0 ? (
          <div className="surgical-no-data">No surgical procedures available.</div>
        ) : (
          <>
            <div className="surgical-checkboxes">
              <div className="checkbox-group">
                {lib.map((s) => (
                  <label key={s.code}>
                    <input
                      type="checkbox"
                      checked={!!selected[s.code]}
                      onChange={(e) => handleToggle(s.code, e.target.checked)}
                    />
                    <span>{s.desc}</span>
                  </label>
                ))}
              </div>
            </div>

            <div className="surgical-codes-table">
              <table>
                <thead>
                  <tr>
                    <th>Code</th>
                    <th>Description</th>
                  </tr>
                </thead>
                <tbody>
                  {lib.map((s) => {
                    if (!selected[s.code]) return null;
                    return (
                      <tr key={s.code}>
                        <td>{s.code}</td>
                        <td>{s.desc}</td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default SurgicalHistory;