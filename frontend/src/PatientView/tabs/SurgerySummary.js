import React, { useState, useEffect } from 'react';
import './SurgicalHistory.css';
import { surgicalService } from '../../services/api';

const defaultSurgeries = [
  { id: 1, code: 'S001', name: 'Appendectomy' },
  { id: 2, code: 'S002', name: 'Cholecystectomy' },
  { id: 3, code: 'S003', name: 'Hernia Repair' },
  { id: 4, code: 'S004', name: 'Cesarean Section' },
  { id: 5, code: 'S005', name: 'Hysterectomy' },
  { id: 6, code: 'S006', name: 'Tonsillectomy' },
  { id: 7, code: 'S007', name: 'Cardiac Surgery' },
  { id: 8, code: 'S008', name: 'Orthopedic Surgery' },
  { id: 9, code: 'S009', name: 'Cataract Surgery' },
  { id: 10, code: 'S010', name: 'Mastectomy' },
  { id: 11, code: 'S011', name: 'Prostatectomy' },
  { id: 12, code: 'S012', name: 'Thyroidectomy' },
  { id: 13, code: 'S013', name: 'Spinal Surgery' },
  { id: 14, code: 'S014', name: 'Gastric Bypass' },
  { id: 15, code: 'S015', name: 'Kidney Surgery' },
  { id: 16, code: 'S016', name: 'Lung Surgery' },
  { id: 17, code: 'S017', name: 'Brain Surgery' },
  { id: 18, code: 'S018', name: 'Joint Replacement' },
  { id: 19, code: 'S998', name: 'Others' }
];

const SurgicalHistory = ({ patientId }) => {
  const [surgeries, setSurgeries] = useState([]);

  const [selectedSurgeries, setSelectedSurgeries] = useState([]);
  const [noneChecked, setNoneChecked] = useState(false);
  const [notes, setNotes] = useState('');
  const [currentDateTime, setCurrentDateTime] = useState(new Date());

  // Load saved selections on mount
  useEffect(() => {
    try {
      const raw = localStorage.getItem('surgicalHistorySelections');
      if (raw) {
        const parsed = JSON.parse(raw);
        setSelectedSurgeries(parsed.selectedSurgeries || []);
        setNoneChecked(parsed.noneChecked || false);
        setNotes(parsed.notes || '');
      }
    } catch (e) {}
  }, []);

  // Fetch surgical library from backend on mount
  useEffect(() => {
    let mounted = true;
    surgicalService.getSurgicalLibrary()
      .then((data) => {
        if (!mounted) return;
        const mapped = (data || []).map((o, idx) => ({
          id: o.id || o.SURGERY_CODE || o.surgery_code || `${idx}`,
          code: o.SURGERY_CODE || o.code || o.surgery_code || o.id || '',
          name: o.SURGERY_DESC || o.desc || o.description || o.name || ''
        }));
        setSurgeries(mapped.length ? mapped : defaultSurgeries);
      })
      .catch((err) => {
        console.error('Failed to load surgical library, using defaults', err);
        setSurgeries(defaultSurgeries);
      });
    return () => { mounted = false; };
  }, []);

  // If patientId provided, try loading saved surgical history from backend
  useEffect(() => {
    if (!patientId) return;
    let mounted = true;
    surgicalService.getPatientSurgicalHistory(patientId)
      .then((data) => {
        if (!mounted) return;
        // accept either { selectedSurgeries: [...] } or direct array
        const sel = data && data.selectedSurgeries ? data.selectedSurgeries : data;
        if (Array.isArray(sel)) {
          // normalize server shape into { id, code, name }
          const mapped = sel.map((it, idx) => ({
            id: it.SurgeryCode || it.surgery_code || it.SURGERY_CODE || it.id || `${idx}`,
            code: it.SurgeryCode || it.surgery_code || it.SURGERY_CODE || it.code || it.id || `${idx}`,
            name: it.SurgeryName || it.surgery_name || it.SURGERY_DESC || it.SURGERY_DESC || it.name || ''
          }));
          setSelectedSurgeries(mapped);
        }
      })
      .catch((err) => {
        console.warn('No patient surgical history from server, falling back to storage', err);
        try {
          const raw = localStorage.getItem('surgicalHistorySelections');
          if (raw) {
            const parsed = JSON.parse(raw);
            setSelectedSurgeries(parsed.selectedSurgeries || []);
            setNoneChecked(parsed.noneChecked || false);
            setNotes(parsed.notes || '');
          }
        } catch (e) {}
      });
    return () => { mounted = false; };
  }, [patientId]);

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentDateTime(new Date());
    }, 1000);
    return () => clearInterval(timer);
  }, []);

  const handleCheckboxChange = (surgery) => {
    if (noneChecked) {
      setNoneChecked(false);
    }

    let next;
    const code = surgery.code || surgery.SURGERY_CODE || surgery.id;
    if (selectedSurgeries.find(s => (s.code || s.id) === code)) {
      next = selectedSurgeries.filter(s => (s.code || s.id) !== code);
    } else {
      // store normalized object
      const entry = { id: code, code: code, name: surgery.name || surgery.SURGERY_DESC || '' };
      next = [...selectedSurgeries, entry];
    }
    setSelectedSurgeries(next);
    saveSelections(next, noneChecked, notes);
  };

  const handleNoneChange = () => {
    const nextNone = !noneChecked;
    setNoneChecked(nextNone);
    const nextSelected = nextNone ? [] : selectedSurgeries;
    setSelectedSurgeries(nextSelected);
    saveSelections(nextSelected, nextNone, notes);
  };

  const handleAddNote = () => {
    if (notes.trim()) {
      const noteEntry = {
        id: `note-${Date.now()}`,
        code: 'S998',
        name: notes.trim()
      };
      const next = [...selectedSurgeries, noteEntry];
      setSelectedSurgeries(next);
      setNotes('');
      saveSelections(next, noneChecked, '');
    }
  };

  const handleSave = () => {
    saveSelections(selectedSurgeries, noneChecked, notes);
    alert('Surgical history saved successfully!');
  };

  const handleClear = () => {
    const next = [];
    setSelectedSurgeries(next);
    setNoneChecked(false);
    setNotes('');
    saveSelections(next, false, '');
  };

  const saveSelections = (selected, none, noteText) => {
    try {
      // persist normalized payload to localStorage
      const payload = { selectedSurgeries: selected, noneChecked: none, notes: noteText };
      localStorage.setItem('surgicalHistorySelections', JSON.stringify(payload));
      // Dispatch a custom event so other components in the same window can react immediately
      try {
        const ev = new CustomEvent('surgicalHistoryUpdated', { detail: payload });
        window.dispatchEvent(ev);
      } catch (e) {}
      // If we have a patientId, persist to backend as well (best-effort)
      if (patientId) {
        // backend expects an array of SurgicalHistoryItem
        const serverItems = (selected || []).map(s => ({
          SurgeryCode: s.code || s.id,
          SurgeryName: s.name || s.SURGERY_DESC || s.SURGERY_DESC || '',
          Notes: noteText || '',
          IsChecked: true,
        }));
        surgicalService.savePatientSurgicalHistory(patientId, serverItems).catch((err) => {
          console.error('Failed to save surgical history to server', err);
        });
      }
    } catch (e) {}
  };

  const formatDateTime = (date) => {
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const year = date.getFullYear();
    return `${month}/${day}/${year}`;
  };

  return (
    <div className="surgical-history-container">
      <div className="surgical-history-content">
        {/* Left Panel - Checkboxes */}
        <div className="surgical-history-left">
          <div className="surgical-header">
            <h3>Surgical History Specifics</h3>
          </div>
          
          <div className="checkbox-list-surgical">
            <div className="checkbox-item-surgical">
              <label>
                <input
                  type="checkbox"
                  checked={noneChecked}
                  onChange={handleNoneChange}
                />
                <span>None</span>
              </label>
            </div>

            {surgeries.map(surgery => (
              <div key={surgery.id || surgery.code} className="checkbox-item-surgical">
                <label>
                  <input
                    type="checkbox"
                    checked={selectedSurgeries.some(s => (s.code || s.id) === (surgery.code || surgery.SURGERY_CODE || surgery.id))}
                    onChange={() => handleCheckboxChange(surgery)}
                    disabled={noneChecked}
                  />
                  <span>{surgery.name}</span>
                </label>
              </div>
            ))}
          </div>
        </div>

        {/* Right Panel - Table and Details */}
        <div className="surgical-history-right">
          <div className="notes-section-surgical">
            <textarea
              className="notes-input-surgical"
              placeholder="Add surgical notes or procedures..."
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              rows="3"
            />
            <span className="char-limit-surgical">Limit your characters to 2000</span>
          </div>

          <div className="action-buttons-top-surgical">
            <button className="btn-add-surgical" onClick={handleAddNote} title="Add">
              Add
            </button>
            <button className="btn-save-surgical" onClick={handleSave} title="Save">
              Save
            </button>
            <button className="btn-clear-surgical" onClick={handleClear} title="Clear">
              Clear
            </button>
          </div>

          <div className="table-section-surgical">
            <table className="surgical-table">
              <thead>
                <tr>
                  <th>Code</th>
                  <th>Description</th>
                </tr>
              </thead>
              <tbody>
                {selectedSurgeries.length > 0 ? (
                  selectedSurgeries.map((surgery, index) => (
                    <tr key={`${surgery.id}-${index}`}>
                      <td>{surgery.code}</td>
                      <td>{surgery.name}</td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan="2" className="no-data-cell-surgical">No items selected</td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* DateTime Footer removed */}
    </div>
  );
};

export default SurgicalHistory;