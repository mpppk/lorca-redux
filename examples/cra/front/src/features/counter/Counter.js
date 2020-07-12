/*global dispatchToServer*/

import React, { useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import {
  decrement,
  increment,
  incrementByAmount,
  incrementAsync,
  selectCount, clickReadDirButton,
} from './counterSlice';
import styles from './Counter.module.css';


export function Counter() {
  const count = useSelector(selectCount);
  const files = useSelector((s) => s.dir);
  const dispatch = useDispatch();
  const [incrementAmount, setIncrementAmount] = useState('2');

  const handleClickIncrementButton = () => {
    console.log('increment');
    if (dispatchToServer !== undefined) {
      dispatchToServer({A: 1});
    }
    dispatch(increment());
  }

  return (
    <div>
      <div className={styles.row}>
        <button
          className={styles.button}
          aria-label="Increment value"
          onClick={handleClickIncrementButton}
        >
          +
        </button>
        <span className={styles.value}>{count}</span>
        <button
          className={styles.button}
          aria-label="Decrement value"
          onClick={() => dispatch(decrement())}
        >
          -
        </button>
      </div>
      <div className={styles.row}>
        <input
          className={styles.textbox}
          aria-label="Set increment amount"
          value={incrementAmount}
          onChange={e => setIncrementAmount(e.target.value)}
        />
        <button
          className={styles.button}
          onClick={() =>
            dispatch(incrementByAmount(Number(incrementAmount) || 0))
          }
        >
          Add Amount
        </button>
        <button
          className={styles.asyncButton}
          onClick={() => dispatch(incrementAsync(Number(incrementAmount) || 0))}
        >
          Add Async
        </button>
        <button
            className={styles.button}
            onClick={() =>
                dispatch(clickReadDirButton('.'))
            }
        >
          Read Current Dir
        </button>
      </div>
      {files.join(' ')}
    </div>
  );
}
