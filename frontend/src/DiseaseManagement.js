import React, { useState, useEffect } from 'react';
import { diseaseService } from './services/api';
import './DiseaseManagement.css';

function DiseaseManagement() {
  const [diseases, setDiseases] = useState([]);
  const [selectedDiseases, setSelectedDiseases] = useState(new Set());
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [currentTime, setCurrentTime] = useState(new Date());
  const [showAddForm, setShowAddForm] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [formData, setFormData] = useState({
    name: '',
    code: '',
    barcode: '',
    category: ''
  });

  const categories = ['', 'Infectious', 'Chronic', 'Mental Health', 'Other'];

  useEffect(() => {
    fetchDiseases();
    const timer = setInterval(() => setCurrentTime(new Date()), 1000);
    return () => clearInterval(timer);
  }, []);

  const fetchDiseases = async () => {
    try {
      setLoading(true);
      const data = await diseaseService.getAllDiseases();
      setDiseases(data || []);
      setError('');
    } catch (err) {
      setError('Failed to fetch diseases: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  const toggleDisease = (id) => {
    const newSelected = new Set(selectedDiseases);
    if (newSelected.has(id)) {
      newSelected.delete(id);
    } else {
      newSelected.add(id);
    }
    setSelectedDiseases(newSelected);
  };

  const handleAddDisease = async (e) => {
    e.preventDefault();
    if (!formData.name.trim() || !formData.code.trim() || !formData.barcode.trim()) {
      setError('Name, Code, and Barcode are required');
      return;
    }
    try {
      setLoading(true);
      if (editingId) {
        await diseaseService.updateDisease(editingId, formData);
        setDiseases(diseases.map(d => d.id === editingId ? { ...formData, id: editingId } : d));
      } else {
        const data = await diseaseService.createDisease(formData);
        setDiseases([...diseases, data]);
      }
      setFormData({ name: '', code: '', barcode: '', category: '' });
      setShowAddForm(false);
      setEditingId(null);
      setError('');
    } catch (err) {
      setError('Failed to save disease: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = (disease) => {
    setFormData(disease);
    setEditingId(disease.id);
    setShowAddForm(true);
  };

  const handleDelete = async (id) => {
    if (window.confirm('Are you sure you want to delete this disease?')) {
      try {
        setLoading(true);
        await diseaseService.deleteDisease(id);
        setDiseases(diseases.filter(d => d.id !== id));
        setSelectedDiseases(new Set([...selectedDiseases].filter(s => s !== id)));
        setError('');
      } catch (err) {
        setError('Failed to delete disease: ' + err.message);
      } finally {
        setLoading(false);
      }
    }
  };

  const handleCancel = () => {
    setShowAddForm(false);
    setEditingId(null);
    setFormData({ name: '', code: '', barcode: '', category: '' });
  };

  const selectedDiseasesList = diseases.filter(d => selectedDiseases.has(d.id));

  const formatDateTime = () => {
    const date = currentTime.toLocaleDateString('en-US', { month: '2-digit', day: '2-digit', year: 'numeric' });
    const time = currentTime.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
    return `${date}  ${time}`;
  };

  return (
    <div className="disease-management-page">
      <div className="disease-content">
        {error && <div className="error-banner">{error}</div>}

        <div className="disease-layout">
          {/* Left Panel - Checkboxes */}
          <div className="disease-left-panel">
            <div className="panel-header">
              <h3>Medical Conditions</h3>
              <button className="add-new-btn" onClick={() => setShowAddForm(true)}>+ Add New</button>
            </div>

            {showAddForm && (
              <div className="add-disease-form">
                <h4>{editingId ? 'Edit Disease' : 'New Disease'}</h4>
                <input
                  type="text"
                  placeholder="Name"
                  value={formData.name}
                  onChange={(e) => setFormData({...formData, name: e.target.value})}
                  disabled={loading}
                />
                <input
                  type="text"
                  placeholder="Code"
                  value={formData.code}
                  onChange={(e) => setFormData({...formData, code: e.target.value})}
                  disabled={loading}
                />
                <input
                  type="text"
                  placeholder="Barcode"
                  value={formData.barcode}
                  onChange={(e) => setFormData({...formData, barcode: e.target.value})}
                  disabled={loading}
                />
                <select
                  value={formData.category}
                  onChange={(e) => setFormData({...formData, category: e.target.value})}
                  disabled={loading}
                >
                  {categories.map(cat => <option key={cat} value={cat}>{cat || 'Category'}</option>)}
                </select>
                <div className="form-buttons">
                  <button className="save-btn" onClick={handleAddDisease} disabled={loading}>Save</button>
                  <button className="cancel-btn" onClick={handleCancel}>Cancel</button>
                </div>
              </div>
            )}

            <div className="checkbox-list">
              {diseases.map(disease => (
                <div key={disease.id} className="checkbox-item">
                  <label>
                    <input
                      type="checkbox"
                      checked={selectedDiseases.has(disease.id)}
                      onChange={() => toggleDisease(disease.id)}
                    />
                    <span>{disease.name}</span>
                  </label>
                </div>
              ))}
            </div>
          </div>

          {/* Right Panel - Table */}
          <div className="disease-right-panel">
            <div className="panel-header">
              <h3>Disease Database</h3>
            </div>

            {diseases.length === 0 ? (
              <div className="no-data">No diseases in database</div>
            ) : (
              <div className="disease-table-wrapper">
                <table className="disease-table">
                  <thead>
                    <tr>
                      <th>Code</th>
                      <th>Name/Description</th>
                      <th>Barcode</th>
                      <th>Category</th>
                      <th>Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {diseases.map(disease => (
                      <tr key={disease.id}>
                        <td onClick={() => {navigator.clipboard.writeText(disease.code)}} style={{cursor: 'pointer', color: '#667eea', fontWeight: 'bold'}}>
                          {disease.code}
                        </td>
                        <td>{disease.name}</td>
                        <td style={{fontSize: '0.8rem', fontFamily: 'monospace'}}>{disease.barcode}</td>
                        <td>{disease.category || '-'}</td>
                        <td className="actions">
                          <button className="edit-link" onClick={() => handleEdit(disease)}>Edit</button>
                          <button className="delete-link" onClick={() => handleDelete(disease.id)}>Delete</button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Bottom DateTime Bar */}
      <div className="datetime-footer">
        <span>{formatDateTime()}</span>
      </div>
    </div>
  );
}

export default DiseaseManagement;
