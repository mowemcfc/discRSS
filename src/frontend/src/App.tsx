import React from 'react';
import logo from './logo.svg';
import './App.css';

import { LoginPage } from './pages/login-page';
import { ProfilePage } from './pages/profile-page';

import {
  Routes,
  Route,
  BrowserRouter
} from 'react-router-dom';
import { ProtectedRoute } from './components/protected';

export const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route 
          path="/" 
          element={<ProtectedRoute component={ProfilePage} />} 
        />

        <Route 
          path="/login" 
          element={<LoginPage />} 
        />
      </Routes>
    </BrowserRouter>
  );
}

