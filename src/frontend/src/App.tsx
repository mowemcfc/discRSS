import React from 'react';
import logo from './logo.svg';
import './App.css';

import LoginButton from './components/login';
import LogoutButton from './components/logout';
import Profile from './components/profile';

function App() {
  return (
    <div className='App'>
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
        <LoginButton></LoginButton>
        <LogoutButton></LogoutButton>
        <Profile></Profile>
    </div>
  );
}

export default App;
