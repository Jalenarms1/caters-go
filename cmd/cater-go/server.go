package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/Jalenarms1/caters-go/internal/handlers"
	"github.com/Jalenarms1/caters-go/internal/types"
)

type Server struct {
	Addr string
	mux  *http.ServeMux
}

func NewServer(addr string) *Server {
	return &Server{
		Addr: addr,
		mux:  http.NewServeMux(),
	}
}

func (s *Server) registerRoutes() {

	s.mux.HandleFunc("POST /api/login", s.catchErrHandlerFunc(handlers.HandleNewAccount))

}

func (s *Server) run() {
	s.registerRoutes()

	fmt.Printf("http://localhost%s", s.Addr)

	middlewareMux := corsMiddleware(s.mux)

	log.Fatal(http.ListenAndServe(s.Addr, middlewareMux))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("CLIENT_DOMAIN"))
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) catchErrHandlerFunc(fn types.ErrHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			fmt.Println(err)
			slog.Info("Error ocurred", "path", r.URL.Path, "req addr", r.RemoteAddr)
			http.Error(w, err.Err.Error(), err.ReturnCode)
		}
	}
}
