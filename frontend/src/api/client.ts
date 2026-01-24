interface CustomConfig extends Omit<RequestInit, "body"> {
  body?: any;
}

export async function client<T>(
  endpoint: string,
  { body, ...customConfig }: CustomConfig = {},
): Promise<T> {
  const token = localStorage.getItem("token");
  const BASE_URL = import.meta.env.VITE_API_BASE_URL;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const config: RequestInit = {
    // Default to POST if body exists, otherwise GET,
    // but allows customConfig.method to override it (for DELETE/PUT)
    method: customConfig.method || (body ? "POST" : "GET"),
    ...customConfig,
    headers: {
      ...headers,
      ...customConfig.headers,
    },
  };

  if (body) {
    config.body = JSON.stringify(body);
  }

  const response = await fetch(`${BASE_URL}${endpoint}`, config);

  if (response.status === 401) {
    // Clear stale session data
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    return Promise.reject(new Error("Session expired. Please login again."));
  }

  if (!response.ok) {
    // Handle cases where the backend might not return JSON on error
    try {
      const errorData = await response.json();
      return Promise.reject(new Error(errorData.message || "Request failed"));
    } catch {
      return Promise.reject(new Error("An unexpected error occurred"));
    }
  }

  if (response.status === 204) {
    return {} as T;
  }

  return response.json();
}
