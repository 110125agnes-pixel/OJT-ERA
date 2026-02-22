import React, { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import './FamilyHistory.css';

const FamilyHistory = ({ patientId }) => {
  const [lib, setLib] = useState([]); // disease library from mdiseases
  const [selected, setSelected] = useState({}); // code => bool
  const [notes, setNotes] = useState('');
  const [tableData, setTableData] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetchLibrary = useCallback(async () => {
    try {
      const res = await axios.get('/api/lib/mdiseases');
      setLib(res.data || []);
    } catch (err) {
      console.error('Failed to fetch library', err);
    }
  }, []);

  const fetchFamily = useCallback(async () => {
    if (!patientId) return;
    try {
      const res = await axios.get(`/api/patients/${patientId}/family-history`);
      const items = res.data || [];
      const nextSelected = {};
      const table = [];
      items.forEach(it => {
        nextSelected[it.disease_code] = !!it.is_checked;
        if (it.is_checked) table.push({ code: it.disease_code, description: it.disease_name || '' });
      });
      setSelected(nextSelected);
      setTableData(table);
      if (items.length > 0 && items[0].notes) setNotes(items[0].notes);
    } catch (err) {
      console.error('Failed to fetch family history', err);
      try {
        const key = `familyHistory_${patientId}`;
        const saved = JSON.parse(localStorage.getItem(key) || '{}');
        if (saved.selected) setSelected(saved.selected);
        if (saved.tableData) setTableData(saved.tableData);
        if (saved.notes) setNotes(saved.notes);
      } catch (e) {}
    }
  }, [patientId]);

  useEffect(() => {
    fetchLibrary();
  }, [fetchLibrary]);

  useEffect(() => {
    fetchFamily();
  }, [fetchFamily]);

  const saveFamily = async (payload) => {
    if (!patientId) return;
    setLoading(true);
    try {
      await axios.post(`/api/patients/${patientId}/family-history`, payload);
    } catch (err) {
      console.error('Failed to save family history', err);
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const handleToggle = (code, checked) => {
    setSelected(prev => {
      const next = { ...prev };
      const desc = (lib.find(d => (d.mdisease_code || d.Code || d.code) === code) || {}).mdisease_desc || '';
      const isNone = code === '999' || (desc && desc.toLowerCase() === 'none');
      if (isNone && checked) {
        // clear others
        for (const k of Object.keys(next)) next[k] = false;
        next[code] = true;
      } else if (isNone && !checked) {
        next[code] = false;
      } else {
        // if None is set, clear it
        for (const k of Object.keys(next)) {
          const kd = (lib.find(d => (d.mdisease_code || d.Code || d.code) === k) || {}).mdisease_desc || '';
          if (kd && kd.toLowerCase() === 'none') next[k] = false;
        }
        next[code] = checked;
      }

      // update tableData from lib
      const table = lib
        .map(d => ({ code: d.mdisease_code || d.Code || d.code, description: d.mdisease_desc || d.Desc || d.desc || '' }))
        .filter(d => next[d.code]);

      setTableData(table);

      // local fallback
      try {
        const key = `familyHistory_${patientId || 'local'}`;
        localStorage.setItem(key, JSON.stringify({ selected: next, tableData: table, notes }));
      } catch (e) {}

      return next;
    });
  };

  const handleAdd = () => {
    const table = lib
      .map(d => ({ code: d.mdisease_code || d.Code || d.code, description: d.mdisease_desc || d.Desc || d.desc || '' }))
      .filter(d => selected[d.code]);
    if (table.length > 0) setTableData(table);
  };

  const handleSave = async () => {
    const payload = tableData.map(t => ({ disease_code: t.code, disease_name: t.description, notes: notes || '', is_checked: true }));
    try {
      // always save locally
      const key = `familyHistory_${patientId || 'local'}`;
      localStorage.setItem(key, JSON.stringify({ selected, tableData, notes }));
    } catch (e) {}

    if (!patientId) {
      alert('Saved locally (no patient selected)');
      return;
    }

    try {
      await saveFamily(payload);
      alert('Family history saved');
    } catch (err) {
      alert('Failed to save to server; changes saved locally');
    }
  };

  const handleClear = () => {
    setSelected({});
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
                  {desc}
                </label>
              );
            })}
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