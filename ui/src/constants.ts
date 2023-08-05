export const API_BASE_URL = process.env.NODE_ENV === "prod" ? "https://api.mangovb.com" : "http://localhost:8080";
export const API_FILTERS_URL = `${API_BASE_URL}/filters`;
export const API_EVENTS_URL = `${API_BASE_URL}/events`;
