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
	r.Get("/api/summer-shifts", h.summerShifts)
	r.Get("/api/summer-shifts/{id}", h.summerShift)
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
	admin.Post("/media", h.uploadMedia)
	admin.Delete("/media/{id}", h.deleteMedia)
	admin.Post("/classes", h.saveAdminClass)
	admin.Get("/classes/{id}/form", h.adminClassForm)
	admin.Put("/classes/{id}/form", h.saveAdminClass)
	admin.Delete("/classes/{id}/media/{mediaId}", h.deleteAdminClassMedia)
	admin.Post("/summer-shifts", h.saveAdminSummerShift)
	admin.Get("/summer-shifts/{id}/form", h.adminSummerShiftForm)
	admin.Put("/summer-shifts/{id}/form", h.saveAdminSummerShift)
	admin.Delete("/summer-shifts/{id}/media/{mediaId}", h.deleteAdminSummerShiftMedia)
	admin.Post("/news", h.saveAdminNews)
	admin.Get("/news/categories", h.adminNewsCategories)
	admin.Get("/news/{id}/form", h.adminNewsForm)
	admin.Put("/news/{id}/form", h.saveAdminNews)
	admin.Delete("/news/{id}/media/{mediaId}", h.deleteAdminNewsMedia)
	admin.Post("/staff", h.saveAdminStaff)
	admin.Get("/staff/{id}/form", h.adminStaffForm)
	admin.Put("/staff/{id}/form", h.saveAdminStaff)
	admin.Delete("/staff/{id}/media/{mediaId}", h.deleteAdminStaffMedia)
	admin.Post("/documents", h.saveAdminDocument)
	admin.Get("/documents/{id}/form", h.adminDocumentForm)
	admin.Put("/documents/{id}/form", h.saveAdminDocument)
	admin.Post("/reviews", h.saveAdminReview)
	admin.Get("/reviews/{id}/form", h.adminReviewForm)
	admin.Put("/reviews/{id}/form", h.saveAdminReview)
	admin.Get("/gallery", h.adminGallery)
	admin.Post("/gallery", h.addAdminGallery)
	admin.Delete("/gallery/{id}", h.deleteAdminGallery)
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
