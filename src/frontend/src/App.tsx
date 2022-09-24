import React from 'react';
import logo from './logo.svg';
import './App.css';

import LoginButton from './components/login';
import LogoutButton from './components/logout';
import Profile from './components/profile';
import { useAuth0 } from "@auth0/auth0-react"

function App() {
  const { isAuthenticated } = useAuth0();
  return (
    isAuthenticated ?
    <div className='App'>
        <LogoutButton></LogoutButton>
        <Profile></Profile>
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <p>
            Edit <code>src/App.tsx</code> and save to reload.
          </p>
          <a
            className="App-link"
            href="https://reactjs.org"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn React
          </a>
        </header>
    </div>
    :
    <LoginButton></LoginButton>
  );
}

export default App;
