import {
  CreateUserRequest, 
  GetUserResponse, 
  CreateUserResponse, 
  ListUsersResponse,
  RegisterRequest,
  LoginRequest,
  AuthResponse,
  User,
  ApiError 
} from '../types/User';

// JWT утилиты
interface JWTPayload {
  exp: number;
  user_id: number;
  email: string;
}

// Поскольку gRPC не работает напрямую в браузере, используем HTTP gateway
// В реальном проекте можно использовать grpc-web или создать REST API прокси

const API_BASE: string = process.env.REACT_APP_API_BASE || 'http://localhost:8081/api/v1';

class UserService {
  private getAuthHeaders(): HeadersInit {
    const token = localStorage.getItem('authToken');
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };
    
    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }
    
    return headers;
  }

  // Проверяем, не истек ли токен
  private isTokenExpired(): boolean {
    const token = localStorage.getItem('authToken');
    if (!token) return true;

    try {
      const payload = JSON.parse(atob(token.split('.')[1])) as JWTPayload;
      const currentTime = Math.floor(Date.now() / 1000);
      return payload.exp < currentTime;
    } catch {
      return true;
    }
  }

  // Автоматический logout при истечении токена
  private checkTokenAndLogout(): void {
    if (this.isTokenExpired()) {
      this.logout();
      window.location.reload(); // Перезагружаем страницу для показа формы входа
    }
  }

  private async handleResponse<T>(response: Response): Promise<T> {
    // Автоматический logout при 401 ошибке
    if (response.status === 401) {
      this.logout();
      window.location.reload();
      throw new Error('Сессия истекла. Войдите заново.');
    }

    if (!response.ok) {
      let errorMessage = `HTTP error! status: ${response.status}`;
      
      try {
        const errorData = await response.json();
        errorMessage = errorData.message || errorMessage;
      } catch {
        // Если не удалось получить JSON, используем базовое сообщение
      }
      
      const apiError: ApiError = {
        message: errorMessage,
        status: response.status
      };
      throw apiError;
    }
    return response.json();
  }

  // Авторизация
  async register(request: RegisterRequest): Promise<AuthResponse> {
    const response = await fetch(`${API_BASE}/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });
    
    const result = await this.handleResponse<AuthResponse>(response);
    
    // Сохраняем токен и пользователя
    if (result.token) {
      localStorage.setItem('authToken', result.token);
      localStorage.setItem('currentUser', JSON.stringify(result.user));
      localStorage.setItem('loginTime', Date.now().toString());
    }
    
    return result;
  }

  async login(request: LoginRequest): Promise<AuthResponse> {
    const response = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });
    
    const result = await this.handleResponse<AuthResponse>(response);
    
    // Сохраняем токен и пользователя
    if (result.token) {
      localStorage.setItem('authToken', result.token);
      localStorage.setItem('currentUser', JSON.stringify(result.user));
      localStorage.setItem('loginTime', Date.now().toString());
    }
    
    return result;
  }

  logout(): void {
    localStorage.removeItem('authToken');
    localStorage.removeItem('currentUser');
    localStorage.removeItem('loginTime');
  }

  getCurrentUser(): User | null {
    // Проверяем токен перед возвратом пользователя
    if (this.isTokenExpired()) {
      this.logout();
      return null;
    }

    const userStr = localStorage.getItem('currentUser');
    if (!userStr) return null;
    
    try {
      return JSON.parse(userStr);
    } catch {
      return null;
    }
  }

  isAuthenticated(): boolean {
    const hasToken = !!localStorage.getItem('authToken');
    if (!hasToken) return false;
    
    // Проверяем, не истек ли токен
    if (this.isTokenExpired()) {
      this.logout();
      return false;
    }
    
    return true;
  }

  // Получаем время входа для отображения
  getLoginTime(): string | null {
    const loginTime = localStorage.getItem('loginTime');
    if (!loginTime) return null;
    
    return new Date(parseInt(loginTime)).toLocaleString('ru-RU');
  }

  // Получаем время истечения токена
  getTokenExpiration(): string | null {
    const token = localStorage.getItem('authToken');
    if (!token) return null;

    try {
      const payload = JSON.parse(atob(token.split('.')[1])) as JWTPayload;
      return new Date(payload.exp * 1000).toLocaleString('ru-RU');
    } catch {
      return null;
    }
  }

  // Пользователи (требуют авторизации)
  async getUser(id: number): Promise<GetUserResponse> {
    this.checkTokenAndLogout();
    
    const response = await fetch(`${API_BASE}/users/${id}`, {
      headers: this.getAuthHeaders(),
    });
    return this.handleResponse<GetUserResponse>(response);
  }

  async createUser(name: string, email: string, password?: string, role?: string): Promise<CreateUserResponse> {
    this.checkTokenAndLogout();
    
    const requestData: CreateUserRequest = { 
      name, 
      email, 
      password: password || 'defaultpassword123',
      role: role || 'user'
    };
    
    const response = await fetch(`${API_BASE}/users`, {
      method: 'POST',
      headers: this.getAuthHeaders(),
      body: JSON.stringify(requestData),
    });
    
    return this.handleResponse<CreateUserResponse>(response);
  }

  async listUsers(): Promise<ListUsersResponse> {
    this.checkTokenAndLogout();
    
    const response = await fetch(`${API_BASE}/users`, {
      headers: this.getAuthHeaders(),
    });
    return this.handleResponse<ListUsersResponse>(response);
  }

  async healthCheck(): Promise<string> {
    const response = await fetch(`${API_BASE.replace('/v1', '')}/health`);
    if (!response.ok) {
      throw new Error(`Health check failed: ${response.status}`);
    }
    return response.text();
  }
}

export default new UserService(); 