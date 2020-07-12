import {createAction, createReducer} from "@reduxjs/toolkit";

const readDir = createAction('SERVER/READ_DIR')

export const dirReducer = createReducer([], {
    [readDir]: (state, action) => action.payload.files
})

export const clickReadDirButton = (dir) => ({
    type: 'APP/CLICK_READ_DIR_BUTTON',
    payload: dir
})
