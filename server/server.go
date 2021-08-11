package server

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// port is a string in Gin format, e.g. ":8600"
func Serve(port string) {
	router := gin.Default() // logging and recovery attached
	setupRouter(router)

	serveWithShutdownSupport(router, port)
}

// based on example from https://github.com/gin-gonic/examples
func serveWithShutdownSupport(router *gin.Engine, port string) {
	ctx, cancel := getNotifyContextForInterruptSignals()
	defer cancel()

	httpServer := startListeningAsync(router, port)

	waitForInterruptSignal(ctx)
	cancel() // restore default behavior on the interrupt signal
	shutDownWithTimeout(httpServer, 5*time.Second)
}

func getNotifyContextForInterruptSignals() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
}

func waitForInterruptSignal(ctx context.Context) {
	<-ctx.Done()
}

func startListeningAsync(router *gin.Engine, port string) http.Server {
	httpServer := http.Server{
		Addr:    port,
		Handler: router,
	}

	log.Printf("Starting listening on port %s", port)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error listening: %s\n", err)
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
