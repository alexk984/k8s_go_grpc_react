import React, { useState } from 'react';
import UserService from '../services/UserService';
import { RegisterRequest, LoginRequest, ApiError } from '../types/User';

interface AuthFormsProps {
  onAuthSuccess: () => void;
}

interface AuthState {
  loading: boolean;
  error: string;
  isLoginMode: boolean;
}

interface FormData {
  name: string;
  email: string;
  password: string;
}

const AuthForms: React.FC<AuthFormsProps> = ({ onAuthSuccess }) => {
  const [authState, setAuthState] = useState<AuthState>({
    loading: false,
    error: '',
    isLoginMode: true,
  });
  
  const [formData, setFormData] = useState<FormData>({
    name: '',
    email: '',
    password: '',
  });

  const handleInputChange = (field: keyof FormData) => 
    (e: React.ChangeEvent<HTMLInputElement>): void => {
      setFormData({ ...formData, [field]: e.target.value });
    };

  const toggleMode = (): void => {
    setAuthState(prev => ({ ...prev, isLoginMode: !prev.isLoginMode, error: '' }));
    setFormData({ name: '', email: '', password: '' });
  };

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault();
    
    if (!formData.email || !formData.password) {
      setAuthState(prev => ({ ...prev, error: 'Заполните все обязательные поля' }));
      return;
    }

    if (!authState.isLoginMode && !formData.name) {
      setAuthState(prev => ({ ...prev, error: 'Имя обязательно для регистрации' }));
      return;
    }

    setAuthState(prev => ({ ...prev, loading: true, error: '' }));

    try {
      if (authState.isLoginMode) {
        const loginRequest: LoginRequest = {
          email: formData.email,
          password: formData.password,
        };
        await UserService.login(loginRequest);
      } else {
        const registerRequest: RegisterRequest = {
          name: formData.name,
          email: formData.email,
          password: formData.password,
        };
        await UserService.register(registerRequest);
      }
      
      onAuthSuccess();
    } catch (err) {
      const apiError = err as ApiError;
      setAuthState(prev => ({ 
        ...prev, 
        error: `Ошибка ${authState.isLoginMode ? 'входа' : 'регистрации'}: ${apiError.message}` 
      }));
    } finally {
      setAuthState(prev => ({ ...prev, loading: false }));
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h2>{authState.isLoginMode ? '🔐 Вход в систему' : '📝 Регистрация'}</h2>
        
        {authState.error && <div className="error">{authState.error}</div>}
        
        <form onSubmit={handleSubmit} className="auth-form">
          {!authState.isLoginMode && (
            <div className="form-group">
              <input
                type="text"
                placeholder="Имя"
                value={formData.name}
                onChange={handleInputChange('name')}
                disabled={authState.loading}
                required={!authState.isLoginMode}
              />
            </div>
          )}
          
          <div className="form-group">
            <input
              type="email"
              placeholder="Email"
              value={formData.email}
              onChange={handleInputChange('email')}
              disabled={authState.loading}
              required
            />
          </div>
          
          <div className="form-group">
            <input
              type="password"
              placeholder="Пароль"
              value={formData.password}
              onChange={handleInputChange('password')}
              disabled={authState.loading}
              required
              minLength={6}
            />
          </div>
          
          <button type="submit" disabled={authState.loading} className="auth-button">
            {authState.loading 
              ? (authState.isLoginMode ? 'Вход...' : 'Регистрация...') 
              : (authState.isLoginMode ? 'Войти' : 'Зарегистрироваться')
            }
          </button>
        </form>
        
        <div className="auth-toggle">
          <p>
            {authState.isLoginMode ? 'Нет аккаунта? ' : 'Уже есть аккаунт? '}
            <button 
              type="button" 
              onClick={toggleMode} 
              disabled={authState.loading}
              className="toggle-button"
            >
              {authState.isLoginMode ? 'Зарегистрироваться' : 'Войти'}
            </button>
          </p>
        </div>
      </div>
    </div>
  );
};

export default AuthForms; 