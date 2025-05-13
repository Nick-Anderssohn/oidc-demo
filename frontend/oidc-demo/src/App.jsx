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

  const slimmedIdentities = userData && Array.isArray(userData.identities)
    ? userData.identities.map(({ identityProviderId, externalId, mostRecentIdToken }) => ({
        identityProviderId,
        externalId,
        email: mostRecentIdToken.email
      }))
    : [];

  return (
    <>
      <h1>OpenID Connect Demo</h1>
      <div className="card">
        {loggedIn && <h2>Logged in! Link another account:</h2>}
        {!loggedIn && <h2>Not logged in! Please log in or create an account via:</h2>}
        <button onClick={() => window.location.href = '/login/google'}>
            Google
        </button>
        <hr style={{ margin: '5px 0' }} />
        {userData && (
          <div className="user-info">
            <h2>User Info</h2>
            <div style={{ textAlign: 'left' }}>
              <p><strong>User Email:</strong> {userData.email}</p>
            </div>
            <h3>Linked Accounts</h3>
            <ul>
              {slimmedIdentities.map((identity, idx) => (
                <li key={idx}>
                  <div style={{ textAlign: 'left' }}>
                    <strong>Provider:</strong> {identity.identityProviderId} <br />
                    <strong>External ID:</strong> {identity.externalId} <br />
                    <strong>Email:</strong> {identity.email}
                  </div>
                </li>
              ))}
            </ul>
          </div>
        )}
        <div style={{ marginTop: '20px' }}>
          {loggedIn && (
            <>
          <button onClick={() => window.location.href = '/logout'}>
            Logout
          </button>
          <button
            onClick={async () => {
              try {
                const response = await fetch('/private/api/me', {
                  method: 'DELETE',
                  credentials: 'include',
                });
                if (response.ok) {
                  setUserData(null);
                  setLoggedIn(false);
                } else {
                  console.error('Failed to delete user');
                }
              } catch (error) {
                console.error('Error deleting user:', error);
              }
            }}
            style={{ marginLeft: '10px', backgroundColor: '#e74c3c', color: 'white' }}
          >
            Hard Delete My Account
          </button>
          </>
          )}
        </div>
      </div>
    </>
  )
}

export default App
