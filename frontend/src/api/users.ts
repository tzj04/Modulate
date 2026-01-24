// src/api/users.ts
import { client } from "./client";

export interface User {
  id: number;
  username: string;
  label: string | null;
  is_deleted: boolean;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface LoginCredentials {
  username: string;
  password: string;
}

export interface RegisterData {
  username: string;
  password: string;
  label?: string;
}

export const userApi = {
  login: (credentials: LoginCredentials) =>
    client<AuthResponse>("/auth/login", {
      method: "POST",
      body: credentials,
    }),

  register: (data: RegisterData) =>
    client<{ message: string }>("/auth/register", {
      method: "POST",
      body: data,
    }),
};
