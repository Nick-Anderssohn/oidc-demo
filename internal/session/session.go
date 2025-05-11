package session

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Nick-Anderssohn/oidc-demo/internal/deps"
	"github.com/Nick-Anderssohn/oidc-demo/internal/sqlc/dal"
	"github.com/Nick-Anderssohn/oidc-demo/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
)

const sessionCookieName = "session_id"

type contextKey string

const sessionContextKey contextKey = "session"
const userIDContextKey contextKey = "user_id"
const sessionLifetimeDays = 1

type Service struct {
	Resolver *deps.Resolver
}

// Middleware to check for a session cookie and add it to the request context
func (s *Service) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the session cookie
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil {
			// If cookie is missing, continue without session
			next.ServeHTTP(w, r)
			return
		}

		sessionRecord, err := s.Resolver.Queries.GetSession(r.Context(), cookie.Value)
		if err != nil {
			// If session is invalid, continue without session
			return
		}

		// Check if session is expired
		sessionStartedAt := sessionRecord.CreatedAt.Time
		expiration := sessionStartedAt.Add(sessionLifetimeDays * 24 * time.Hour)

		if sessionStartedAt.After(expiration) {
			// If session is expired, delete it and continue without session
			if err := s.Resolver.Queries.DeleteSession(r.Context(), cookie.Value); err != nil {
				http.Error(w, "Failed to delete expired session", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:   sessionCookieName,
				MaxAge: -1, // Delete the cookie
			})
			next.ServeHTTP(w, r)
			return
		}

		// Add session and user IDs to request context
		ctx1 := context.WithValue(r.Context(), sessionContextKey, cookie.Value)
		ctx2 := context.WithValue(ctx1, userIDContextKey, sessionRecord.UserID.String())

		next.ServeHTTP(w, r.WithContext(ctx2))
	})
}

// Middleware to enforce session authentication
func (s *Service) RequireSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve session ID from context
		sessionID, ok := r.Context().Value(sessionContextKey).(string)
		if !ok || sessionID == "" {
			http.Error(w, "Unauthorized: No session found", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

func UserIDFromContext(ctx context.Context) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(ctx.Value(userIDContextKey)); err != nil {
		return uuid, err
	}
	return uuid, nil
}

func (s *Service) SaveNewSessionCookie(ctx context.Context, userID pgtype.UUID, w http.ResponseWriter) error {
	sessionId, err := util.GenerateSecureID()
	if err != nil {
		return err
	}

	// Create a new session in the database
	err = s.Resolver.Queries.InsertSession(ctx, dal.InsertSessionParams{
		ID:     sessionId,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	var secure bool
	if strings.HasPrefix(s.Resolver.Config.APIConfig.BaseURL, "https://") {
		secure = true
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionId,
		Expires:  time.Now().Add(sessionLifetimeDays * 24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Secure:   secure,
		HttpOnly: true,
	})

	return nil
}

func (s *Service) Logout(ctx context.Context, w http.ResponseWriter) error {
	sessionID, ok := ctx.Value(sessionContextKey).(string)
	if !ok || sessionID == "" {
		return nil
	}

	// Delete the session from the database
	err := s.Resolver.Queries.DeleteSession(ctx, sessionID)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:   sessionCookieName,
		MaxAge: -1, // Delete the cookie
	})

	return nil
}
