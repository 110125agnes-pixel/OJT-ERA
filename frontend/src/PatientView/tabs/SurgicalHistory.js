import React, { useState, useEffect } from 'react';
import './SurgicalHistory.css';

const SurgicalHistory = () => {
  const [surgeries] = useState([
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
  ]);

  const [selectedSurgeries, setSelectedSurgeries] = useState([]);
  const [noneChecked, setNoneChecked] = useState(false);
  const [notes, setNotes] = useState('');
  const [currentDateTime, setCurrentDateTime] = useState(new Date());

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

    if (selectedSurgeries.find(s => s.id === surgery.id)) {
      setSelectedSurgeries(selectedSurgeries.filter(s => s.id !== surgery.id));
    } else {
      setSelectedSurgeries([...selectedSurgeries, surgery]);
    }
  };

  const handleNoneChange = () => {
    setNoneChecked(!noneChecked);
    if (!noneChecked) {
      setSelectedSurgeries([]);
    }
  };

  const handleAddNote = () => {
    if (notes.trim()) {
      const noteEntry = {
        id: Date.now(),
        code: 'S998',
        name: notes.trim()
      };
      setSelectedSurgeries([...selectedSurgeries, noteEntry]);
      setNotes('');
    }
  };

  const handleSave = () => {
    alert('Surgical history saved successfully!');
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
              <div key={surgery.id} className="checkbox-item-surgical">
                <label>
                  <input
                    type="checkbox"
                    checked={selectedSurgeries.some(s => s.id === surgery.id)}
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
            <button className="btn-clear-surgical" onClick={() => setSelectedSurgeries([])} title="Clear">
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