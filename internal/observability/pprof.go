package observability

import (
	"context"
	"course-go-avito-SitnikovArtem06/internal/logger"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

func StartPprof(addr string, logger logger.Logger) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.DefaultServeMux)

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalf("pprof listen: %v", err)
		}
		logger.Log(fmt.Sprintf("pprof on http://%s/debug/pprof/", addr))
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log(fmt.Sprintf("pprof serve: %v", err))
			os.Exit(1)
		}
	}()

	return srv
}

func StopServer(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
