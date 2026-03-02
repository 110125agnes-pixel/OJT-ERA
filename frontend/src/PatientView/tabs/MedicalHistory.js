import React, { useState, useEffect, useCallback } from "react";
import axios from "axios";
import "./MedicalHistory.css";
import Button from '@mui/material/Button';
import ButtonGroup from '@mui/material/ButtonGroup';

// MedicalHistory tab — shows checkboxes for all diseases in tsekap_lib_mdiseases.
// Selections auto-save to the backend on every toggle (no Save button needed).
// Data persists in the patient_medhist DB table as a pipe-separated 1|0 string.
function MedicalHistory({ patientId }) {
  // lib: full disease list loaded from tsekap_lib_mdiseases via /api/lib/mdiseases
  const [lib, setLib] = useState([]);
  // selected: map of { disease_code: true/false } representing current checkbox states
  const [selected, setSelected] = useState({});

  // fetchLibrary: loads all available diseases from the backend library table.
  // Runs once on mount and refreshes every 5 seconds so new/deleted diseases reflect live.
  const fetchLibrary = useCallback(async () => {
    try {
      const res = await axios.get('/api/lib/mdiseases');
      setLib(res.data || []);
    } catch (err) {
      console.error('Failed to fetch library', err);
    }
  }, []);

  // fetchPatientSelections: loads this patient's saved checkbox state from the backend.
  // Called on mount so previously saved selections are restored after refresh or re-open.
  // Backend is the single source of truth — no localStorage fallback.
  const fetchPatientSelections = useCallback(async () => {
    if (!patientId) return;
    try {
      const res = await axios.get(`/api/patients/${patientId}/medical-history`);
      const map = {};
      // Convert the array response [{disease_code, is_checked}] into a lookup map
      (res.data || []).forEach((it) => {
        map[it.disease_code] = !!it.is_checked;
      });
      setSelected(map);
    } catch (err) {
      console.error('Failed to fetch patient selections', err);
    }
  }, [patientId]);

  // On mount: load the disease library and start the 5-second refresh interval
  useEffect(() => {
    fetchLibrary();
    const t = setInterval(fetchLibrary, 5000);
    return () => clearInterval(t); // cleanup on unmount
  }, [fetchLibrary]);

  // On mount (and whenever patientId changes): load saved selections from backend
  useEffect(() => {
    fetchPatientSelections();
  }, [fetchPatientSelections]);

  // saveSelections: POSTs the full current checkbox state to the backend immediately.
  // Takes a snapshot of lib at click time (currentLib) to avoid stale-closure issues.
  // Guard: if library is empty, skip save to prevent wiping the patient's data.
  const saveSelections = useCallback(async (nextSelected, currentLib) => {
    if (!patientId) return;
    const libToUse = currentLib || lib;
    if (!libToUse || libToUse.length === 0) return; // don't save if library not loaded yet
    try {
      // Build array of all diseases with their current checked state
      const items = libToUse.map((d) => ({
        disease_code: d.mdisease_code || d.Code || d.code,
        disease_name: d.mdisease_desc || d.Desc || d.desc,
        is_checked: !!nextSelected[(d.mdisease_code || d.Code || d.code)]
      }));
      await axios.post(`/api/patients/${patientId}/medical-history`, items);
    } catch (err) {
      console.error('Failed to save selections', err);
    }
  }, [patientId, lib]);

  // handleToggle: called when a checkbox is clicked.
  // Special logic: checking "None" (code 999) unchecks everything else.
  // Checking anything else unchecks "None" if it was active.
  // After updating state, immediately calls saveSelections to persist to DB.
  const handleToggle = (code, checked) => {
    const snapshot = lib; // capture current lib before any async state updates
    setSelected((prev) => {
      const next = { ...prev };
      const isNone = code === '999' || code === 'None';

      if (isNone && checked) {
        // "None" checked: uncheck all other diseases
        for (const k of Object.keys(next)) next[k] = false;
        next[code] = true;
      } else if (isNone && !checked) {
        next[code] = false;
      } else {
        // Any other disease checked: auto-uncheck "None" if it was active
        if (Object.keys(next).some(k => (k === '999' || k === 'None') && next[k])) {
          next['999'] = false;
          next['None'] = false;
        }
        next[code] = checked;
      }

      // Auto-save immediately on every toggle — no manual Save button needed
      saveSelections(next, snapshot);
      return next;
    });
  };

  return (
    <div className="medical-content">
      <div className="medical-header">
        <h3>Medical History Specifics</h3>
        
      </div>

      <div className="medical-history-section">
        <div className="medical-checkboxes">
          {/* Render one checkbox per disease in lib (order from tsekap_lib_mdiseases ORDER BY mdisease_code).
              All checkboxes are disabled while "None" (999) is selected — only "None" stays clickable. */}
          <div className="checkbox-group">
            {lib.map((d) => {
              const code = d.mdisease_code || d.Code || d.code;
              const desc = d.mdisease_desc || d.Desc || d.desc || '';
              const isNone = code === '999' || (desc && desc.toLowerCase() === 'none');
              const isOthers = code === '998';
              return (
                <label key={code}>
                  <input
                    type="checkbox"
                    checked={!!selected[code]}
                    // Disable all disease checkboxes while "None" is active (except "None" itself)
                    disabled={!!selected['999'] && !isNone}
                    onChange={(e) => handleToggle(code, e.target.checked)}
                  />
                  <span>{desc}</span>
                </label>
              );
            })}
          </div>

          <div className="medical-details">
            <textarea placeholder="Enter medical history details here..." rows="10"></textarea>
          </div>
        </div>

        {/* Summary table: shows only the diseases currently checked by this patient.
            Filters lib by selected[code] — source of truth is still the selected state map. */}
        <div className="medical-codes-table">
          <table>
            <thead>
              <tr>
                <th>Code</th>
                <th>Description</th>
              </tr>
            </thead>
            <tbody>
              {lib.map((d) => {
                const code = d.mdisease_code || d.Code || d.code;
                const desc = d.mdisease_desc || d.Desc || d.desc || '';
                if (!selected[code]) return null; // skip unchecked diseases
                return (
                  <tr key={code}>
                    <td>{code}</td>
                    <td>{desc}</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

export default MedicalHistory;
