import { configureStore } from '@reduxjs/toolkit';
import {createAction, createReducer} from "@reduxjs/toolkit";
import {makeLorcaMiddleware, setupServerActionHandler} from "./lib";

const readDir = createAction('SERVER/READ_DIR')

export const clickReadDirButton = (dir) => ({
  type: 'APP/CLICK_READ_DIR_BUTTON',
  payload: dir
})

const dir = createReducer([], {
  [readDir]: (state, action) => action.payload.files
})
const reducer = {dir}

configureStore({reducer});

const makeStore = ({reducer}) => {
  const store =  configureStore({reducer, middleware: [makeLorcaMiddleware()]});
  setupServerActionHandler(store);
  return store;
}

export const redux = makeStore({reducer})
