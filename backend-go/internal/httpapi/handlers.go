package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Fuga12/we-are-friends/backend-go/internal/store"
	"github.com/go-chi/chi/v5"
)

type mediaURLs interface {
	URL(context.Context, string) (string, error)
}
type mediaStorage interface {
	mediaURLs
	Upload(context.Context, string, string, int64, io.Reader) error
	Delete(context.Context, string) error
}

type handler struct {
	store *store.Store
	media mediaURLs
}
type adminContextKey struct{}
type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type publicationRequest struct {
	IsPublished bool `json:"isPublished"`
}
type contentUpdateRequest struct {
	Title       string `json:"title"`
	IsPublished bool   `json:"isPublished"`
	Position    int    `json:"position"`
}

type classDTO struct {
	ID              int64    `json:"id"`
	Title           string   `json:"title"`
	AgeRange        string   `json:"ageRange"`
	Description     string   `json:"description"`
	Schedule        string   `json:"schedule"`
	Price           string   `json:"price"`
	Emoji           string   `json:"emoji"`
	Color           string   `json:"color"`
	Photo           string   `json:"photo,omitempty"`
	Photos          []string `json:"photos,omitempty"`
	FullDescription []string `json:"fullDescription,omitempty"`
	Highlights      []string `json:"highlights,omitempty"`
}
type newsDTO struct {
	ID       int64    `json:"id"`
	Title    string   `json:"title"`
	Date     string   `json:"date"`
	Category string   `json:"category"`
	Excerpt  string   `json:"excerpt"`
	Emoji    string   `json:"emoji"`
	Photo    string   `json:"photo,omitempty"`
	Photos   []string `json:"photos,omitempty"`
	Content  []string `json:"content,omitempty"`
}
type serviceDTO struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Tagline     string          `json:"tagline"`
	Description string          `json:"description"`
	Features    []string        `json:"features"`
	Price       string          `json:"price"`
	Emoji       string          `json:"emoji"`
	To          string          `json:"to"`
	Color       string          `json:"color"`
	Photo       string          `json:"photo,omitempty"`
	Photos      []string        `json:"photos,omitempty"`
	Detail      json.RawMessage `json:"detail,omitempty"`
}
type staffDTO struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Bio      string `json:"bio"`
	Initials string `json:"initials"`
	Color    string `json:"color"`
	Photo    string `json:"photo,omitempty"`
}
type documentDTO struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	FileType    string `json:"fileType"`
	FileSize    string `json:"fileSize"`
	URL         string `json:"url"`
}
type galleryDTO struct {
	ID    int64  `json:"id"`
	Image string `json:"image"`
}
type reviewDTO struct {
	ID     int64  `json:"id"`
	Author string `json:"author"`
	Role   string `json:"role"`
	Text   string `json:"text"`
	Rating int16  `json:"rating"`
}
type mediaDTO struct {
	ID        int64  `json:"id"`
	URL       string `json:"url"`
	MIME      string `json:"mime"`
	SizeBytes int64  `json:"sizeBytes"`
	Alt       string `json:"alt"`
}
type mediaPageDTO struct {
	Items   []mediaDTO `json:"items"`
	Page    int        `json:"page"`
	PerPage int        `json:"perPage"`
	Total   int64      `json:"total"`
}

func (h handler) classes(w http.ResponseWriter, r *http.Request) {
	values, err := h.store.Classes(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	out := make([]classDTO, 0, len(values))
	for _, v := range values {
		dto, err := h.classDTO(r.Context(), v)
		if err != nil {
			h.failure(w, err)
			return
		}
		out = append(out, dto)
	}
	writeJSON(w, http.StatusOK, out)
}
func (h handler) class(w http.ResponseWriter, r *http.Request) {
	v, err := h.store.Class(r.Context(), pathID(r))
	if err != nil {
		h.failure(w, err)
		return
	}
	dto, err := h.classDTO(r.Context(), v)
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, dto)
}
func (h handler) classDTO(ctx context.Context, v store.Class) (classDTO, error) {
	photo, err := h.oneURL(ctx, v.CoverKey)
	if err != nil {
		return classDTO{}, err
	}
	photos, err := h.urls(ctx, v.Photos)
	if err != nil {
		return classDTO{}, err
	}
	return classDTO{v.ID, v.Title, v.AgeRange, v.Description, v.Schedule, v.Price, v.Emoji, v.Color, photo, photos, v.FullDescription, v.Highlights}, nil
}

func (h handler) news(w http.ResponseWriter, r *http.Request) {
	values, err := h.store.News(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	out := make([]newsDTO, 0, len(values))
	for _, v := range values {
		dto, err := h.newsDTO(r.Context(), v)
		if err != nil {
			h.failure(w, err)
			return
		}
		out = append(out, dto)
	}
	writeJSON(w, http.StatusOK, out)
}
func (h handler) newsItem(w http.ResponseWriter, r *http.Request) {
	v, err := h.store.NewsItem(r.Context(), pathID(r))
	if err != nil {
		h.failure(w, err)
		return
	}
	dto, err := h.newsDTO(r.Context(), v)
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, dto)
}
func (h handler) newsDTO(ctx context.Context, v store.News) (newsDTO, error) {
	photo, err := h.oneURL(ctx, v.CoverKey)
	if err != nil {
		return newsDTO{}, err
	}
	photos, err := h.urls(ctx, v.Photos)
	if err != nil {
		return newsDTO{}, err
	}
	return newsDTO{v.ID, v.Title, v.Date, v.Category, v.Excerpt, v.Emoji, photo, photos, v.Content}, nil
}

func (h handler) services(w http.ResponseWriter, r *http.Request) {
	values, err := h.store.Services(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	out := make([]serviceDTO, 0, len(values))
	for _, v := range values {
		dto, err := h.serviceDTO(r.Context(), v)
		if err != nil {
			h.failure(w, err)
			return
		}
		out = append(out, dto)
	}
	writeJSON(w, http.StatusOK, out)
}
func (h handler) service(w http.ResponseWriter, r *http.Request) {
	v, err := h.store.Service(r.Context(), pathID(r))
	if err != nil {
		h.failure(w, err)
		return
	}
	dto, err := h.serviceDTO(r.Context(), v)
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, dto)
}
func (h handler) serviceDTO(ctx context.Context, v store.Service) (serviceDTO, error) {
	photo, err := h.oneURL(ctx, v.CoverKey)
	if err != nil {
		return serviceDTO{}, err
	}
	photos, err := h.urls(ctx, v.Photos)
	if err != nil {
		return serviceDTO{}, err
	}
	return serviceDTO{v.ID, v.Title, v.Tagline, v.Description, v.Features, v.Price, v.Emoji, v.Path, v.Color, photo, photos, v.Detail}, nil
}

func (h handler) staff(w http.ResponseWriter, r *http.Request) {
	values, err := h.store.Staff(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	out := make([]staffDTO, 0, len(values))
	for _, v := range values {
		photo, err := h.oneURL(r.Context(), v.PhotoKey)
		if err != nil {
			h.failure(w, err)
			return
		}
		out = append(out, staffDTO{v.ID, v.Name, v.Role, v.Bio, v.Initials, v.Color, photo})
	}
	writeJSON(w, http.StatusOK, out)
}
func (h handler) documents(w http.ResponseWriter, r *http.Request) {
	values, err := h.store.Documents(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	out := make([]documentDTO, 0, len(values))
	for _, v := range values {
		url, err := h.media.URL(r.Context(), v.ObjectKey)
		if err != nil {
			h.failure(w, err)
			return
		}
		kind, size := fileType(v.MIME), fileSize(v.SizeBytes)
		if v.DisplayFileType != nil {
			kind = *v.DisplayFileType
		}
		if v.DisplayFileSize != nil {
			size = *v.DisplayFileSize
		}
		out = append(out, documentDTO{v.ID, v.Title, v.Description, kind, size, url})
	}
	writeJSON(w, http.StatusOK, out)
}
func (h handler) gallery(w http.ResponseWriter, r *http.Request) {
	values, err := h.store.Gallery(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	out := make([]galleryDTO, 0, len(values))
	for _, v := range values {
		url, err := h.media.URL(r.Context(), v.ObjectKey)
		if err != nil {
			h.failure(w, err)
			return
		}
		out = append(out, galleryDTO{v.ID, url})
	}
	writeJSON(w, http.StatusOK, out)
}
func (h handler) reviews(w http.ResponseWriter, r *http.Request) {
	values, err := h.store.Reviews(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	out := make([]reviewDTO, 0, len(values))
	for _, v := range values {
		out = append(out, reviewDTO{v.ID, v.Author, v.Role, v.Text, v.Rating})
	}
	writeJSON(w, http.StatusOK, out)
}
func (h handler) site(w http.ResponseWriter, r *http.Request) {
	data, err := h.store.Site(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(data)
}

func (h handler) adminOverview(w http.ResponseWriter, r *http.Request) {
	counts, err := h.store.AdminCounts(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{
		"classes": counts.Classes, "news": counts.News, "services": counts.Services,
		"staff": counts.Staff, "documents": counts.Documents, "gallery": counts.Gallery,
	})
}
func (h handler) login(w http.ResponseWriter, r *http.Request) {
	var input loginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || strings.TrimSpace(input.Login) == "" || input.Password == "" {
		h.validationFailure(w, "login", "Укажите логин и пароль")
		return
	}
	session, err := h.store.Login(r.Context(), strings.TrimSpace(input.Login), input.Password)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: APIError{Code: "invalid_credentials", Message: "Неверный логин или пароль"}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": session.Token, "expiresAt": session.ExpiresAt, "user": map[string]any{"id": session.User.ID, "login": session.User.Login, "role": session.User.Role}})
}
func (h handler) logout(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(adminContextKey{}).(store.AdminUser)
	if err := h.store.RevokeSessions(r.Context(), user.ID); err != nil {
		h.failure(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h handler) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if raw == "" {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: APIError{Code: "unauthorized", Message: "Требуется авторизация"}})
			return
		}
		user, err := h.store.Session(r.Context(), raw)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: APIError{Code: "unauthorized", Message: "Сессия истекла или отозвана"}})
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), adminContextKey{}, user)))
	})
}

const maxUploadSize = 20 << 20

func (h handler) mediaItems(w http.ResponseWriter, r *http.Request) {
	page, perPage, err := pageParams(r)
	if err != nil {
		h.validationFailure(w, "page", err.Error())
		return
	}
	items, total, err := h.store.Media(r.Context(), strings.TrimSpace(r.URL.Query().Get("q")), perPage, (page-1)*perPage)
	if err != nil {
		h.failure(w, err)
		return
	}
	out := make([]mediaDTO, 0, len(items))
	for _, item := range items {
		url, err := h.media.URL(r.Context(), item.ObjectKey)
		if err != nil {
			h.failure(w, err)
			return
		}
		out = append(out, mediaDTO{ID: item.ID, URL: url, MIME: item.MIME, SizeBytes: item.SizeBytes, Alt: item.Alt})
	}
	writeJSON(w, http.StatusOK, mediaPageDTO{Items: out, Page: page, PerPage: perPage, Total: total})
}

func (h handler) uploadMedia(w http.ResponseWriter, r *http.Request) {
	storage, ok := h.media.(mediaStorage)
	if !ok {
		h.failure(w, errors.New("media storage does not support uploads"))
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize+1024)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		h.validationFailure(w, "file", "Размер файла не должен превышать 20 МБ")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		h.validationFailure(w, "file", "Выберите файл для загрузки")
		return
	}
	defer file.Close()
	mime := header.Header.Get("Content-Type")
	if !allowedMediaType(mime) {
		h.validationFailure(w, "file", "Поддерживаются изображения, PDF и видео")
		return
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		h.validationFailure(w, "file", "У файла должно быть расширение")
		return
	}
	key := "uploads/" + time.Now().UTC().Format("2006/01/02") + "/" + uuid.NewString() + ext
	if err := storage.Upload(r.Context(), key, mime, header.Size, file); err != nil {
		h.failure(w, err)
		return
	}
	item, err := h.store.CreateMedia(r.Context(), key, mime, header.Size, strings.TrimSpace(r.FormValue("alt")))
	if err != nil {
		_ = storage.Delete(r.Context(), key)
		h.failure(w, err)
		return
	}
	url, err := h.media.URL(r.Context(), key)
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, mediaDTO{ID: item.ID, URL: url, MIME: item.MIME, SizeBytes: item.SizeBytes, Alt: item.Alt})
}

func (h handler) deleteMedia(w http.ResponseWriter, r *http.Request) {
	storage, ok := h.media.(mediaStorage)
	if !ok {
		h.failure(w, errors.New("media storage does not support deletion"))
		return
	}
	key, err := h.store.DeleteMedia(r.Context(), pathID(r))
	if err != nil {
		h.failure(w, err)
		return
	}
	if err := storage.Delete(r.Context(), key); err != nil {
		h.failure(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h handler) adminItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.store.AdminItems(r.Context(), chi.URLParam(r, "resource"))
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h handler) setPublication(w http.ResponseWriter, r *http.Request) {
	var input publicationRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: APIError{Code: "invalid_request", Message: "Ожидается JSON с isPublished"}})
		return
	}
	if err := h.store.SetPublished(r.Context(), chi.URLParam(r, "resource"), pathID(r), input.IsPublished); err != nil {
		h.failure(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h handler) updateAdminItem(w http.ResponseWriter, r *http.Request) {
	var input contentUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || strings.TrimSpace(input.Title) == "" {
		h.validationFailure(w, "title", "Название обязательно")
		return
	}
	if input.Position < 0 {
		h.validationFailure(w, "position", "Порядок не может быть отрицательным")
		return
	}
	if err := h.store.UpdateAdminItem(r.Context(), chi.URLParam(r, "resource"), pathID(r), store.ContentUpdate{Title: strings.TrimSpace(input.Title), IsPublished: input.IsPublished, Position: input.Position}); err != nil {
		h.failure(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h handler) deleteAdminItem(w http.ResponseWriter, r *http.Request) {
	if err := h.store.DeleteAdminItem(r.Context(), chi.URLParam(r, "resource"), pathID(r)); err != nil {
		h.failure(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h handler) adminContent(w http.ResponseWriter, r *http.Request) {
	data, err := h.store.AdminContent(r.Context(), chi.URLParam(r, "resource"), pathID(r))
	if err != nil {
		h.failure(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(data)
}
func (h handler) updateAdminContent(w http.ResponseWriter, r *http.Request) {
	var data json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil || !json.Valid(data) {
		h.validationFailure(w, "content", "Ожидается корректный JSON")
		return
	}
	if err := h.store.UpdateAdminContent(r.Context(), chi.URLParam(r, "resource"), pathID(r), data); err != nil {
		h.failure(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h handler) home(w http.ResponseWriter, r *http.Request) {
	keys, err := h.store.HomePhotoKeys(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	photos, err := h.urls(r.Context(), keys)
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string][]string{"heroPhotos": photos})
}

func (h handler) oneURL(ctx context.Context, key *string) (string, error) {
	if key == nil {
		return "", nil
	}
	return h.media.URL(ctx, *key)
}
func (h handler) urls(ctx context.Context, keys []string) ([]string, error) {
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		url, err := h.media.URL(ctx, key)
		if err != nil {
			return nil, err
		}
		out = append(out, url)
	}
	return out, nil
}
func (h handler) failure(w http.ResponseWriter, err error) {
	if errors.Is(err, store.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: APIError{Code: "not_found", Message: "Ресурс не найден"}})
		return
	}
	if errors.Is(err, store.ErrMediaInUse) {
		writeJSON(w, http.StatusConflict, ErrorResponse{Error: APIError{Code: "media_in_use", Message: "Файл используется в материалах и не может быть удалён"}})
		return
	}
	writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: APIError{Code: "internal_error", Message: "Внутренняя ошибка сервера"}})
}
func (h handler) validationFailure(w http.ResponseWriter, field, message string) {
	writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: APIError{Code: "validation_error", Message: message, Field: field}})
}
func pageParams(r *http.Request) (int, int, error) {
	page, perPage := 1, 24
	var err error
	if value := r.URL.Query().Get("page"); value != "" {
		page, err = strconv.Atoi(value)
		if err != nil || page < 1 {
			return 0, 0, errors.New("Параметр page должен быть положительным числом")
		}
	}
	if value := r.URL.Query().Get("perPage"); value != "" {
		perPage, err = strconv.Atoi(value)
		if err != nil || perPage < 1 || perPage > 100 {
			return 0, 0, errors.New("Параметр perPage должен быть числом от 1 до 100")
		}
	}
	return page, perPage, nil
}
func allowedMediaType(value string) bool {
	return strings.HasPrefix(value, "image/") || strings.HasPrefix(value, "video/") || value == "application/pdf"
}
func pathID(r *http.Request) int64 {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		return 0
	}
	return id
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func fileType(mime string) string {
	if mime == "application/pdf" {
		return "PDF"
	}
	parts := strings.Split(mime, "/")
	if len(parts) == 2 {
		return strings.ToUpper(parts[1])
	}
	return strings.ToUpper(mime)
}
func fileSize(bytes int64) string {
	const mb = 1024 * 1024
	const kb = 1024
	if bytes >= mb {
		return strings.ReplaceAll(strconv.FormatFloat(float64(bytes)/mb, 'f', 1, 64), ".", ",") + " МБ"
	}
	return strconv.FormatFloat(float64(bytes)/kb, 'f', 0, 64) + " КБ"
}
