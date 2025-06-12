import React, { useState, useEffect } from 'react';
import './App.css';
import UserService from './services/UserService';
import AuthForms from './components/AuthForms';
import UserProfile from './components/UserProfile';
import ProtectedRoute from './components/ProtectedRoute';
import { User, CreateUserRequest, ApiError } from './types/User';

interface AppState {
  users: User[];
  loading: boolean;
  error: string;
  newUser: CreateUserRequest;
  currentUser: User | null;
  isAuthenticated: boolean;
}

const App: React.FC = () => {
  const [state, setState] = useState<AppState>({
    users: [],
    loading: false,
    error: '',
    newUser: { name: '', email: '', password: '', role: 'user' },
    currentUser: null,
    isAuthenticated: false,
  });

  useEffect(() => {
    // Проверяем авторизацию при загрузке
    const isAuth = UserService.isAuthenticated();
    const user = UserService.getCurrentUser();
    
    setState(prev => ({
      ...prev,
      isAuthenticated: isAuth,
      currentUser: user,
    }));

    if (isAuth) {
      loadUsers();
    }
  }, []);

  const loadUsers = async (): Promise<void> => {
    setState(prev => ({ ...prev, loading: true }));
    try {
      const response = await UserService.listUsers();
      setState(prev => ({
        ...prev,
        users: response.users || [],
        error: '',
      }));
    } catch (err) {
      const apiError = err as ApiError;
      setState(prev => ({
        ...prev,
        error: 'Ошибка загрузки пользователей: ' + apiError.message,
      }));
    } finally {
      setState(prev => ({ ...prev, loading: false }));
    }
  };

  const handleAuthSuccess = (): void => {
    const user = UserService.getCurrentUser();
    setState(prev => ({
      ...prev,
      isAuthenticated: true,
      currentUser: user,
      error: '',
    }));
    loadUsers();
  };

  const handleLogout = (): void => {
    setState(prev => ({
      ...prev,
      isAuthenticated: false,
      currentUser: null,
      users: [],
      error: '',
    }));
  };

  const handleCreateUser = async (e: React.FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault();
    if (!state.newUser.name || !state.newUser.email || !state.newUser.password) {
      setState(prev => ({ ...prev, error: 'Заполните все поля' }));
      return;
    }

    setState(prev => ({ ...prev, loading: true }));
    try {
      await UserService.createUser(
        state.newUser.name, 
        state.newUser.email, 
        state.newUser.password,
        state.newUser.role
      );
      setState(prev => ({
        ...prev,
        newUser: { name: '', email: '', password: '', role: 'user' },
        error: '',
      }));
      await loadUsers();
    } catch (err) {
      const apiError = err as ApiError;
      setState(prev => ({
        ...prev,
        error: 'Ошибка создания пользователя: ' + apiError.message,
      }));
    } finally {
      setState(prev => ({ ...prev, loading: false }));
    }
  };

  const handleGetUser = async (id: number): Promise<void> => {
    setState(prev => ({ ...prev, loading: true }));
    try {
      const response = await UserService.getUser(id);
      alert(`Пользователь: ${response.user?.name} (${response.user?.email})`);
      setState(prev => ({ ...prev, error: '' }));
    } catch (err) {
      const apiError = err as ApiError;
      setState(prev => ({
        ...prev,
        error: 'Ошибка получения пользователя: ' + apiError.message,
      }));
    } finally {
      setState(prev => ({ ...prev, loading: false }));
    }
  };

  const handleInputChange = (field: keyof CreateUserRequest) => 
    (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>): void => {
      setState(prev => ({
        ...prev,
        newUser: { ...prev.newUser, [field]: e.target.value }
      }));
    };

  const formatDate = (timestamp: number): string => {
    return new Date(timestamp * 1000).toLocaleString('ru-RU');
  };

  // Если не авторизован, показываем форму входа
  if (!state.isAuthenticated) {
    return (
      <div className="App">
        <header className="App-header">
          <h1>🚀 K8s Go gRPC React приложение</h1>
          <p>Микросервисная архитектура с авторизацией</p>
        </header>
        <AuthForms onAuthSuccess={handleAuthSuccess} />
      </div>
    );
  }

  // Авторизованное состояние
  return (
    <div className="App">
      <header className="App-header">
        <h1>🚀 K8s Go gRPC React приложение</h1>
        <p>Микросервисная архитектура с мониторингом</p>
        {state.currentUser && (
          <UserProfile user={state.currentUser} onLogout={handleLogout} />
        )}
      </header>

      <main className="main-content">
        {state.error && <div className="error">{state.error}</div>}

        {/* Создание пользователя только для админов */}
        <ProtectedRoute 
          requiredRole="admin"
          fallback={
            <div className="info-message">
              <h3>ℹ️ Информация</h3>
              <p>Создание пользователей доступно только администраторам.</p>
            </div>
          }
        >
          <section className="create-user">
            <h2>Создать пользователя (Админ)</h2>
            <form onSubmit={handleCreateUser}>
              <div className="form-group">
                <input
                  type="text"
                  placeholder="Имя"
                  value={state.newUser.name}
                  onChange={handleInputChange('name')}
                  disabled={state.loading}
                />
              </div>
              <div className="form-group">
                <input
                  type="email"
                  placeholder="Email"
                  value={state.newUser.email}
                  onChange={handleInputChange('email')}
                  disabled={state.loading}
                />
              </div>
              <div className="form-group">
                <input
                  type="password"
                  placeholder="Пароль"
                  value={state.newUser.password}
                  onChange={handleInputChange('password')}
                  disabled={state.loading}
                />
              </div>
              <div className="form-group">
                <select
                  value={state.newUser.role}
                  onChange={handleInputChange('role')}
                  disabled={state.loading}
                >
                  <option value="user">👤 Пользователь</option>
                  <option value="moderator">🛡️ Модератор</option>
                  <option value="admin">👑 Администратор</option>
                </select>
              </div>
              <button type="submit" disabled={state.loading}>
                {state.loading ? 'Создание...' : 'Создать пользователя'}
              </button>
            </form>
          </section>
        </ProtectedRoute>

        <section className="users-list">
          <h2>Список пользователей</h2>
          <button onClick={loadUsers} disabled={state.loading} className="refresh-btn">
            {state.loading ? 'Загрузка...' : 'Обновить список'}
          </button>
          
          {state.users.length === 0 && !state.loading ? (
            <p>Пока нет пользователей</p>
          ) : (
            <div className="users-grid">
              {state.users.map((user: User) => (
                <div key={user.id} className="user-card">
                  <div className="user-card-header">
                    <span className="user-role-badge" data-role={user.role}>
                      {user.role === 'admin' ? '👑' : '👤'} {user.role}
                    </span>
                    <span className={`user-status ${user.is_active ? 'active' : 'inactive'}`}>
                      {user.is_active ? '✅' : '❌'}
                    </span>
                  </div>
                  <h3>{user.name}</h3>
                  <p>{user.email}</p>
                  <p>ID: {user.id}</p>
                  <p>Создан: {formatDate(user.created_at)}</p>
                  <button 
                    onClick={() => handleGetUser(user.id)}
                    disabled={state.loading}
                    className="get-user-btn"
                  >
                    Получить детали
                  </button>
                </div>
              ))}
            </div>
          )}
        </section>

        <section className="monitoring-links">
          <h2>🔗 Ссылки на мониторинг</h2>
          <div className="links-grid">
            <a href="http://localhost:3001" target="_blank" rel="noopener noreferrer" className="monitoring-link">
              📊 Grafana
            </a>
            <a href="http://localhost:9091" target="_blank" rel="noopener noreferrer" className="monitoring-link">
              🔍 Prometheus
            </a>
            <a href="http://localhost:9000" target="_blank" rel="noopener noreferrer" className="monitoring-link">
              📝 Graylog
            </a>
            <a href="http://localhost:8081/health" target="_blank" rel="noopener noreferrer" className="monitoring-link">
              ❤️ Health Check
            </a>
          </div>
        </section>
      </main>
    </div>
  );
};

export default App; 