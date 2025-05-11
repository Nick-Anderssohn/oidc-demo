import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import { useEffect } from 'react'

function App() {
  const [userData, setUserData] = useState(null)
  const [count, setCount] = useState(0)

  useEffect(() => {
    fetch('/private/api/me', {credentials: 'include'})
      .then(response => response.json())
      .then(data => {
        setUserData(data)
      })
      .catch(error => {
        console.error('Error fetching /private/api/me:', error)
      })
  }, [])

  return (
    <>
      <div>
        <a href="https://vite.dev" target="_blank">
          <img src={viteLogo} className="logo" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        {userData && (
          <div className="user-info">
            <h2>User Info</h2>
            <textarea
              readOnly
              style={{ width: '975px', height: '388px', fontFamily: 'monospace', fontSize: '1rem' }}
              value={JSON.stringify(userData, null, 2)}
            />
          </div>
        )}
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/App.jsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </>
  )
}

export default App
