import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import './Immunization.css';

function Immunization({ patientId }) {
  const [childLib, setChildLib] = useState([]);
  const [youngLib, setYoungLib] = useState([]);
  const [pregLib, setPregLib] = useState([]);
  const [elderlyLib, setElderlyLib] = useState([]);

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

  useEffect(() => {
    fetchLibs();
  }, [fetchLibs]);

  const mapCode = (item) => item.code || item.vaccine_code || item.IMM_CODE || item.imm_code || item.Code || item.code;
  const mapName = (item) => item.name || item.vaccine_name || item.IMM_DESC || item.imm_desc || item.Desc || item.desc || '';

  return (
    <div className="immunization-content">
      <div className="immunization-sections">
        <div className="immunization-left">
          {/* Children Section */}
          <div className="immunization-section">
            <h4>1. Children</h4>
            <div className="immunization-grid">
              {childLib.map((item) => {
                const code = mapCode(item);
                const name = mapName(item);
                return (
                  <label key={code}>
                    <input type="checkbox" />
                    <span>{name}</span>
                  </label>
                );
              })}
            </div>
          </div>

          {/* Young Section */}
          <div className="immunization-section">
            <h4>2. Young</h4>
            <div className="immunization-grid">
              {youngLib.map((item) => {
                const code = mapCode(item);
                const name = mapName(item);
                return (
                  <label key={code}>
                    <input type="checkbox" />
                    <span>{name}</span>
                  </label>
                );
              })}
            </div>
          </div>

          {/* Pregnant Section */}
          <div className="immunization-section">
            <h4>3. Pregnant</h4>
            <div className="immunization-grid">
              {pregLib.map((item) => {
                const code = mapCode(item);
                const name = mapName(item);
                return (
                  <label key={code}>
                    <input type="checkbox" />
                    <span>{name}</span>
                  </label>
                );
              })}
            </div>
          </div>

          {/* Elderly Section */}
          <div className="immunization-section">
            <h4>4. Elderly</h4>
            <div className="immunization-grid">
              {elderlyLib.map((item) => {
                const code = mapCode(item);
                const name = mapName(item);
                return (
                  <label key={code}>
                    <input type="checkbox" />
                    <span>{name}</span>
                  </label>
                );
              })}
            </div>
          </div>
        </div>

        <div className="immunization-right">
          {/* Other Immunization Section */}
          <div className="immunization-section other-section">
            <h4>5. Other Immunization</h4>
            <textarea placeholder="Other immunization notes" rows="10"></textarea>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Immunization;