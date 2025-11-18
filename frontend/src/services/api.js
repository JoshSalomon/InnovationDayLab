// API client for making HTTP requests with session cookie handling

const API_BASE = '/api';

async function apiRequest(endpoint, options = {}) {
    const url = `${API_BASE}${endpoint}`;
    const config = {
        ...options,
        credentials: 'include', // Include cookies in requests
        headers: {
            'Content-Type': 'application/json',
            ...options.headers,
        },
    };

    if (options.body && typeof options.body === 'object') {
        config.body = JSON.stringify(options.body);
    }

    const response = await fetch(url, config);

    if (!response.ok) {
        const error = await response.json().catch(() => ({ error: 'Unknown error' }));
        throw new Error(error.error || `HTTP ${response.status}`);
    }

    return response.json();
}

export const api = {
    get: (endpoint) => apiRequest(endpoint, { method: 'GET' }),
    post: (endpoint, data) => apiRequest(endpoint, { method: 'POST', body: data }),
    put: (endpoint, data) => apiRequest(endpoint, { method: 'PUT', body: data }),
    delete: (endpoint) => apiRequest(endpoint, { method: 'DELETE' }),
};
