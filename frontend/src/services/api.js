import axios from 'axios';

const API_BASE_URL = '/api';

// Item API service
export const itemService = {
  // Get all items
  getAllItems: async () => {
    const response = await axios.get(`${API_BASE_URL}/items`);
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

// Disease API service
export const diseaseService = {
  // Get all diseases
  getAllDiseases: async () => {
    const response = await axios.get(`${API_BASE_URL}/diseases`);
    return response.data;
  },

  // Create a new disease
  createDisease: async (disease) => {
    const response = await axios.post(`${API_BASE_URL}/diseases`, disease);
    return response.data;
  },

  // Update a disease
  updateDisease: async (id, disease) => {
    const response = await axios.put(`${API_BASE_URL}/diseases/${id}`, disease);
    return response.data;
  },

  // Delete a disease
  deleteDisease: async (id) => {
    const response = await axios.delete(`${API_BASE_URL}/diseases/${id}`);
    return response.data;
  },

  // Get diseases for an employee
  getEmployeeDiseases: async (employeeID) => {
    const response = await axios.get(`${API_BASE_URL}/employees/${employeeID}/diseases`);
    return response.data;
  },

  // Add disease to employee
  addDiseaseToEmployee: async (employeeID, diseaseID, dateDiagnosed = '') => {
    const response = await axios.post(`${API_BASE_URL}/employees/${employeeID}/diseases`, {
      disease_id: diseaseID,
      date_diagnosed: dateDiagnosed
    });
    return response.data;
  },

  // Remove disease from employee
  removeDiseaseFromEmployee: async (employeeID, diseaseID) => {
    const response = await axios.delete(`${API_BASE_URL}/employees/${employeeID}/diseases/${diseaseID}`);
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
