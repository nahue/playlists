// Authentication utility functions

export const auth = {
  // Get the current token from localStorage
  getToken() {
    return localStorage.getItem('token');
  },

  // Get the current user from localStorage
  getUser() {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  },

  // Check if user is authenticated
  isAuthenticated() {
    return !!this.getToken();
  },

  // Set authentication data
  setAuth(token, user) {
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
  },

  // Clear authentication data
  clearAuth() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  // Make authenticated API request
  async apiRequest(url, options = {}) {
    const token = this.getToken();
    if (!token) {
      throw new Error('No authentication token');
    }

    const defaultOptions = {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
        ...options.headers
      }
    };

    const response = await fetch(url, { ...defaultOptions, ...options });
    
    if (response.status === 401) {
      this.clearAuth();
      window.location.href = '/login';
      throw new Error('Authentication expired');
    }

    return response;
  },

  // Logout and redirect
  logout() {
    this.clearAuth();
    window.location.href = '/login';
  }
};

// Redirect to login if not authenticated
export function requireAuth() {
  if (!auth.isAuthenticated()) {
    window.location.href = '/login';
    return false;
  }
  return true;
}

// Check if user should be redirected away from login page
export function redirectIfAuthenticated() {
  if (auth.isAuthenticated()) {
    window.location.href = '/';
    return true;
  }
  return false;
} 