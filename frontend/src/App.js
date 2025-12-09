import React, { useState, useEffect } from 'react';
import './App.css';
import axios from 'axios';

function App() {
  const [items, setItems] = useState([]);
  const [newItem, setNewItem] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [editingId, setEditingId] = useState(null);
  const [editingName, setEditingName] = useState('');

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
      setError('Failed to fetch items: ' + err.message);
      console.error('Error fetching items:', err);
    } finally {
      setLoading(false);
    }
  };

  const addItem = async (e) => {
    e.preventDefault();
    if (!newItem.trim()) return;

    try {
      setLoading(true);
      const response = await axios.post('/api/items', { name: newItem });
      setItems([...items, response.data]);
      setNewItem('');
      setError('');
    } catch (err) {
      setError('Failed to add item: ' + err.message);
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
    setEditingName(item.name);
  };

  const cancelEditing = () => {
    setEditingId(null);
    setEditingName('');
  };

  const updateItem = async (id) => {
    if (!editingName.trim()) return;

    try {
      setLoading(true);
      const response = await axios.put(`/api/items/${id}`, { name: editingName });
      setItems(items.map(item => item.id === id ? response.data : item));
      setEditingId(null);
      setEditingName('');
      setError('');
    } catch (err) {
      setError('Failed to update item: ' + err.message);
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
          <input
            type="text"
            value={newItem}
            onChange={(e) => setNewItem(e.target.value)}
            placeholder="Enter new item..."
            disabled={loading}
          />
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
                      <input
                        type="text"
                        value={editingName}
                        onChange={(e) => setEditingName(e.target.value)}
                        disabled={loading}
                        className="edit-input"
                      />
                      <div className="button-group">
                        <button 
                          onClick={() => updateItem(item.id)}
                          disabled={loading}
                          className="save-btn"
                        >
                          Save
                        </button>
                        <button 
                          onClick={cancelEditing}
                          disabled={loading}
                          className="cancel-btn"
                        >
                          Cancel
                        </button>
                      </div>
                    </>
                  ) : (
                    <>
                      <span>{item.name}</span>
                      <div className="button-group">
                        <button 
                          onClick={() => startEditing(item)}
                          disabled={loading}
                          className="edit-btn"
                        >
                          Edit
                        </button>
                        <button 
                          onClick={() => deleteItem(item.id)}
                          disabled={loading}
                          className="delete-btn"
                        >
                          Delete
                        </button>
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
