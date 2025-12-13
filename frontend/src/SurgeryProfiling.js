import React, { useState, useEffect } from 'react';
import './SurgeryProfiling.css';
import { surgeryService } from './services/api';

function SurgeryProfiling() {
  const [surgeries, setSurgeries] = useState([]);
  const [newSurgery, setNewSurgery] = useState({
    patient_name: '',
    surgery_type: '',
    surgeon_name: '',
    surgery_date: '',
    surgery_time: '',
    duration: '',
    status: '',
    notes: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [editingId, setEditingId] = useState(null);
  const [editingSurgery, setEditingSurgery] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');

  const surgeryTypes = ['', 'Orthopedic', 'Cardiac', 'Neurosurgery', 'General', 'Plastic', 'Vascular', 'Thoracic', 'Pediatric', 'Other'];
  const statusOptions = ['', 'Scheduled', 'In Progress', 'Completed', 'Cancelled', 'Postponed'];

  useEffect(() => {
    fetchSurgeries();
  }, []);

  const fetchSurgeries = async () => {
    try {
      setLoading(true);
      const data = await surgeryService.getAllSurgeries();
      setSurgeries(data || []);
      setError('');
    } catch (err) {
      setError('Failed to fetch surgeries: ' + (err.response?.data?.error || err.message));
      console.error('Error fetching surgeries:', err);
    } finally {
      setLoading(false);
    }
  };

  const addSurgery = async (e) => {
    e.preventDefault();
    if (!newSurgery.patient_name.trim() || !newSurgery.surgery_type.trim() || !newSurgery.surgeon_name.trim()) {
      setError('Patient Name, Surgery Type, and Surgeon Name are required');
      return;
    }
    try {
      setLoading(true);
      const data = await surgeryService.createSurgery(newSurgery);
      setSurgeries([...surgeries, data]);
      setNewSurgery({
        patient_name: '',
        surgery_type: '',
        surgeon_name: '',
        surgery_date: '',
        surgery_time: '',
        duration: '',
        status: '',
        notes: ''
      });
      setError('');
    } catch (err) {
      setError('Failed to add surgery: ' + (err.response?.data?.error || err.message));
      console.error('Error adding surgery:', err);
    } finally {
      setLoading(false);
    }
  };

  const deleteSurgery = async (id) => {
    const confirmDelete = window.confirm('Are you sure you want to delete this surgery record?');
    if (!confirmDelete) return;

    try {
      setLoading(true);
      await surgeryService.deleteSurgery(id);
      setSurgeries(surgeries.filter(surgery => surgery.id !== id));
      setError('');
    } catch (err) {
      setError('Failed to delete surgery: ' + err.message);
      console.error('Error deleting surgery:', err);
    } finally {
      setLoading(false);
    }
  };

  const startEditing = (surgery) => {
    setEditingId(surgery.id);
    setEditingSurgery({ ...surgery });
  };

  const cancelEditing = () => {
    setEditingId(null);
    setEditingSurgery(null);
  };

  const updateSurgery = async (id) => {
    if (!editingSurgery.patient_name.trim() || !editingSurgery.surgery_type.trim() || !editingSurgery.surgeon_name.trim()) {
      setError('Patient Name, Surgery Type, and Surgeon Name are required');
      return;
    }
    try {
      setLoading(true);
      const data = await surgeryService.updateSurgery(id, editingSurgery);
      setSurgeries(surgeries.map(surgery => surgery.id === id ? data : surgery));
      setEditingId(null);
      setEditingSurgery(null);
      setError('');
    } catch (err) {
      setError('Failed to update surgery: ' + (err.response?.data?.error || err.message));
      console.error('Error updating surgery:', err);
    } finally {
      setLoading(false);
    }
  };

  const filteredSurgeries = surgeries.filter(surgery =>
    surgery.patient_name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    surgery.surgeon_name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    surgery.surgery_type?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="surgery-container">
      <h2>🏥 Surgery Management</h2>
      <p className="subtitle">Schedule and manage surgery operations</p>
      
      {error && <div className="error">{error}</div>}
      
      <form onSubmit={addSurgery} className="add-form">
        <h3>Schedule New Surgery</h3>
        <div className="form-grid">
          <input
            type="text"
            placeholder="Patient Name *"
            value={newSurgery.patient_name}
            onChange={(e) => setNewSurgery({...newSurgery, patient_name: e.target.value})}
            required
          />
          <select
            value={newSurgery.surgery_type}
            onChange={(e) => setNewSurgery({...newSurgery, surgery_type: e.target.value})}
            required
          >
            {surgeryTypes.map(type => (
              <option key={type} value={type}>{type || 'Select Surgery Type *'}</option>
            ))}
          </select>
          <input
            type="text"
            placeholder="Surgeon Name *"
            value={newSurgery.surgeon_name}
            onChange={(e) => setNewSurgery({...newSurgery, surgeon_name: e.target.value})}
            required
          />
          <input
            type="date"
            placeholder="Surgery Date"
            value={newSurgery.surgery_date}
            onChange={(e) => setNewSurgery({...newSurgery, surgery_date: e.target.value})}
          />
          <input
            type="time"
            placeholder="Surgery Time"
            value={newSurgery.surgery_time}
            onChange={(e) => setNewSurgery({...newSurgery, surgery_time: e.target.value})}
          />
          <input
            type="text"
            placeholder="Duration (e.g., 2 hours)"
            value={newSurgery.duration}
            onChange={(e) => setNewSurgery({...newSurgery, duration: e.target.value})}
          />
          <select
            value={newSurgery.status}
            onChange={(e) => setNewSurgery({...newSurgery, status: e.target.value})}
          >
            {statusOptions.map(status => (
              <option key={status} value={status}>{status || 'Select Status'}</option>
            ))}
          </select>
          <input
            type="text"
            placeholder="Notes"
            value={newSurgery.notes}
            onChange={(e) => setNewSurgery({...newSurgery, notes: e.target.value})}
            className="notes-input"
          />
        </div>
        <button type="submit" disabled={loading} className="submit-btn">
          {loading ? '⏳ Adding...' : '➕ Schedule Surgery'}
        </button>
      </form>

      <div className="surgery-list-section">
        <div className="list-header">
          <h3>Surgery Records ({filteredSurgeries.length})</h3>
          <input
            type="text"
            placeholder="🔍 Search by patient, surgeon, or type..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="search-input"
          />
        </div>

        {loading && <div className="loading">Loading surgeries...</div>}
        
        {!loading && filteredSurgeries.length === 0 && (
          <div className="no-data">No surgery records found. Schedule your first surgery above.</div>
        )}

        {!loading && filteredSurgeries.length > 0 && (
          <div className="table-wrapper">
            <table className="surgery-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Patient Name</th>
                  <th>Surgery Type</th>
                  <th>Surgeon</th>
                  <th>Date</th>
                  <th>Time</th>
                  <th>Duration</th>
                  <th>Status</th>
                  <th>Notes</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredSurgeries.map((surgery) => (
                  <tr key={surgery.id}>
                    {editingId === surgery.id ? (
                      <>
                        <td>{surgery.id}</td>
                        <td>
                          <input
                            type="text"
                            value={editingSurgery.patient_name}
                            onChange={(e) => setEditingSurgery({...editingSurgery, patient_name: e.target.value})}
                            className="edit-input"
                          />
                        </td>
                        <td>
                          <select
                            value={editingSurgery.surgery_type}
                            onChange={(e) => setEditingSurgery({...editingSurgery, surgery_type: e.target.value})}
                            className="edit-input"
                          >
                            {surgeryTypes.map(type => (
                              <option key={type} value={type}>{type || 'Select Type'}</option>
                            ))}
                          </select>
                        </td>
                        <td>
                          <input
                            type="text"
                            value={editingSurgery.surgeon_name}
                            onChange={(e) => setEditingSurgery({...editingSurgery, surgeon_name: e.target.value})}
                            className="edit-input"
                          />
                        </td>
                        <td>
                          <input
                            type="date"
                            value={editingSurgery.surgery_date}
                            onChange={(e) => setEditingSurgery({...editingSurgery, surgery_date: e.target.value})}
                            className="edit-input"
                          />
                        </td>
                        <td>
                          <input
                            type="time"
                            value={editingSurgery.surgery_time}
                            onChange={(e) => setEditingSurgery({...editingSurgery, surgery_time: e.target.value})}
                            className="edit-input"
                          />
                        </td>
                        <td>
                          <input
                            type="text"
                            value={editingSurgery.duration}
                            onChange={(e) => setEditingSurgery({...editingSurgery, duration: e.target.value})}
                            className="edit-input"
                          />
                        </td>
                        <td>
                          <select
                            value={editingSurgery.status}
                            onChange={(e) => setEditingSurgery({...editingSurgery, status: e.target.value})}
                            className="edit-input"
                          >
                            {statusOptions.map(status => (
                              <option key={status} value={status}>{status || 'Select Status'}</option>
                            ))}
                          </select>
                        </td>
                        <td>
                          <input
                            type="text"
                            value={editingSurgery.notes}
                            onChange={(e) => setEditingSurgery({...editingSurgery, notes: e.target.value})}
                            className="edit-input"
                          />
                        </td>
                        <td>
                          <button onClick={() => updateSurgery(surgery.id)} className="save-btn">💾</button>
                          <button onClick={cancelEditing} className="cancel-btn">❌</button>
                        </td>
                      </>
                    ) : (
                      <>
                        <td>{surgery.id}</td>
                        <td>{surgery.patient_name}</td>
                        <td>
                          <span className="badge surgery-type">{surgery.surgery_type}</span>
                        </td>
                        <td>{surgery.surgeon_name}</td>
                        <td>{surgery.surgery_date || 'N/A'}</td>
                        <td>{surgery.surgery_time || 'N/A'}</td>
                        <td>{surgery.duration || 'N/A'}</td>
                        <td>
                          <span className={`badge status-${surgery.status?.toLowerCase().replace(' ', '-')}`}>
                            {surgery.status || 'N/A'}
                          </span>
                        </td>
                        <td>{surgery.notes || '-'}</td>
                        <td>
                          <button onClick={() => startEditing(surgery)} className="edit-btn">✏️</button>
                          <button onClick={() => deleteSurgery(surgery.id)} className="delete-btn">🗑️</button>
                        </td>
                      </>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}

export default SurgeryProfiling;
