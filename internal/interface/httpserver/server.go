package httpserver

import (
	"context"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/di"
	"github.com/zuxt268/berry/internal/interface/handlers"
	xmiddleware "github.com/zuxt268/berry/internal/interface/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	frontend "github.com/zuxt268/berry/frontend"
)

type Server struct {
	dbClose              func()
	server               *http.Server
	router               *chi.Mux
	authMiddleware       *xmiddleware.AuthMiddleware
	userAuthHandler      *handlers.UserAuthHandler
	operatorAuthHandler  *handlers.OperatorAuthHandler
	ga4AuthHandler       *handlers.GA4AuthHandler
	gscAuthHandler       *handlers.GSCAuthHandler
	gbpAuthHandler       *handlers.GBPAuthHandler
	instagramAuthHandler *handlers.InstagramAuthHandler
	lineAuthHandler      *handlers.LineAuthHandler
	userHandler          *handlers.UserHandler
	operatorHandler      *handlers.OperatorHandler
	reportHandler        *handlers.ReportHandler
}

func NewServer(
	addr string,
	container *di.Container,
) *Server {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{config.Env.FrontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	srv := &Server{
		router:               r,
		authMiddleware:       container.AuthMiddleware,
		dbClose:              container.DBClose,
		userHandler:          container.UserHandler,
		operatorHandler:      container.OperatorHandler,
		userAuthHandler:      container.UserAuthHandler,
		operatorAuthHandler:  container.OperatorAuthHandler,
		ga4AuthHandler:       container.GA4AuthHandler,
		gscAuthHandler:       container.GSCAuthHandler,
		gbpAuthHandler:       container.GBPAuthHandler,
		instagramAuthHandler: container.InstagramAuthHandler,
		lineAuthHandler:      container.LineAuthHandler,
		reportHandler:        container.ReportHandler,
		server: &http.Server{
			Addr:    addr,
			Handler: r,
		},
	}

	srv.registerRoutes()
	_ = srv.setupStaticFileServer()

	return srv
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.dbClose()
	return s.server.Shutdown(ctx)
}

// registerRoutes registers all HTTP routes
func (s *Server) registerRoutes() {
	s.router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	s.router.Route("/api", func(r chi.Router) {
		// 認証必要エンドポイント
		r.Group(func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Route("/auth/google", func(r chi.Router) {
					r.Get("/login", s.userAuthHandler.GoogleLogin)
					r.Get("/callback", s.userAuthHandler.GoogleCallback)
					r.Get("/logout", s.userAuthHandler.GoogleLogout)
				})

				r.Group(func(r chi.Router) {
					r.Use(s.authMiddleware.RequireAuth)
					r.Get("/", s.userHandler.Gets)
					r.Get("/{uid}", s.userHandler.GetByUID)
					r.Post("/", s.userHandler.Create)
					r.Put("/{uid}", s.userHandler.Update)
					r.Delete("/{uid}", s.userHandler.Delete)
				})

				r.Route("/ga4", func(r chi.Router) {
					r.Group(func(r chi.Router) {
						r.Get("/auth/connect", s.ga4AuthHandler.Connect)
						r.Get("/auth/callback", s.ga4AuthHandler.Callback)
						r.Get("/connections", s.ga4AuthHandler.GetConnections)
						r.Delete("/connections/{uid}", s.ga4AuthHandler.Disconnect)
					})
					r.Group(func(r chi.Router) {
						r.Use(s.authMiddleware.RequireAuth)
						r.Get("/reports", s.reportHandler.GA4Reports)
					})
				})

				r.Route("/gsc", func(r chi.Router) {
					r.Group(func(r chi.Router) {
						r.Get("/auth/connect", s.gscAuthHandler.Connect)
						r.Get("/auth/callback", s.gscAuthHandler.Callback)
						r.Get("/connections", s.gscAuthHandler.GetConnections)
						r.Delete("/connections/{uid}", s.gscAuthHandler.Disconnect)
					})
					r.Group(func(r chi.Router) {
						r.Use(s.authMiddleware.RequireAuth)
						r.Get("/reports", s.reportHandler.GSCReports)
					})
				})

				r.Route("/gbp", func(r chi.Router) {
					r.Group(func(r chi.Router) {
						r.Get("/auth/connect", s.gbpAuthHandler.Connect)
						r.Get("/auth/callback", s.gbpAuthHandler.Callback)
						r.Get("/connections", s.gbpAuthHandler.GetConnections)
						r.Delete("/connections/{uid}", s.gbpAuthHandler.Disconnect)
					})
					r.Group(func(r chi.Router) {
						r.Use(s.authMiddleware.RequireAuth)
						r.Get("/reports", s.reportHandler.GBPReports)
					})
				})

				r.Route("/instagram", func(r chi.Router) {
					r.Group(func(r chi.Router) {
						r.Get("/auth/connect", s.instagramAuthHandler.Connect)
						r.Get("/auth/callback", s.instagramAuthHandler.Callback)
						r.Get("/connections", s.instagramAuthHandler.GetConnections)
						r.Delete("/connections/{uid}", s.instagramAuthHandler.Disconnect)
					})
					r.Group(func(r chi.Router) {
						r.Use(s.authMiddleware.RequireAuth)
						r.Get("/reports", s.reportHandler.InstagramReports)
					})
				})

				r.Route("/line", func(r chi.Router) {
					r.Use(s.authMiddleware.RequireAuth)
					r.Post("/connect", s.lineAuthHandler.Connect)
					r.Get("/connections", s.lineAuthHandler.GetConnections)
					r.Delete("/connections/{uid}", s.lineAuthHandler.Disconnect)
					r.Get("/reports", s.reportHandler.LineReports)
				})
			})
			r.Route("/operators", func(r chi.Router) {
				r.Route("/auth/google", func(r chi.Router) {
					r.Get("/login", s.operatorAuthHandler.GoogleLogin)
					r.Get("/callback", s.operatorAuthHandler.GoogleCallback)
					r.Get("/logout", s.operatorAuthHandler.GoogleLogout)
				})
				r.Group(func(r chi.Router) {
					r.Use(s.authMiddleware.RequireOperatorAuth)
					r.Get("/", s.operatorHandler.Gets)
					r.Get("/{uid}", s.operatorHandler.GetByUID)
					r.Post("/", s.operatorHandler.Create)
					r.Put("/{uid}", s.operatorHandler.Update)
					r.Delete("/{uid}", s.operatorHandler.Delete)
				})
			})

		})
	})
}

func (s *Server) setupStaticFileServer() error {
	staticFiles := frontend.GetStaticFiles()
	fsys, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		return err
	}
	fileServer := http.FileServer(http.FS(fsys))
	s.router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "index.html"
		}
		_, err := fsys.Open(strings.TrimPrefix(path, "/"))
		if os.IsNotExist(err) {
			// ファイルが存在しない場合はindex.htmlを返す（SPA用）
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
	return nil
}