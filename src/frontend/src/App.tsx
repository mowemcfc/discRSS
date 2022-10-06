import React from 'react';
import './App.css';

import { LoginPage } from './pages/login-page';
import { UserPage } from './pages/user-page';

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
          element={<ProtectedRoute component={UserPage} />} 
        />

        <Route 
          path="/login" 
          element={<LoginPage />}
        />
      </Routes>
    </BrowserRouter>
  );
}

