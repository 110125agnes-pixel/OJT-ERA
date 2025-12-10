
import React, { useState, useEffect } from 'react';
import './App.css';
import axios from 'axios';

function App() {
  const [items, setItems] = useState([]);
  const [newItem, setNewItem] = useState({
    lastname: '',
    firstname: '',
    middlename: '',
    suffix: '',
    birthdate: '',
    sex: '',
    civil_status: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [editingId, setEditingId] = useState(null);
  const [editingItem, setEditingItem] = useState(null);

  const sexOptions = ['', 'Male', 'Female', 'Other'];
  const civilStatusOptions = ['', 'Single', 'Married', 'Divorced', 'Widowed'];

  useEffect(() => {
    fetchItems();
  }, []);

  const fetchItems = async () => {
    try {
      setLoading(true);
      const response = await axios.get('/api/items');
      setItems(response.data || []);
      setError('');
    } catch (err) {
      setError('Failed to fetch items: ' + (err.response?.data?.error || err.message));
      console.error('Error fetching items:', err);
    } finally {
      setLoading(false);
    }
  };

  const addItem = async (e) => {
    e.preventDefault();
    if (!newItem.lastname.trim() || !newItem.firstname.trim()) {
      setError('Lastname and Firstname are required');
      return;
    }
    try {
      setLoading(true);
      const response = await axios.post('/api/items', newItem);
      setItems([...items, response.data]);
      setNewItem({
        lastname: '',
        firstname: '',
        middlename: '',
        suffix: '',
        birthdate: '',
        sex: '',
        civil_status: ''
      });
      setError('');
    } catch (err) {
      setError('Failed to add item: ' + (err.response?.data?.error || err.message));
      console.error('Error adding item:', err);
    } finally {
      setLoading(false);
    }
  };

  const deleteItem = async (id) => {
    const confirmDelete = window.confirm('Are you sure you want to delete this item?');
    if (!confirmDelete) return;

    try {
      setLoading(true);
      await axios.delete(`/api/items/${id}`);
      setItems(items.filter(item => item.id !== id));
      setError('');
    } catch (err) {
      setError('Failed to delete item: ' + err.message);
      console.error('Error deleting item:', err);
    } finally {
      setLoading(false);
    }
  };

  const startEditing = (item) => {
    setEditingId(item.id);
    setEditingItem({ ...item });
  };

  const cancelEditing = () => {
    setEditingId(null);
    setEditingItem(null);
  };

  const updateItem = async (id) => {
    if (!editingItem.lastname.trim() || !editingItem.firstname.trim()) {
      setError('Lastname and Firstname are required');
      return;
    }
    try {
      setLoading(true);
      const response = await axios.put(`/api/items/${id}`, editingItem);
      setItems(items.map(item => item.id === id ? response.data : item));
      setEditingId(null);
      setEditingItem(null);
      setError('');
    } catch (err) {
      setError('Failed to update item: ' + (err.response?.data?.error || err.message));
      console.error('Error updating item:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>React + Go + MySQL</h1>
        <p>Full Stack Application</p>
      </header>
      
      <main className="App-main">
        {error && <div className="error">{error}</div>}
        
        <form onSubmit={addItem} className="add-form">
          <div className="form-row">
            <input type="text" value={newItem.lastname} onChange={e => setNewItem({ ...newItem, lastname: e.target.value })} placeholder="Lastname" disabled={loading} required />
            <input type="text" value={newItem.firstname} onChange={e => setNewItem({ ...newItem, firstname: e.target.value })} placeholder="Firstname" disabled={loading} required />
            <input type="text" value={newItem.middlename} onChange={e => setNewItem({ ...newItem, middlename: e.target.value })} placeholder="Middlename" disabled={loading} />
            <input type="text" value={newItem.suffix} onChange={e => setNewItem({ ...newItem, suffix: e.target.value })} placeholder="Suffix" disabled={loading} />
          </div>
          <div className="form-row">
            <input type="date" value={newItem.birthdate} onChange={e => setNewItem({ ...newItem, birthdate: e.target.value })} placeholder="Birthdate" disabled={loading} />
            <select value={newItem.sex} onChange={e => setNewItem({ ...newItem, sex: e.target.value })} disabled={loading} required>
              {sexOptions.map(opt => <option key={opt} value={opt}>{opt || 'Sex'}</option>)}
            </select>
            <select value={newItem.civil_status} onChange={e => setNewItem({ ...newItem, civil_status: e.target.value })} disabled={loading} required>
              {civilStatusOptions.map(opt => <option key={opt} value={opt}>{opt || 'Civil Status'}</option>)}
            </select>
          </div>
          <button type="submit" disabled={loading}>
            {loading ? 'Adding...' : 'Add Item'}
          </button>
        </form>

        <div className="items-list">
          <h2>Items</h2>
          {loading && items.length === 0 ? (
            <p>Loading...</p>
          ) : items.length === 0 ? (
            <p>No items yet. Add one above!</p>
          ) : (
            <ul>
              {items.map((item) => (
                <li key={item.id}>
                  {editingId === item.id ? (
                    <>
                      <div className="form-row">
                        <input type="text" value={editingItem.lastname} onChange={e => setEditingItem({ ...editingItem, lastname: e.target.value })} placeholder="Lastname" disabled={loading} required />
                        <input type="text" value={editingItem.firstname} onChange={e => setEditingItem({ ...editingItem, firstname: e.target.value })} placeholder="Firstname" disabled={loading} required />
                        <input type="text" value={editingItem.middlename} onChange={e => setEditingItem({ ...editingItem, middlename: e.target.value })} placeholder="Middlename" disabled={loading} />
                        <input type="text" value={editingItem.suffix} onChange={e => setEditingItem({ ...editingItem, suffix: e.target.value })} placeholder="Suffix" disabled={loading} />
                      </div>
                      <div className="form-row">
                        <input type="date" value={editingItem.birthdate} onChange={e => setEditingItem({ ...editingItem, birthdate: e.target.value })} placeholder="Birthdate" disabled={loading} />
                        <select value={editingItem.sex} onChange={e => setEditingItem({ ...editingItem, sex: e.target.value })} disabled={loading} required>
                          {sexOptions.map(opt => <option key={opt} value={opt}>{opt || 'Sex'}</option>)}
                        </select>
                        <select value={editingItem.civil_status} onChange={e => setEditingItem({ ...editingItem, civil_status: e.target.value })} disabled={loading} required>
                          {civilStatusOptions.map(opt => <option key={opt} value={opt}>{opt || 'Civil Status'}</option>)}
                        </select>
                      </div>
                      <div className="button-group">
                        <button onClick={() => updateItem(item.id)} disabled={loading} className="save-btn">Save</button>
                        <button onClick={cancelEditing} disabled={loading} className="cancel-btn">Cancel</button>
                      </div>
                    </>
                  ) : (
                    <>
                      <div className="item-fields">
                        <span><b>{item.lastname}, {item.firstname} {item.middlename} {item.suffix}</b></span><br />
                        <span>Birthdate: {item.birthdate || '-'}</span> | <span>Sex: {item.sex || '-'}</span> | <span>Civil Status: {item.civil_status || '-'}</span>
                      </div>
                      <div className="button-group">
                        <button onClick={() => startEditing(item)} disabled={loading} className="edit-btn">Edit</button>
                        <button onClick={() => deleteItem(item.id)} disabled={loading} className="delete-btn">Delete</button>
                      </div>
                    </>
                  )}
                </li>
              ))}
            </ul>
          )}
        </div>
      </main>
    </div>
  );
}

export default App;
