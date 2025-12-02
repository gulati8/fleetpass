import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { AuthProvider, useAuth } from './AuthContext';
import * as api from '../services/api';

// Mock the API
jest.mock('../services/api');

// Test component that uses auth context
const TestComponent = () => {
  const { user, isAuthenticated, login, logout } = useAuth();

  return (
    <div>
      <div data-testid="auth-status">{isAuthenticated ? 'authenticated' : 'not authenticated'}</div>
      <div data-testid="user-email">{user?.email || 'no user'}</div>
      <button onClick={() => login('test@example.com', 'password')}>Login</button>
      <button onClick={logout}>Logout</button>
    </div>
  );
};

describe('AuthContext', () => {
  beforeEach(() => {
    localStorage.clear();
    jest.clearAllMocks();
  });

  test('initially not authenticated when no stored token', () => {
    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    );

    expect(screen.getByTestId('auth-status')).toHaveTextContent('not authenticated');
    expect(screen.getByTestId('user-email')).toHaveTextContent('no user');
  });

  test('loads user from localStorage on mount', async () => {
    const mockUser = { email: 'test@example.com', id: '123' };
    localStorage.setItem('token', 'mock-token');
    localStorage.setItem('user', JSON.stringify(mockUser));

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    );

    await waitFor(() => {
      expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated');
      expect(screen.getByTestId('user-email')).toHaveTextContent('test@example.com');
    });
  });

  test('login sets user and token', async () => {
    const mockResponse = {
      token: 'new-token',
      user: { email: 'test@example.com', id: '123' }
    };

    api.authAPI.login = jest.fn().mockResolvedValue(mockResponse);

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    );

    screen.getByText('Login').click();

    await waitFor(() => {
      expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated');
      expect(screen.getByTestId('user-email')).toHaveTextContent('test@example.com');
      expect(localStorage.getItem('token')).toBe('new-token');
    });
  });

  test('logout clears user and token', async () => {
    const mockUser = { email: 'test@example.com', id: '123' };
    localStorage.setItem('token', 'mock-token');
    localStorage.setItem('user', JSON.stringify(mockUser));

    render(
      <AuthProvider>
        <TestComponent />
      </AuthProvider>
    );

    await waitFor(() => {
      expect(screen.getByTestId('auth-status')).toHaveTextContent('authenticated');
    });

    screen.getByText('Logout').click();

    await waitFor(() => {
      expect(screen.getByTestId('auth-status')).toHaveTextContent('not authenticated');
      expect(screen.getByTestId('user-email')).toHaveTextContent('no user');
      expect(localStorage.getItem('token')).toBeNull();
      expect(localStorage.getItem('user')).toBeNull();
    });
  });
});
