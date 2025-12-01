import axios from 'axios';

// const api = axios.create({
//   baseURL: '/api',
// });

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
});

// Add token to requests if available
const token = localStorage.getItem('token');
if (token) {
  api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
}

export default api;



