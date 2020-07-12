/*global dispatchToServer*/

import { configureStore } from '@reduxjs/toolkit';
import counterReducer, {dirReducer} from '../features/counter/counterSlice';

const reducer = {
  counter: counterReducer,
  dir: dirReducer
}

configureStore({reducer});

const makeLorcaMiddleware = (dispatchToServer) => (_store) => (next) => (action) => {
  dispatchToServer(action)
  next(action)
}

const makeServerActionHandler = (store) => {
  return (action) => {
    store.dispatch(action);
  }
}

const makeStore = ({reducer}) => {
  const store =  configureStore({reducer, middleware: [makeLorcaMiddleware(dispatchToServer)]});
  window.handleServerAction = makeServerActionHandler(store);
  return store;
}

export const store = makeStore({reducer})
