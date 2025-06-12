import React from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';

// Mock UserService to avoid API calls in tests
jest.mock('./services/UserService', () => ({
  listUsers: jest.fn().mockResolvedValue({ users: [] }),
  createUser: jest.fn().mockResolvedValue({}),
  getUser: jest.fn().mockResolvedValue({ user: null }),
}));

describe('App Component', () => {
  test('renders main heading', () => {
    render(<App />);
    const headingElement = screen.getByText(/K8s Go gRPC React приложение/i);
    expect(headingElement).toBeInTheDocument();
  });

  test('renders create user section', () => {
    render(<App />);
    const createUserHeading = screen.getByText(/Создать пользователя/i);
    expect(createUserHeading).toBeInTheDocument();
  });

  test('renders users list section', () => {
    render(<App />);
    const usersListHeading = screen.getByText(/Список пользователей/i);
    expect(usersListHeading).toBeInTheDocument();
  });

  test('renders monitoring links section', () => {
    render(<App />);
    const monitoringHeading = screen.getByText(/Ссылки на мониторинг/i);
    expect(monitoringHeading).toBeInTheDocument();
  });

  test('has form inputs', () => {
    render(<App />);
    const nameInput = screen.getByPlaceholderText(/Имя/i);
    const emailInput = screen.getByPlaceholderText(/Email/i);
    
    expect(nameInput).toBeInTheDocument();
    expect(emailInput).toBeInTheDocument();
  });
}); 