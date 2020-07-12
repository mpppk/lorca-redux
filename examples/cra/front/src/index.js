import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import {redux} from './redux';
import { Provider } from 'react-redux';

ReactDOM.render(
  <React.StrictMode>
    <Provider store={redux}>
      <App />
    </Provider>
  </React.StrictMode>,
  document.getElementById('root')
);

