package mw

import (
	"context"

	"github.com/deixis/spine/net/http"
)

// NodeInfo is a middelware that adds information about the node serving the request
type NodeInfo struct {
	Name    string
	Version string
}

// ReturnNodeInfo adds information about the node serving the request on the response header
func (mw *NodeInfo) ReturnNodeInfo(next http.ServeFunc) http.ServeFunc {
	return func(ctx context.Context, w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Node-Name", mw.Name)
		w.Header().Set("Node-Version", mw.Version)
		next(ctx, w, req)
	}
}
