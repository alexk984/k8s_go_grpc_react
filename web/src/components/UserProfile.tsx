import React, { useState, useEffect } from 'react';
import UserService from '../services/UserService';
import { User } from '../types/User';

interface UserProfileProps {
  user: User;
  onLogout: () => void;
}

const UserProfile: React.FC<UserProfileProps> = ({ user, onLogout }) => {
  const [tokenInfo, setTokenInfo] = useState<{
    loginTime: string | null;
    expirationTime: string | null;
  }>({
    loginTime: null,
    expirationTime: null,
  });

  useEffect(() => {
    setTokenInfo({
      loginTime: UserService.getLoginTime(),
      expirationTime: UserService.getTokenExpiration(),
    });
  }, [user]);

  const handleLogout = (): void => {
    UserService.logout();
    onLogout();
  };

  const getRoleColor = (role: string): string => {
    switch (role) {
      case 'admin': return '#e74c3c';
      case 'moderator': return '#f39c12';
      case 'user': return '#27ae60';
      default: return '#95a5a6';
    }
  };

  const getRoleIcon = (role: string): string => {
    switch (role) {
      case 'admin': return '👑';
      case 'moderator': return '🛡️';
      case 'user': return '👤';
      default: return '❓';
    }
  };

  return (
    <div className="user-profile">
      <div className="profile-card">
        <div className="profile-header">
          <h3>👋 Добро пожаловать, {user.name}!</h3>
          <div className="user-role" style={{ color: getRoleColor(user.role) }}>
            {getRoleIcon(user.role)} {user.role.toUpperCase()}
          </div>
        </div>
        
        <div className="profile-info">
          <div className="info-item">
            <strong>📧 Email:</strong> {user.email}
          </div>
          <div className="info-item">
            <strong>🆔 ID:</strong> {user.id}
          </div>
          {user.created_at && (
            <div className="info-item">
              <strong>📅 Регистрация:</strong> {new Date(user.created_at * 1000).toLocaleString('ru-RU')}
            </div>
          )}
        </div>

        {/* JWT Token Info */}
        <div className="token-info">
          <h4>🔐 Информация о сессии</h4>
          {tokenInfo.loginTime && (
            <div className="info-item">
              <strong>🕐 Время входа:</strong> {tokenInfo.loginTime}
            </div>
          )}
          {tokenInfo.expirationTime && (
            <div className="info-item">
              <strong>⏰ Токен действует до:</strong> {tokenInfo.expirationTime}
            </div>
          )}
          <div className="info-item">
            <strong>🔑 Статус токена:</strong> 
            <span className="token-status active"> Активен</span>
          </div>
        </div>
        
        <button onClick={handleLogout} className="logout-button">
          🚪 Выйти
        </button>
      </div>
    </div>
  );
};

export default UserProfile; 