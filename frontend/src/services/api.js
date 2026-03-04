import axios from 'axios';

const API_BASE_URL = '/api';

export const authService = {
  login: async (payload) => {
    const response = await axios.post(`${API_BASE_URL}/auth/login`, payload);
    return response.data;
  },

  signup: async (payload) => {
    const response = await axios.post(`${API_BASE_URL}/auth/signup`, payload);
    return response.data;
  }
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

// Abdomen API service
export const abdomenService = {
  // Get all abdomen options
  getAllAbdomens: async () => {
    const response = await axios.get(`${API_BASE_URL}/abdomens`);
    return response.data;
  }
};
