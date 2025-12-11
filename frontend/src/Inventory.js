import React, { useState, useEffect } from 'react';
import axios from 'axios';
import './Inventory.css';

function Inventory() {
  const [items, setItems] = useState([]);
  const [newItem, setNewItem] = useState({
    itemName: '',
    category: '',
    brand: '',
    quantity: '',
    unit: '',
    price: ''
  });
  const [editingId, setEditingId] = useState(null);
  const [editingItem, setEditingItem] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');

  // Fetch inventory items on component mount
  useEffect(() => {
    fetchItems();
  }, []);

  const fetchItems = async () => {
    try {
      const response = await axios.get('http://localhost:8080/api/inventory');
      // Map backend field names to frontend format
      const mappedItems = response.data.map(item => ({
        id: item.id,
        itemName: item.item_name,
        category: item.category,
        brand: item.brand,
        quantity: item.quantity,
        unit: item.unit,
        price: item.price
      }));
      setItems(mappedItems);
    } catch (error) {
      console.error('Error fetching items:', error);
      alert('Failed to load inventory items');
    }
  };

  const categoryOptions = [
    '', 
    'Motherboard', 
    'Processor (CPU)', 
    'RAM (Memory)', 
    'Graphics Card (GPU)',
    'Storage (HDD/SSD)',
    'Power Supply (PSU)',
    'Computer Case',
    'Cooling System',
    'Monitor',
    'Keyboard',
    'Mouse',
    'Headset',
    'Webcam',
    'Speakers',
    'Cables & Accessories',
    'Networking',
    'Other'
  ];
  const brandOptions = [
    '',
    'Intel',
    'AMD',
    'ASUS',
    'MSI',
    'Gigabyte',
    'NVIDIA',
    'Corsair',
    'Kingston',
    'Samsung',
    'Western Digital',
    'Seagate',
    'Logitech',
    'Razer',
    'HyperX',
    'Cooler Master',
    'NZXT',
    'Thermaltake',
    'Acer',
    'Dell',
    'LG',
    'ViewSonic',
    'Other'
  ];
  const unitOptions = ['', 'pcs', 'unit', 'set', 'box'];

  const addItem = async (e) => {
    e.preventDefault();
    if (!newItem.itemName.trim() || !newItem.quantity) {
      alert('Item name and quantity are required');
      return;
    }
    
    try {
      // Send to backend with correct field names
      const response = await axios.post('http://localhost:8080/api/inventory', {
        item_name: newItem.itemName,
        category: newItem.category,
        brand: newItem.brand,
        quantity: parseInt(newItem.quantity) || 0,
        unit: newItem.unit,
        price: parseFloat(newItem.price) || 0
      });
      
      // Map response back to frontend format
      const addedItem = {
        id: response.data.id,
        itemName: response.data.item_name,
        category: response.data.category,
        brand: response.data.brand,
        quantity: response.data.quantity,
        unit: response.data.unit,
        price: response.data.price
      };
      
      setItems([addedItem, ...items]);
      setNewItem({
        itemName: '',
        category: '',
        brand: '',
        quantity: '',
        unit: '',
        price: ''
      });
    } catch (error) {
      console.error('Error adding item:', error);
      alert('Failed to add item');
    }
  };

  const deleteItem = async (id) => {
    if (window.confirm('Are you sure you want to delete this item?')) {
      try {
        await axios.delete(`http://localhost:8080/api/inventory/${id}`);
        setItems(items.filter(item => item.id !== id));
      } catch (error) {
        console.error('Error deleting item:', error);
        alert('Failed to delete item');
      }
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
    if (!editingItem.itemName.trim() || !editingItem.quantity) {
      alert('Item name and quantity are required');
      return;
    }
    
    try {
      const response = await axios.put(`http://localhost:8080/api/inventory/${id}`, {
        item_name: editingItem.itemName,
        category: editingItem.category,
        brand: editingItem.brand,
        quantity: parseInt(editingItem.quantity) || 0,
        unit: editingItem.unit,
        price: parseFloat(editingItem.price) || 0
      });
      
      // Map response back to frontend format
      const updatedItem = {
        id: response.data.id,
        itemName: response.data.item_name,
        category: response.data.category,
        brand: response.data.brand,
        quantity: response.data.quantity,
        unit: response.data.unit,
        price: response.data.price
      };
      
      setItems(items.map(item => item.id === id ? updatedItem : item));
      setEditingId(null);
      setEditingItem(null);
    } catch (error) {
      console.error('Error updating item:', error);
      alert('Failed to update item');
    }
  };

  const filteredItems = items.filter(item =>
    item.itemName?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    item.category?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="inventory-container">
      <h2>📦 Inventory Management</h2>
      <p className="subtitle">Track and manage inventory items</p>
      
      <form onSubmit={addItem} className="add-form">
        <h3>Add New Computer Hardware Item</h3>
        <div className="form-grid">
          <input
            type="text"
            placeholder="Product Name *"
            value={newItem.itemName}
            onChange={(e) => setNewItem({...newItem, itemName: e.target.value})}
            required
          />
          <select
            value={newItem.category}
            onChange={(e) => setNewItem({...newItem, category: e.target.value})}
            required
          >
            {categoryOptions.map((option, idx) => (
              <option key={idx} value={option}>
                {option || 'Select Category *'}
              </option>
            ))}
          </select>
          <select
            value={newItem.brand}
            onChange={(e) => setNewItem({...newItem, brand: e.target.value})}
          >
            {brandOptions.map((option, idx) => (
              <option key={idx} value={option}>
                {option || 'Select Brand'}
              </option>
            ))}
          </select>
          <input
            type="number"
            placeholder="Quantity *"
            value={newItem.quantity}
            onChange={(e) => setNewItem({...newItem, quantity: e.target.value})}
            min="0"
            required
          />
          <select
            value={newItem.unit}
            onChange={(e) => setNewItem({...newItem, unit: e.target.value})}
          >
            {unitOptions.map((option, idx) => (
              <option key={idx} value={option}>
                {option || 'Select Unit'}
              </option>
            ))}
          </select>
          <input
            type="number"
            placeholder="Price (₱)"
            value={newItem.price}
            onChange={(e) => setNewItem({...newItem, price: e.target.value})}
            min="0"
            step="0.01"
          />
          <button type="submit" className="submit-btn">
            Add Item
          </button>
        </div>
      </form>

      <div className="search-section">
        <input
          type="text"
          placeholder="🔍 Search by name..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="search-input"
        />
      </div>

      <div className="items-section">
        <h3>Inventory Items ({filteredItems.length})</h3>
        {filteredItems.length === 0 ? (
          <p className="empty-state">No inventory items yet. Add your first item above!</p>
        ) : (
          <div className="items-table-container">
            <table className="items-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Product Name</th>
                  <th>Category</th>
                  <th>Brand</th>
                  <th>Quantity</th>
                  <th>Unit</th>
                  <th>Price</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredItems.map(item => (
                  <tr key={item.id}>
                    {editingId === item.id ? (
                      <>
                        <td>{item.id}</td>
                        <td>
                          <input
                            type="text"
                            value={editingItem.itemName}
                            onChange={(e) => setEditingItem({...editingItem, itemName: e.target.value})}
                            placeholder="Product Name"
                          />
                        </td>
                        <td>
                          <select
                            value={editingItem.category}
                            onChange={(e) => setEditingItem({...editingItem, category: e.target.value})}
                          >
                            {categoryOptions.map((opt, idx) => <option key={idx} value={opt}>{opt || 'Category'}</option>)}
                          </select>
                        </td>
                        <td>
                          <select
                            value={editingItem.brand}
                            onChange={(e) => setEditingItem({...editingItem, brand: e.target.value})}
                          >
                            {brandOptions.map((opt, idx) => <option key={idx} value={opt}>{opt || 'Brand'}</option>)}
                          </select>
                        </td>
                        <td>
                          <input
                            type="number"
                            value={editingItem.quantity}
                            onChange={(e) => setEditingItem({...editingItem, quantity: e.target.value})}
                            placeholder="Quantity"
                            min="0"
                          />
                        </td>
                        <td>
                          <select
                            value={editingItem.unit}
                            onChange={(e) => setEditingItem({...editingItem, unit: e.target.value})}
                          >
                            {unitOptions.map((opt, idx) => <option key={idx} value={opt}>{opt || 'Unit'}</option>)}
                          </select>
                        </td>
                        <td>
                          <input
                            type="number"
                            value={editingItem.price}
                            onChange={(e) => setEditingItem({...editingItem, price: e.target.value})}
                            placeholder="Price"
                            min="0"
                            step="0.01"
                          />
                        </td>
                        <td>
                          <button onClick={() => updateItem(item.id)} className="save-btn">Save</button>
                          <button onClick={cancelEditing} className="cancel-btn">Cancel</button>
                        </td>
                      </>
                    ) : (
                      <>
                        <td>{item.id}</td>
                        <td>{item.itemName}</td>
                        <td>{item.category || '-'}</td>
                        <td>{item.brand || '-'}</td>
                        <td>{item.quantity}</td>
                        <td>{item.unit || '-'}</td>
                        <td>₱{item.price ? item.price.toLocaleString('en-PH', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) : '0.00'}</td>
                        <td>
                          <button onClick={() => startEditing(item)} className="edit-btn">Edit</button>
                          <button onClick={() => deleteItem(item.id)} className="delete-btn">Delete</button>
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

export default Inventory;
