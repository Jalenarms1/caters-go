package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/Jalenarms1/caters-go/internal/handlers"
	"github.com/Jalenarms1/caters-go/internal/types"
	"github.com/golang-jwt/jwt/v5"
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

	s.mux.HandleFunc("POST /api/login", s.catchErrHandlerFunc(handlers.HandleLogin, false))
	s.mux.HandleFunc("POST /api/signup", s.catchErrHandlerFunc(handlers.HandleNewAccount, false))
	s.mux.HandleFunc("POST /api/login/v2", s.catchErrHandlerFunc(handlers.HandleLoginV2, false))
	s.mux.HandleFunc("POST /api/signup/v2", s.catchErrHandlerFunc(handlers.HandleSignupV2, false))
	s.mux.HandleFunc("GET /api/get-me", s.catchErrHandlerFunc(handlers.HandleGetMe, false))

	s.mux.HandleFunc("POST /api/food-shop", s.catchErrHandlerFunc(handlers.HandleNewFoodShop, true))
	s.mux.HandleFunc("POST /api/food-shop-schedule", s.catchErrHandlerFunc(handlers.HandlerNewFoodShopSchedule, true))
	s.mux.HandleFunc("GET /api/food-shop-schedule/{foodShopId}", s.catchErrHandlerFunc(handlers.HandleGetFoodShopSchedule, true))
	s.mux.HandleFunc("DELETE /api/food-shop-schedule/{slotId}", s.catchErrHandlerFunc(handlers.HandleDeleteScheduleSlot, true))

	// s.mux.HandleFunc("/pages/login", s.catchErrHandlerFunc(handlers.HandleLoginPage))
	// s.mux.HandleFunc("/pages/signup", s.catchErrHandlerFunc(handlers.HandleSignupPage))

}

func (s *Server) run() {

	s.registerRoutes()

	fmt.Printf("http://localhost%s\n", s.Addr)

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	fs := http.FileServer(http.Dir(path.Join(wd, "public")))

	s.mux.Handle("/public/", http.StripPrefix("/public", fs))

	middlw := corsMiddleware(s.mux)

	log.Fatal(http.ListenAndServe(s.Addr, middlw))
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

		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		var ctx context.Context
		if len(authHeader) > 1 {
			userKeyToken := authHeader[1]
			fmt.Println(userKeyToken)
			claims := &types.Claims{}
			jwtToken, err := jwt.ParseWithClaims(userKeyToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Printf("User Id: %s", claims.Uid)
			fmt.Println(jwtToken.Valid)
			fmt.Println(claims.Uid)
			if claims.Uid != "" && jwtToken.Valid {
				fmt.Println("claims")
				fmt.Println(claims)

				ctx = context.WithValue(context.Background(), types.AuthKey, claims.Uid)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)

	})
}

func (s *Server) catchErrHandlerFunc(fn types.ErrHandlerFunc, isProtected bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if isProtected {
			uid := r.Context().Value(types.AuthKey)
			fmt.Println("user id")
			fmt.Println(uid)
			if uid == nil {
				http.Error(w, "no authentication", http.StatusForbidden)
				return
			}
		}

		if err := fn(w, r); err != nil {
			fmt.Println(err)
			slog.Info("Error ocurred", "path", r.URL.Path, "req addr", r.RemoteAddr, "Error", err.Err.Error())
			http.Error(w, err.Err.Error(), err.ReturnCode)
		}
	}
}
