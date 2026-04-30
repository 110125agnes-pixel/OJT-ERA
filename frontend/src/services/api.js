import axios from 'axios';

// In development when running on localhost, call backend directly to avoid proxy issues.
// In production, requests should go to the same origin under `/api`.
const API_BASE_URL = (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1')
  ? 'http://localhost:8080/api'
  : '/api';


// Auth API service
export const authService = {
  login: async (payload) => {
    const response = await axios.post(`${API_BASE_URL}/auth/login`, payload);
    return response.data;
  },
  signup: async (payload) => {
    const response = await axios.post(`${API_BASE_URL}/auth/signup`, payload);
    return response.data;
  },
};

// Item API service
export const itemService = {
  // Get all items
  getAllItems: async () => {
    const response = await axios.get(`${API_BASE_URL}/items`);
    return response.data;
  },

  // Get a single item by ID
  getItem: async (id) => {
    const response = await axios.get(`${API_BASE_URL}/items/${id}`);
    return response.data;
  },

  // Create a new item
  createItem: async (item) => {
    const response = await axios.post(`${API_BASE_URL}/items`, item);
    return response.data;
  },

  // Update an existing item
  updateItem: async (id, item) => {
    const response = await axios.put(`${API_BASE_URL}/items/${id}`, item);
    return response.data;
  },

  // Delete an item
  deleteItem: async (id) => {
    const response = await axios.delete(`${API_BASE_URL}/items/${id}`);
    return response.data;
  }
};

// Inventory API service
export const inventoryService = {
  // Get all inventory items
  getAllInventory: async () => {
    const response = await axios.get(`${API_BASE_URL}/inventory`);
    return response.data;
  },

  // Create a new inventory item
  createInventoryItem: async (item) => {
    const response = await axios.post(`${API_BASE_URL}/inventory`, item);
    return response.data;
  },

  // Update an existing inventory item
  updateInventoryItem: async (id, item) => {
    const response = await axios.put(`${API_BASE_URL}/inventory/${id}`, item);
    return response.data;
  },

  // Delete an inventory item
  deleteInventoryItem: async (id) => {
    const response = await axios.delete(`${API_BASE_URL}/inventory/${id}`);
    return response.data;
  }
};

// Surgical history API service
export const surgicalService = {
  // Fetch surgical library (list of known surgeries)
  getSurgicalLibrary: async () => {
    const response = await axios.get(`${API_BASE_URL}/lib/surgery`);
    return response.data;
  },

  // Get saved surgical history for a patient
  getPatientSurgicalHistory: async (patientId) => {
    const response = await axios.get(`${API_BASE_URL}/patients/${patientId}/surgical-history`);
    return response.data;
  },

  // Save surgical history for a patient
  savePatientSurgicalHistory: async (patientId, payload) => {
    const response = await axios.post(`${API_BASE_URL}/patients/${patientId}/surgical-history`, payload);
    return response.data;
  }
};
export const libService = {
  // Get digital rectal library options
  getDigitalRectal: async () => {
    const response = await axios.get(`${API_BASE_URL}/lib/digital_rectal`);
    return response.data;
  },
  // Get genitourinary library options
  getGenitourinary: async () => {
    const response = await axios.get(`${API_BASE_URL}/lib/genitourinary`);
    return response.data;
  },
};
