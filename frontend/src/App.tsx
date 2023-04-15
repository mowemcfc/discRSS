import React from 'react';
import './App.css';

import { UserPage } from './pages/user-page';

import {
  Routes,
  Route,
  BrowserRouter
} from 'react-router-dom';
import { ProtectedRoute } from './components/protected';
import { HomePage } from './pages/home-page';

export const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/"
          element={<HomePage />}
        />

        <Route
          path="/account"
          element={<ProtectedRoute component={UserPage} />}
        />
      </Routes>
    </BrowserRouter>
  );
}
