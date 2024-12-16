package api

import (
	_ "auth_rest/docs"
	"auth_rest/internal/api/handler"
	"auth_rest/internal/services"
	"auth_rest/internal/utils"
	"context"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"time"
)

const (
	maxHeaderBytes = 1 << 20 // 1 MB
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
	router     *mux.Router
}

func NewServer() *Server {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return &Server{
		httpServer: &http.Server{
			Addr:           ":8080",
			MaxHeaderBytes: maxHeaderBytes,
			ReadTimeout:    readTimeout,
			WriteTimeout:   writeTimeout,
			Handler:        router,
		},
		router: router,
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) HandleAuth(service services.AuthService, ctx *utils.AppContext) {
	s.router.HandleFunc("/auth/login/", handler.Login(service, ctx)).Methods(http.MethodPost)
	s.router.HandleFunc("/auth/refresh/", handler.Refresh(service, ctx)).Methods(http.MethodPost)
}
