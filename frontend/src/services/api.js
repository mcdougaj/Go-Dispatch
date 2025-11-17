import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: `${API_URL}/api`,
  headers: {
    'Content-Type': 'application/json',
  },
});

export const fetchVehicles = async () => {
  const response = await api.get('/vehicles');
  return response.data || [];
};

export const fetchDrivers = async () => {
  const response = await api.get('/drivers');
  return response.data || [];
};

export const fetchAssets = async () => {
  const response = await api.get('/assets');
  return response.data || [];
};

export const fetchGeofences = async () => {
  const response = await api.get('/geofences');
  return response.data || [];
};

export const sendQuery = async (query) => {
  const response = await api.post('/query', { query });
  return response.data;
};

export const geocode = async (address) => {
  const response = await api.get('/geocode', {
    params: { address },
  });
  return response.data;
};

export const reverseGeocode = async (lat, lon) => {
  const response = await api.get('/reverse-geocode', {
    params: { lat, lon },
  });
  return response.data;
};

export const calculateDistance = async (originLat, originLon, destLat, destLon) => {
  const response = await api.get('/distance', {
    params: {
      origin_lat: originLat,
      origin_lon: originLon,
      dest_lat: destLat,
      dest_lon: destLon,
    },
  });
  return response.data;
};

export default api;
