package restfibonacciserver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type App struct {
	asyncStores map[string]*AsyncStore
}

func NewApp() *App {
	return &App{
		asyncStores: make(map[string]*AsyncStore),
	}
}

func (a *App) Start() {
	router := mux.NewRouter()
	router.Use(generateRequestIDMiddleware)
	router.HandleFunc("/fibonacci/sync/{number:[0-9]+}", a.fibonacciSyncHandler)
	router.HandleFunc("/fibonacci/async/{number:[0-9]+}", a.fibonacciAsyncHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Server is running on :8080")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "listen: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down REST fibonacci server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func generateRequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("request-id") == "" {
			requestID := uuid.New().String()
			r.Header.Add("request-id", requestID)
		}
		next.ServeHTTP(w, r)
	})
}
