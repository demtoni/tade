package api

import (
	"database/sql"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/demtoni/tade/internal/config"
	"github.com/demtoni/tade/internal/database"
	"github.com/demtoni/tade/webapp"
	"github.com/rvinnie/yookassa-sdk-go/yookassa"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/gorilla/sessions"
	_ "modernc.org/sqlite"
)

type Server struct {
	config  *config.Config
	router  *chi.Mux
	queries *database.Queries
	store   *sessions.CookieStore
	kassa   *yookassa.PaymentHandler
}

func New(cfg *config.Config) (*Server, error) {
	s := &Server{
		config: cfg,
		router: chi.NewRouter(),
	}

	db, err := sql.Open("sqlite", s.config.PathToDB)
	if err != nil {
		return nil, err
	}
	s.queries = database.New(db)

	s.store = sessions.NewCookieStore([]byte(s.config.SessionSecret))

	s.kassa = yookassa.NewPaymentHandler(yookassa.NewClient(s.config.YooShopID, s.config.YooApiKey))

	s.router.Route("/api", func(r chi.Router) {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
		r.Post("/register", s.Register)
		r.Post("/login", s.Login)
		r.Route("/me", func(r chi.Router) {
			r.Use(s.AuthCtx)
			r.Get("/", s.GetUserInfo)
			r.Put("/password", s.ChangePassword)
			r.Route("/services", func(r chi.Router) {
				r.Get("/", s.ListUserServices)
				r.Get("/locations", s.ListLocations)
				r.Post("/", s.CreateService)
				r.Get("/{id}", s.GetService)
			})
			r.Post("/balance", s.AddBalance)
			r.Get("/transactions", s.GetTransactionList)
			r.Post("/invites", s.GenerateInvite)
			r.Get("/invites", s.ListInvites)
		})
	})

	frontend, _ := fs.Sub(webapp.Content, "dist")
	s.router.Handle("/*", http.FileServer(http.FS(frontend)))

	return s, nil
}

func (s *Server) Run() error {
	paymentsTicker := time.NewTicker(60 * time.Second)
	serviceTicker := time.NewTicker(60 * time.Second)

	defer paymentsTicker.Stop()
	defer serviceTicker.Stop()

	go func() {
		for {
			select {
			case <-paymentsTicker.C:
				if err := s.ProcessPayments(); err != nil {
					log.Println(err)
				}
			case <-serviceTicker.C:
				if err := s.CheckServices(); err != nil {
					log.Println(err)
				}
			}
		}
	}()

	return http.ListenAndServe(s.config.ServerAddr, s.router)
}
