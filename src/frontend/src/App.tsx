import React from 'react';
import logo from './logo.svg';
import './App.css';

import { LoginPage } from './components/login';
import { LogoutButton } from './components/logout';
import UserProfile from './components/profile';
import { Auth0ContextInterface, useAuth0, withAuth0 } from "@auth0/auth0-react"

import {
  Routes,
  Route,
  redirect,
  createBrowserRouter,
  RouterProvider,
  LoaderFunction,
  BrowserRouter
} from 'react-router-dom';
import { ProtectedRoute } from './components/protected';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route 
          path="/" 
          element={<ProtectedRoute component={UserProfile} />} 
        />

        <Route 
          path="/login" 
          element={<LoginPage />} 
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
