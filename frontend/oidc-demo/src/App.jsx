import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'
import { useEffect } from 'react'

function App() {
  const [userData, setUserData] = useState(null)
  const [loggedIn, setLoggedIn] = useState(false)

  useEffect(() => {
    fetch('/private/api/me', {credentials: 'include'})
      .then(response => response.json())
      .then(data => {
        setUserData(data)
        setLoggedIn(true)
      })
      .catch(error => {
        console.error('Error fetching /private/api/me:', error)
        setLoggedIn(false)
      })
  }, [])

  var googleIdentityExists = userData &&
          Array.isArray(userData.identities) &&
            userData.identities.some(
              (identity) => identity.identityProviderId === "google"
            )

  return (
    <>
      <h1>OpenID Connect Demo</h1>
      <div className="card">
        {loggedIn && <h2>Logged in! Link another account:</h2>}
        {!loggedIn && <h2>Not logged in! Please log in or create an account via:</h2>}
        {!googleIdentityExists && (
          <button onClick={() => window.location.href = '/login/google'}>
            Google
          </button>
        )}
        {googleIdentityExists && (
          <p>all account types already linked, no new ones left.</p>
        )}
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
        <div style={{ marginTop: '20px' }}>
          {loggedIn && (
          <button onClick={() => window.location.href = '/logout'}>
            Logout
          </button>
          )}
        </div>
      </div>
    </>
  )
}

export default App
