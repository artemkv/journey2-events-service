package server

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// Starts serving requests on a specified port with graceful shutdown support
// Blocks the calling thread
// port is a string in Gin format, e.g. ":8600"
//    when port is an empty string, serves on default HTTP port
// callback is called after the server has been setup to serve
//    callback is passed the actual port the server is listening on
func Serve(router *gin.Engine, port string, callback func()) {
	// based on example from https://github.com/gin-gonic/examples
	ctx, restoreInterrupt := getNotifyContextForInterruptSignals()
	defer restoreInterrupt()

	httpServer := startServingAsync(router, port)
	if callback != nil {
		callback()
	}

	waitForInterruptSignal(ctx)
	restoreInterrupt()
	shutDownWithTimeout(httpServer, 5*time.Second)
}

func getNotifyContextForInterruptSignals() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}

func waitForInterruptSignal(ctx context.Context) {
	<-ctx.Done()
}

func startServingAsync(router *gin.Engine, port string) http.Server {
	log.Printf("Starting server on port %s", port)

	httpServer := http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error serving: %s\n", err)
		}
	}()

	return httpServer
}

func shutDownWithTimeout(httpServer http.Server, timeout time.Duration) {
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
