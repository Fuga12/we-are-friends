package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Fuga12/we-are-friends/backend-go/internal/config"
	"github.com/Fuga12/we-are-friends/backend-go/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// NewRouter wires cross-cutting HTTP middleware and API routes. Entity
// handlers are added here as their store queries are implemented.
func NewRouter(cfg config.Config, database *store.Store, media mediaURLs) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	h := handler{store: database, media: media}
	r.Get("/api", healthz)
	r.Get("/api/healthz", healthz)
	r.Get("/api/home", h.home)
	r.Get("/api/news", h.news)
	r.Get("/api/news/{id}", h.newsItem)
	r.Get("/api/classes", h.classes)
	r.Get("/api/classes/{id}", h.class)
	r.Get("/api/services", h.services)
	r.Get("/api/services/{id}", h.service)
	r.Get("/api/staff", h.staff)
	r.Get("/api/documents", h.documents)
	r.Get("/api/gallery", h.gallery)
	r.Get("/api/reviews", h.reviews)
	r.Get("/api/site", h.site)

	r.Post("/api/admin/login", h.login)
	admin := chi.NewRouter()
	admin.Use(h.requireAuth)
	admin.Post("/logout", h.logout)
	admin.Get("/overview", h.adminOverview)
	admin.Get("/media", h.mediaItems)
	admin.Post("/media", h.uploadMedia)
	admin.Delete("/media/{id}", h.deleteMedia)
	admin.Get("/{resource}", h.adminItems)
	admin.Get("/{resource}/{id}/content", h.adminContent)
	admin.Patch("/{resource}/{id}/content", h.updateAdminContent)
	admin.Patch("/{resource}/{id}", h.setPublication)
	admin.Put("/{resource}/{id}", h.updateAdminItem)
	admin.Delete("/{resource}/{id}", h.deleteAdminItem)
	r.Mount("/api/admin", admin)
	return r
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
