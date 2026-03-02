import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import './Immunization.css';

function Immunization({ patientId }) {
  const [childLib, setChildLib] = useState([]);
  const [youngLib, setYoungLib] = useState([]);
  const [pregLib, setPregLib] = useState([]);
  const [elderlyLib, setElderlyLib] = useState([]);
  const [savedItems, setSavedItems] = useState([]);

  // Map code and name with fallbacks
  const mapCode = (item) => item.code || item.vaccine_code || item.IMM_CODE || item.imm_code || item.Code || item.code;
  const mapName = (item) => item.name || item.vaccine_name || item.IMM_DESC || item.imm_desc || item.Desc || item.desc || '';

  // Fetch all immunization libraries
  const fetchLibs = useCallback(async () => {
    try {
      const [cRes, yRes, pRes, eRes] = await Promise.all([
        axios.get('/api/lib/immchild'),
        axios.get('/api/lib/immyoungw'),
        axios.get('/api/lib/immpregw'),
        axios.get('/api/lib/immelderly'),
      ]);
      setChildLib(cRes.data || []);
      setYoungLib(yRes.data || []);
      setPregLib(pRes.data || []);
      setElderlyLib(eRes.data || []);
    } catch (err) {
      console.error('Failed to fetch immunization libraries', err);
    }
  }, []);

  // Fetch saved immunizations for the patient
  const fetchSavedImmunization = useCallback(async () => {
    try {
      const res = await axios.get(`/api/patients/${patientId}/immunization`);
      setSavedItems(res.data || []);
    } catch (err) {
      console.error('Failed to fetch saved immunizations', err);
    }
  }, [patientId]);

  useEffect(() => {
    fetchLibs();
    fetchSavedImmunization();
  }, [fetchLibs, fetchSavedImmunization]);

  // Check if item is selected (honor the saved is_checked flag)
  const isChecked = (item) => {
    const s = savedItems.find(si => mapCode(si) === mapCode(item));
    if (!s) return false;
    return !!(s.is_checked || s.IsChecked || s.checked || s.isChecked);
  };

  // Handle checkbox toggle
  const handleCheckboxChange = (lib, setLib, item) => {
    // Update library state
    const updatedLib = lib.map(v =>
      mapCode(v) === mapCode(item)
        ? { ...v, is_checked: !v.is_checked, category: item.category }
        : v
    );
    setLib(updatedLib);

    // Update savedItems state
    setSavedItems(prev => {
      const idx = prev.findIndex(s => mapCode(s) === mapCode(item));
      if (idx === -1) {
        // add with checked=true
        return [...prev, { ...item, is_checked: true }];
      }

      // toggle the is_checked flag on the existing item
      const updated = [...prev];
      const cur = updated[idx];
      const curVal = !!(cur.is_checked || cur.IsChecked || cur.checked || cur.isChecked);
      updated[idx] = { ...cur, is_checked: !curVal };
      return updated;
    });
  };

  // Save current selections to backend
  const handleSave = async () => {
    try {
      // Build payload using the library order so backend can map codes correctly
      const allLibs = [...childLib, ...youngLib, ...pregLib, ...elderlyLib];
      const payload = allLibs.map(libItem => {
        const code = mapCode(libItem);
        const name = mapName(libItem);
        const saved = savedItems.find(s => (s.vaccine_code || s.VaccineCode || s.IMM_CODE || s.code) === code);
        const is_checked = !!(saved && (saved.is_checked || saved.IsChecked || saved.checked || saved.isChecked));
        return {
          vaccine_code: code,
          vaccine_name: name,
          category: libItem.category || libItem.Category || libItem.grp || '',
          is_checked,
        };
      });

      await axios.post(`/api/patients/${patientId}/immunization`, payload);
      // Optionally refetch saved items to sync state
      fetchSavedImmunization();
      alert('Immunization saved');
    } catch (err) {
      console.error('Failed to save immunizations', err);
      alert('Save failed');
    }
  };

  // Render a library section
  const renderSection = (title, lib, setLib) => (
    <div className="immunization-section">
      <h4>{title}</h4>
      <div className="immunization-grid">
        {lib.map(item => {
          const code = mapCode(item);
          const name = mapName(item);
          return (
            <label key={code}>
              <input
                type="checkbox"
                checked={isChecked(item)}
                onChange={() => handleCheckboxChange(lib, setLib, item)}
              />
              <span>{name}</span>
            </label>
          );
        })}
      </div>
    </div>
  );

  return (
    <div className="immunization-content">
      <div className="immunization-sections">
        <div className="immunization-left">
          {renderSection('1. Children', childLib, setChildLib)}
          {renderSection('2. Young', youngLib, setYoungLib)}
          {renderSection('3. Pregnant', pregLib, setPregLib)}
          {renderSection('4. Elderly', elderlyLib, setElderlyLib)}
        </div>

        <div className="immunization-right">
          <div className="immunization-section other-section">
            <h4>5. Other Immunization</h4>
            <textarea
              placeholder="Other immunization notes"
              rows="10"
              value={savedItems.find(s => s.other_notes)?.other_notes || ''}
              onChange={(e) => {
                const notes = e.target.value;
                setSavedItems(prev => {
                  const otherIndex = prev.findIndex(s => s.other_notes !== undefined);
                  if (otherIndex > -1) {
                    const updated = [...prev];
                    updated[otherIndex].other_notes = notes;
                    return updated;
                  } else {
                    return [...prev, { other_notes: notes }];
                  }
                });
              }}
            ></textarea>
          </div>
          <div style={{marginTop:16}}>
            <button onClick={handleSave}>Save Immunization</button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Immunization;