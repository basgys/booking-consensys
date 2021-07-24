package auth

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/basgys/booking-consensys/app/iam"
	"github.com/basgys/booking-consensys/pkg/jwtutil"
	"github.com/deixis/errors"
	"github.com/deixis/errors/httperrors"
	"github.com/deixis/pkg/utc"
	"github.com/deixis/spine/net/http"
	"github.com/form3tech-oss/jwt-go"
)

func (s *Service) HandleHTTP(srv *http.Server) {
	h := httpHandler{
		svc: s,
	}

	// Attach middlewares
	jwtm := &jwtutil.HTTPMiddleware{Secret: s.Secret}
	srv.Append(jwtm.Parse)
	srv.Append(h.accountMiddleware)

	// Mount endpoint
	srv.HandleFunc("/auth/challenge", http.POST, h.challenge)
	srv.HandleFunc("/auth/authorise", http.POST, h.authorise)
}

type httpHandler struct {
	svc *Service
}

type httpChallengeRequest struct {
	Address *iam.Address `json:"address"`
}

type httpChallengeResponse struct {
	Challenge string `json:"challenge"`
}

func (h *httpHandler) challenge(
	ctx context.Context, w http.ResponseWriter, req *http.Request,
) {
	defer req.HTTP.Body.Close()
	r := httpChallengeRequest{}
	if err := unmarshalJSON(req.HTTP.Body, &r); err != nil {
		httperrors.Marshal(req.HTTP, w, err)
		return
	}

	ch, err := h.svc.Challenge(ctx, r.Address)
	if err != nil {
		httperrors.Marshal(req.HTTP, w, err)
		return
	}

	w.JSON(http.StatusCreated, &httpChallengeResponse{
		Challenge: ch,
	})
}

type httpAuthoriseRequest struct {
	Address   *iam.Address `json:"address"`
	Signature string       `json:"signature"`
}

type httpAuthoriseResponse struct {
	Token string `json:"token"`
}

func (h *httpHandler) authorise(
	ctx context.Context, w http.ResponseWriter, req *http.Request,
) {
	defer req.HTTP.Body.Close()
	r := httpAuthoriseRequest{}
	if err := unmarshalJSON(req.HTTP.Body, &r); err != nil {
		httperrors.Marshal(req.HTTP, w, err)
		return
	}

	token, err := h.svc.Authorise(ctx, r.Address, r.Signature)
	if err != nil {
		httperrors.Marshal(req.HTTP, w, err)
		return
	}

	w.JSON(http.StatusOK, &httpAuthoriseResponse{
		Token: token,
	})
}

func (h *httpHandler) accountMiddleware(next http.ServeFunc) http.ServeFunc {
	return func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		token, ok := jwtutil.FromContext(ctx)
		if !ok {
			// Unauthenticated request
			next(ctx, w, req)
			return
		}

		// Lengthy way to load account key from JWT
		if err := token.Claims.Valid(); err != nil {
			httperrors.Marshal(req.HTTP, w, err)
			return
		}
		claims, ok := token.Claims.(*jwt.StandardClaims)
		if !ok {
			httperrors.Marshal(req.HTTP, w, errors.Bad(&errors.FieldViolation{
				Field:       "jwt",
				Description: "Invalid JWT claims format. Expect standard claims",
			}))
			return
		}
		if claims.Id == "" {
			httperrors.Marshal(req.HTTP, w, errors.Bad(&errors.FieldViolation{
				Field:       "id",
				Description: "Missing ID in JWT claims",
			}))
			return
		}
		if claims.ExpiresAt < int64(utc.Now()) {
			httperrors.Marshal(req.HTTP, w, errors.Unauthenticated)
			return
		}

		addr, err := iam.ParseAddress(claims.Id)
		if err != nil {
			httperrors.Marshal(req.HTTP, w, err)
			return
		}

		// Load account from repository
		acc, err := h.svc.accounts.Get(ctx, addr)
		switch {
		case err == nil:
			// Good
		case errors.IsNotFound(err):
			// The JWT points to an account that does not exist
			httperrors.Marshal(req.HTTP, w, errors.Unauthenticated)
			return
		default:
			httperrors.Marshal(req.HTTP, w, err)
			return
		}

		// Attach account to context
		ctx = iam.WithContext(ctx, acc)
		next(ctx, w, req)
	}
}

func unmarshalJSON(r io.Reader, v interface{}) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, v); err != nil {
		return errors.WithBad(err)
	}
	return nil
}
