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

  // Check if item is selected
  const isChecked = (item) => savedItems.some(s => mapCode(s) === mapCode(item));

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
      const exists = prev.some(s => mapCode(s) === mapCode(item));
      if (!exists) {
        return [...prev, item];
      } else {
        return prev.filter(s => mapCode(s) !== mapCode(item));
      }
    });
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
        </div>
      </div>
    </div>
  );
}

export default Immunization;