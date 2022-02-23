package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/datawire/dlib/dhttp"
	"github.com/datawire/dlib/dlog"
)

func main() {
	ctx := context.Background()
	if err := Main(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s: error: %v\n", os.Args[0], err)
		os.Exit(1)
	}
}

func Main(ctx context.Context) error {
	sc := &dhttp.ServerConfig{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = io.Copy(w, rand.Reader)
		}),
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	dlog.Infof(ctx, "started up and listening on %v", listener.Addr())
	return sc.Serve(ctx, listener)
}
