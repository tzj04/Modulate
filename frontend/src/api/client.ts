interface CustomConfig extends Omit<RequestInit, "body"> {
  body?: any;
}

export async function client<T>(
  endpoint: string,
  { body, ...customConfig }: CustomConfig = {},
): Promise<T> {
  const BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";
  let token = localStorage.getItem("token");

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }

  const config: RequestInit = {
    method: customConfig.method || (body ? "POST" : "GET"),
    ...customConfig,
    headers: {
      ...headers,
      ...customConfig.headers,
    },
    credentials: "include",
  };

  if (body) {
    config.body = JSON.stringify(body);
  }

  let response = await fetch(`${BASE_URL}${endpoint}`, config);

  if (
    response.status === 401 &&
    !endpoint.includes("/auth/refresh") &&
    !endpoint.includes("/auth/login")
  ) {
    try {
      // Attempt to get a new access token
      const refreshResponse = await fetch(`${BASE_URL}/auth/refresh`, {
        method: "POST",
        credentials: "include", // Send the HttpOnly cookie
      });

      if (refreshResponse.ok) {
        const { access_token } = await refreshResponse.json();

        // Save new token
        localStorage.setItem("token", access_token);

        // Retry the original request with the new token
        const newHeaders = {
          ...config.headers,
          Authorization: `Bearer ${access_token}`,
        };

        response = await fetch(`${BASE_URL}${endpoint}`, {
          ...config,
          headers: newHeaders,
        });
      } else {
        throw new Error("Refresh failed");
      }
    } catch (err) {
      // If refresh fails, clear session and redirect
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      window.location.href = "/login";
      return Promise.reject(new Error("Session expired. Please login again."));
    }
  }

  if (!response.ok) {
    try {
      const errorData = await response.json();
      return Promise.reject(new Error(errorData.message || "Request failed"));
    } catch {
      return Promise.reject(new Error("An unexpected error occurred"));
    }
  }

  if (response.status === 204) return {} as T;
  return response.json();
}
