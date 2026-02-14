package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	sessionInfra "github.com/Jaxongir1006/Chat-X-v2/internal/infra/postgres/repo/session"
	"github.com/rs/zerolog"
)

type ctxKey string

type metaKey string

const (
	CtxUserID    ctxKey = "user_id"
	CtxSessionID ctxKey = "session_id"
	CtxIP        metaKey = "ip"
	CtxUserAgent metaKey = "user_agent"
	CtxDevice    metaKey = "device"
)

type RequestMeta struct {
	IP        string
	UserAgent string
	Device    string
}

func WithMeta(ctx context.Context, m RequestMeta) context.Context {
	ctx = context.WithValue(ctx, CtxIP, m.IP)
	ctx = context.WithValue(ctx, CtxUserAgent, m.UserAgent)
	ctx = context.WithValue(ctx, CtxDevice, m.Device)
	return ctx
}

func MetaFromContext(ctx context.Context) (RequestMeta, bool) {
	ip, ok1 := ctx.Value(CtxIP).(string)
	ua, ok2 := ctx.Value(CtxUserAgent).(string)
	device, ok3 := ctx.Value(CtxDevice).(string)
	if !ok1 || !ok2 || !ok3 {
		return RequestMeta{}, false
	}
	return RequestMeta{IP: ip, UserAgent: ua, Device: device}, true
}

func MetaMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.UserAgent()
		meta := RequestMeta{
			IP:        clientIP(r),
			UserAgent: ua,
			Device:    deviceLabel(ua),
		}
		next.ServeHTTP(w, r.WithContext(WithMeta(r.Context(), meta)))
	})
}

type AuthMiddleware struct {
	Sessions       sessionInfra.SessionStore
	RequireRefresh bool
}

func NewAuthMiddleware(sessions sessionInfra.SessionStore, requireRefresh bool) *AuthMiddleware {
	return &AuthMiddleware{Sessions: sessions, RequireRefresh: requireRefresh}
}

func (m *AuthMiddleware) WrapAccess(next http.Handler) http.Handler {
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
		if sess.RevokedAt != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		now := time.Now()
		if now.After(sess.AccessTokenExp) {
			http.Error(w, "access token expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), CtxUserID, sess.UserID)
		ctx = context.WithValue(ctx, CtxSessionID, sess.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) WrapRefresh(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		sess, err := m.Sessions.GetByRefreshToken(r.Context(), refresh)
		if err != nil || sess == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if sess.RevokedAt != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		now := time.Now()
		if now.After(sess.RefreshTokenExp) {
			http.Error(w, "refresh token expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), CtxUserID, sess.UserID)
		ctx = context.WithValue(ctx, CtxSessionID, sess.ID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


// logging
type responseWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += n
	return n, err
}

func Logging(l zerolog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(rw, r)

		d := time.Since(start)

		evt := l.Info()
		switch {
		case rw.status >= 500:
			evt = l.Error()
		case rw.status >= 400:
			evt = l.Warn()
		}

		meta, ok := MetaFromContext(r.Context())
		if !ok {
			ua := r.UserAgent()
			meta = RequestMeta{IP: clientIP(r), UserAgent: ua, Device: deviceLabel(ua)}
		}

		evt.
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", rw.status).
			Int("bytes", rw.bytes).
			Int64("dur_ms", d.Milliseconds()).
			Str("ip", meta.IP).
			Str("ua", meta.UserAgent).
			Str("device", meta.Device).
			Msg("http")
	})
}

// Helpers
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

func clientIP(r *http.Request) string {
	// X-Forwarded-For: "client, proxy1, proxy2"
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		first := strings.TrimSpace(strings.Split(xff, ",")[0])
		if first != "" {
			return first
		}
	}
	if xrip := strings.TrimSpace(r.Header.Get("X-Real-IP")); xrip != "" {
		return xrip
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && host != "" {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

func UserIDFromContext(ctx context.Context) (uint64, bool) {
	v := ctx.Value(CtxUserID)
	id, ok := v.(uint64)
	return id, ok
}

func parseClient(ua string) string {
	u := strings.ToLower(ua)
	switch {
	case strings.Contains(u, "postman"):
		return "Postman"
	case strings.Contains(u, "chrome") && !strings.Contains(u, "edg"):
		return "Chrome"
	case strings.Contains(u, "safari") && !strings.Contains(u, "chrome"):
		return "Safari"
	case strings.Contains(u, "firefox"):
		return "Firefox"
	case strings.Contains(u, "edg"):
		return "Edge"
	default:
		return "Unknown Client"
	}
}

func deviceLabel(ua string) string {
	return parseClient(ua) + " on " + parseDevice(ua)
}

func parseDevice(ua string) string {
	ua = strings.ToLower(ua)
	switch {
	case strings.Contains(ua, "android"):
		return "Android"
	case strings.Contains(ua, "iphone"), strings.Contains(ua, "ipad"), strings.Contains(ua, "ios"):
		return "iOS"
	case strings.Contains(ua, "windows"):
		return "Windows"
	case strings.Contains(ua, "macintosh"), strings.Contains(ua, "mac os"):
		return "Mac"
	case strings.Contains(ua, "linux"):
		return "Linux"
	default:
		return "Unknown"
	}
}
