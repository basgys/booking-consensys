package jwtutil

import (
	"context"
	"strings"

	"github.com/deixis/errors"
	"github.com/deixis/errors/httperrors"
	"github.com/deixis/spine/net/http"
	"github.com/form3tech-oss/jwt-go"
)

var (
	// StandardMethod specifies the signature method used by this package
	StandardMethod = jwt.SigningMethodHS256

	// ErrBadSigningMethod is returned when the given alg is wrong
	ErrBadSigningMethod = errors.New("bad signing method")
)

// HTTPMiddleware parses JWT from HTTP header
type HTTPMiddleware struct {
	Secret string
}

// Parse parses JWT from HTTP header and add it to the request context
func (mw *HTTPMiddleware) Parse(next http.ServeFunc) http.ServeFunc {
	return func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		// Extract authorization from header
		authHeader := req.HTTP.Header.Get("Authorization")
		if authHeader == "" {
			// Request has no token
			next(ctx, w, req)
			return
		}

		// Check auth header validity
		authHeaderParts := strings.Fields(authHeader)
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			httperrors.Marshal(req.HTTP, w, errors.Bad(&errors.FieldViolation{
				Field:       "Authorization",
				Description: "Authorization header format must be \"Bearer {token}\"",
			}))
			return
		}
		tokenString := authHeaderParts[1]

		// Parse and validate JWT
		claims := jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrBadSigningMethod
			}
			return []byte(mw.Secret), nil
		})
		if err != nil {
			httperrors.Marshal(req.HTTP, w, errors.PermissionDenied)
			return
		}

		// Add JWT to context and move on
		ctx = WithContext(ctx, token)
		next(ctx, w, req)
	}
}
