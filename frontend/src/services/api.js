import axios from 'axios';

const API_BASE_URL = '/api';

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

// Pertinent Physical Exam admin CRUD service (operate by patno)
export const pertinentService = {
  // List recent entries
  list: async () => {
    const response = await axios.get(`${API_BASE_URL}/pertinent-physical-exams`);
    return response.data;
  },
  // Get a single record by patno
  getByPatno: async (patno) => {
    const response = await axios.get(`${API_BASE_URL}/pertinent-physical-exam/${encodeURIComponent(patno)}`);
    return response.data;
  },
  // Update / upsert a record by patno
  updateByPatno: async (patno, payload) => {
    const response = await axios.put(`${API_BASE_URL}/pertinent-physical-exam/${encodeURIComponent(patno)}`, payload);
    return response.data;
  },
  // Delete a record by patno
  deleteByPatno: async (patno) => {
    const response = await axios.delete(`${API_BASE_URL}/pertinent-physical-exam/${encodeURIComponent(patno)}`);
    return response.data;
  }
};
