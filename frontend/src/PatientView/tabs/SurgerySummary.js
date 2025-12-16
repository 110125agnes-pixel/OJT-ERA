import React, { useEffect, useState } from 'react';
import './DiseaseSummary.css';

const SurgerySummary = () => {
  const [selected, setSelected] = useState([]);
  const [notes, setNotes] = useState('');

  useEffect(() => {
    try {
      const raw = localStorage.getItem('surgicalHistorySelections');
      if (raw) {
        const parsed = JSON.parse(raw);
        setSelected(parsed.selectedSurgeries || []);
        setNotes(parsed.notes || '');
      }
    } catch (e) {
      setSelected([]);
      setNotes('');
    }
  }, []);

  // Listen for updates dispatched from SurgicalHistory in the same window
  useEffect(() => {
    const handler = (ev) => {
      try {
        const payload = ev?.detail;
        if (payload) {
          setSelected(payload.selectedSurgeries || []);
          setNotes(payload.notes || '');
        } else {
          // fallback to localStorage
          const raw = localStorage.getItem('surgicalHistorySelections');
          if (raw) {
            const parsed = JSON.parse(raw);
            setSelected(parsed.selectedSurgeries || []);
            setNotes(parsed.notes || '');
          }
        }
      } catch (e) {
        setSelected([]);
        setNotes('');
      }
    };

    window.addEventListener('surgicalHistoryUpdated', handler);
    return () => window.removeEventListener('surgicalHistoryUpdated', handler);
  }, []);

  return (
    <div className="disease-summary">
      <h3>Selected Surgeries Summary</h3>
      {selected.length === 0 ? (
        <p>No surgeries selected.</p>
      ) : (
        <div className="summary-list">
          <table>
            <thead>
              <tr>
                <th>Code</th>
                <th>Procedure</th>
                <th>Description</th>
              </tr>
            </thead>
            <tbody>
              {selected.map((s, i) => (
                <tr key={s.id + '-' + i}>
                  <td>{s.code}</td>
                  <td>{s.name}</td>
                  <td>{s.name}</td>
                </tr>
              ))}
            </tbody>
          </table>
          {notes && (
            <div style={{ marginTop: 12 }}>
              <strong>Notes:</strong>
              <div>{notes}</div>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default SurgerySummary;
