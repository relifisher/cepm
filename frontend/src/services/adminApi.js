import axios from 'axios';

const API_BASE_URL = 'http://localhost:8090/api/v1';

// Helper to set Authorization header with JWT token
const getAuthHeaders = () => {
  const token = localStorage.getItem('jwtToken');
  return token ? { 'Authorization': `Bearer ${token}` } : {};
};

// User Management
export const getUsers = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/admin/users`, { headers: getAuthHeaders() });
    return response.data;
  } catch (error) {
    console.error('Error fetching users:', error);
    throw error;
  }
};

export const updateUser = async (userId, userData) => {
  try {
    const response = await axios.put(`${API_BASE_URL}/admin/users/${userId}`, userData, { headers: getAuthHeaders() });
    return response.data;
  } catch (error) {
    console.error(`Error updating user ${userId}:`, error);
    throw error;
  }
};

// Department Management
export const createDepartment = async (departmentData) => {
  try {
    const response = await axios.post(`${API_BASE_URL}/admin/departments`, departmentData, { headers: getAuthHeaders() });
    return response.data;
  } catch (error) {
    console.error('Error creating department:', error);
    throw error;
  }
};

export const getDepartments = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/admin/departments`, { headers: getAuthHeaders() });
    return response.data;
  } catch (error) {
    console.error('Error fetching departments:', error);
    throw error;
  }
};

export const getRoles = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/admin/roles`, { headers: getAuthHeaders() });
    return response.data;
  } catch (error) {
    console.error('Error fetching roles:', error);
    throw error;
  }
};

// System Settings
export const updateSystemSetting = async (settingData) => {
  try {
    const response = await axios.put(`${API_BASE_URL}/admin/settings`, settingData, { headers: getAuthHeaders() });
    return response.data;
  } catch (error) {
    console.error('Error updating system setting:', error);
    throw error;
  }
};
