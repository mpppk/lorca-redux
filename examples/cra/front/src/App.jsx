import React from 'react';
import logo from './logo.svg';
import './App.css';
import styles from "./global.module.css";
import {clickReadDirButton} from "./redux";
import {useDispatch, useSelector} from "react-redux";

const toFilesText = (files, maxNum) => {
  const fileMaxNum = files.length >=maxNum ? maxNum: files.length;
  let fileNames = files.slice(0, fileMaxNum).map(f => `"${f}"`).join(', ');
  fileNames += files.length > maxNum ? `, and ${files.length-maxNum} files` : ''
  return fileNames;
}

const formatDir = (s) => {
  return toFilesText(s.dir, 2)
}

function App() {
  const dispatch = useDispatch();
  const fileNames = useSelector(formatDir);

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <button
            className={styles.button}
            onClick={() =>
                dispatch(clickReadDirButton('.'))
            }
        >
          Read Current Dir
        </button>
        {fileNames}
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <span>
          <span>Learn </span>
          <a
            className="App-link"
            href="https://github.com/mpppk/lorca-fsa"
            target="_blank"
            rel="noopener noreferrer"
          >
            lorca-fsa
          </a>
        </span>
      </header>
    </div>
  );
}

export default App;
