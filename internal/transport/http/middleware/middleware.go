package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	sessionInfra "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/session"
)

type ctxKey string

const CtxUserID ctxKey = "user_id"
const CtxSessionID ctxKey = "session_id"

type AuthMiddleware struct {
	Sessions       sessionInfra.SessionStore
	RequireRefresh bool // set true only if you REALLY want refresh on every request
}

func NewAuthMiddleware(sessions sessionInfra.SessionStore, requireRefresh bool) *AuthMiddleware {
	return &AuthMiddleware{
		Sessions:       sessions,
		RequireRefresh: requireRefresh,
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (m AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		access := extractBearer(r.Header.Get("Authorization"))
		if access == "" {
			if c, err := r.Cookie("access_token"); err == nil {
				access = c.Value
			}
		}
		if access == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		sess, err := m.Sessions.GetByAccessToken(r.Context(), access)
		if err != nil || sess == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		now := time.Now()
		if now.After(sess.AccessTokenExp) {
			http.Error(w, "access token expired", http.StatusUnauthorized)
			return
		}

		if m.RequireRefresh {
			refresh := ""
			if c, err := r.Cookie("refresh_token"); err == nil {
				refresh = c.Value
			}
			if refresh == "" {
				refresh = r.Header.Get("X-Refresh-Token")
			}
			if refresh == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			rs, err := m.Sessions.GetByRefreshToken(r.Context(), refresh)
			if err != nil || rs == nil || rs.UserID != sess.UserID {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			if now.After(rs.RefreshTokenExp) {
				http.Error(w, "refresh token expired", http.StatusUnauthorized)
				return
			}
		}

		_ = m.Sessions.UpdateLastUsed(r.Context(), sess.ID, now)

		ctx := context.WithValue(r.Context(), CtxUserID, sess.UserID)
		ctx = context.WithValue(ctx, CtxSessionID, sess.ID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractBearer(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: 0}

		next.ServeHTTP(rw, r)

		d := time.Since(start)

		log.Printf(
			`%s %s -> %d (%dB) in %s | ip=%s | %q`,
			r.Method,
			r.URL.Path,
			rw.status,
			rw.bytes,
			d,
			clientIP(r),
			r.UserAgent(),
		)
	})
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return xrip
	}
	return r.RemoteAddr
}
