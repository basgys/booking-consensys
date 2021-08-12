import React from 'react';
import ReactDOM from 'react-dom';
import {
  BrowserRouter as Router, Link, Route, Switch
} from "react-router-dom";
import Room from './booking/room';
import Rooms from './booking/rooms';
import './index.css';
import reportWebVitals from './reportWebVitals';
import Auth, { AuthForm } from './user/auth';

ReactDOM.render(
  <React.StrictMode>
    <Router>
      <Auth>
        <div className="app">
          <nav className="app__nav">
            <ul className="app__nav-items">
              <li className="app__nav-item">
                <Link to="/rooms">Rooms</Link>
              </li>
            </ul>
            <ul className="app__nav-items">
              <li className="app__nav-item">
                <Link to="/auth">Login</Link>
              </li>
            </ul>
          </nav>
          <main role="main" className="app__main">
            <Switch>
              <Route path="/auth">
                <AuthForm />
              </Route>
              <Route path="/rooms/:id">
                <Room />
              </Route>
              <Route path="/">
                <Rooms />
              </Route>
            </Switch>
          </main>
        </div>
      </Auth>
    </Router>
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
