package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"

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

// asciiReader wraps an underlying io.Reader, and munges all read
// bytes to be ASCII non-control characters.
type asciiReader struct {
	inner io.Reader
}

func (a asciiReader) Read(dat []byte) (int, error) {
	n, err := a.inner.Read(dat)
	for i := 0; i < n; i++ {
		// the bit 8 is always zero for ascii
		dat[i] &^= 0b1000_0000
		// bits 6 and 7 are only ever clear for control
		// characters, ensure one of them is set
		if dat[i]&0b0110_0000 == 0 {
			dat[i] |= 0b0100_0000
		}
	}
	return n, err
}

func Main(ctx context.Context) error {

	sc := &dhttp.ServerConfig{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var requestID [8]byte
			_, _ = rand.Read(requestID[:])
			ctx := dlog.WithField(r.Context(), "request_id", fmt.Sprintf("%0x", requestID))

			dataSource := io.Reader(asciiReader{rand.Reader})
			if size, err := strconv.ParseInt(r.URL.Query().Get("size"), 10, 0); err == nil && size > 0 {
				w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
				dataSource = io.LimitReader(dataSource, size)
			}

			dlog.Infof(ctx, "begin: %s %s %s", r.Method, r.RequestURI, r.Proto)

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			if flusher, ok := w.(http.Flusher); ok {
				// don't buffer the request to auto-determine Content-Length
				flusher.Flush()
			}
			n, _ := io.Copy(w, dataSource)

			dlog.Infof(ctx, "end: wrote %d bytes", n)
		}),
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	dlog.Infof(ctx, "started up and listening on %v", listener.Addr())
	return sc.Serve(ctx, listener)
}
