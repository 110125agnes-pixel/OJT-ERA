import React, { useState, useEffect, useCallback } from "react";
import axios from "axios";
import "./MedicalHistory.css";

function MedicalHistory({ patientId }) {
  const [lib, setLib] = useState([]); // {mdisease_code, mdisease_desc}
  const [selected, setSelected] = useState({}); // code => bool
  const [loading, setLoading] = useState(false);

  const fetchLibrary = useCallback(async () => {
    try {
      const res = await axios.get('/api/lib/mdiseases');
      setLib(res.data || []);
    } catch (err) {
      console.error('Failed to fetch library', err);
    }
  }, []);

  const fetchPatientSelections = useCallback(async () => {
    if (!patientId) return;
    try {
      const res = await axios.get(`/api/patients/${patientId}/medical-history`);
      const map = {};
      (res.data || []).forEach((it) => {
        map[it.disease_code] = !!it.is_checked;
      });
      setSelected(map);
    } catch (err) {
      console.error('Failed to fetch patient selections', err);
      // fallback: load from localStorage so selections persist even without backend
      try {
        const key = `medicalHistorySelections_${patientId}`;
        const saved = JSON.parse(localStorage.getItem(key) || '{}');
        setSelected(saved || {});
      } catch (e) {
        setSelected({});
      }
    }
  }, [patientId]);

  useEffect(() => {
    fetchLibrary();
    const t = setInterval(fetchLibrary, 5000);
    return () => clearInterval(t);
  }, [fetchLibrary]);

  useEffect(() => {
    fetchPatientSelections();
  }, [fetchPatientSelections]);

  const saveSelections = async (nextSelected) => {
    if (!patientId) return;
    // always save locally as a fallback for persistence
    try {
      const key = `medicalHistorySelections_${patientId}`;
      localStorage.setItem(key, JSON.stringify(nextSelected));
    } catch (e) {}
    setLoading(true);
    try {
      const items = lib.map((d) => ({
        disease_code: d.mdisease_code || d.Code || d.code,
        disease_name: d.mdisease_desc || d.Desc || d.desc,
        is_checked: !!nextSelected[(d.mdisease_code || d.Code || d.code)]
      }));
      await axios.post(`/api/patients/${patientId}/medical-history`, items);
    } catch (err) {
      console.error('Failed to save selections', err);
    } finally {
      setLoading(false);
    }
  };

  const handleToggle = (code, checked) => {
    setSelected((prev) => {
      const next = { ...prev };
      const isNone = code === '999' || code === 'None';
      if (isNone && checked) {
        for (const k of Object.keys(next)) next[k] = false;
        next[code] = true;
      } else if (isNone && !checked) {
        next[code] = false;
      } else {
        if (Object.keys(next).some(k => (k === '999' || k === 'None') && next[k])) {
          next['999'] = false;
          next['None'] = false;
        }
        next[code] = checked;
      }
      saveSelections(next);
      return next;
    });
  };

  return (
    <div className="medical-content">
      <div className="medical-history-section">
        <h3>Medical History Specifics</h3>
        <div className="medical-checkboxes">
          <div className="checkbox-group">
            {lib.map((d) => {
              const code = d.mdisease_code || d.Code || d.code;
              const desc = d.mdisease_desc || d.Desc || d.desc || '';
              const isNone = code === '999' || (desc && desc.toLowerCase() === 'none');
              return (
                <label key={code}>
                  <input
                    type="checkbox"
                    checked={!!selected[code]}
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
                if (!selected[code]) return null;
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
