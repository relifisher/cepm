import axios from 'axios';

const apiClient = axios.create({
  baseURL: 'http://localhost:8090/api/v1', // Now running frontend locally, use full backend URL
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add a request interceptor to include the JWT token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('jwtToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Personal Performance
export const createPerformanceReview = (reviewData) => {
  return apiClient.post('/reviews', reviewData);
};

export const getPerformanceReview = (id) => {
  return apiClient.get(`/reviews/${id}`);
};

export const listUserReviews = (userId) => {
  // In a real app, userId would be handled by the session, not a query param
  return apiClient.get(`/reviews?userId=${userId}`);
};

export const submitPerformanceReview = (id, userId) => {
  return apiClient.post(`/reviews/${id}/submit?userId=${userId}`);
};

// Manager/Team Performance
export const listTeamReviews = (managerId) => {
  return apiClient.get(`/team/reviews?managerId=${managerId}`);
};

export const scorePerformanceReview = (id, scoreData) => {
  return apiClient.post(`/reviews/${id}/score`, scoreData);
};

export const approvePerformanceReview = (id, approverId, comment = '') => {
  return apiClient.post(`/reviews/${id}/approve?approverId=${approverId}`, { comment });
};

export const rejectPerformanceReview = (id, approverId, comment = '') => {
  return apiClient.post(`/reviews/${id}/reject?approverId=${approverId}`, { comment });
};

export const getReviewByPeriod = (userId, period) => {
  return apiClient.get(`/reviews/by-period?userId=${userId}&period=${period}`);
};

export const updatePerformanceReview = (id, reviewData) => {
  return apiClient.put(`/reviews/${id}`, reviewData);
};

export const getAllSubmittedReviews = (userId) => {
  return apiClient.get(`/reviews/all-submitted?userId=${userId}`);
};

export const getAllReviewsByPeriod = (period) => {
  return apiClient.get(`/reviews/all-by-period?period=${period}`);
};

export default apiClient;
