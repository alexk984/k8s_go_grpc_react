import React from 'react';
import UserService from '../services/UserService';
import { User } from '../types/User';

interface ProtectedRouteProps {
  children: React.ReactNode;
  requiredRole?: string;
  fallback?: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ 
  children, 
  requiredRole, 
  fallback 
}) => {
  const isAuthenticated = UserService.isAuthenticated();
  const currentUser: User | null = UserService.getCurrentUser();

  // Если не авторизован
  if (!isAuthenticated || !currentUser) {
    return (
      <div className="access-denied">
        <h3>🔒 Доступ запрещен</h3>
        <p>Необходима авторизация для доступа к этому разделу.</p>
        {fallback}
      </div>
    );
  }

  // Если требуется определенная роль
  if (requiredRole && currentUser.role !== requiredRole) {
    return (
      <div className="access-denied">
        <h3>🚫 Недостаточно прав</h3>
        <p>
          Требуется роль: <strong>{requiredRole}</strong><br/>
          Ваша роль: <strong>{currentUser.role}</strong>
        </p>
        {fallback}
      </div>
    );
  }

  // Проверяем, не истек ли токен перед рендерингом
  return <>{children}</>;
};

export default ProtectedRoute; 