package intercepts

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
)

func Recover(ctx context.Context, spec connect.Spec, header http.Header, req any) error {
	return nil
}
