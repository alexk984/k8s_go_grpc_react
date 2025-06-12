export interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  is_active: boolean;
  created_at: number;
}

export interface CreateUserRequest {
  name: string;
  email: string;
  password?: string;
  role?: string;
}

// Авторизация
export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: User;
  message: string;
}

export interface GetUserResponse {
  user: User;
  message: string;
}

export interface CreateUserResponse {
  user: User;
  message: string;
}

export interface ListUsersResponse {
  users: User[];
  total: number;
}

export interface ApiError {
  message: string;
  status?: number;
} 