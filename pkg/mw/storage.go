package mw

import (
	"context"

	"github.com/deixis/spine/net/http"
	"github.com/deixis/storage/kvdb"
)

type Store struct {
	S kvdb.Store
}

func (mw *Store) Inject(next http.ServeFunc) http.ServeFunc {
	return func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		ctx = kvdb.WithContext(ctx, mw.S)
		next(ctx, w, req)
	}
}
