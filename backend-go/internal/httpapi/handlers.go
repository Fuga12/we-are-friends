package httpapi

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	mediafiles "github.com/Fuga12/we-are-friends/backend-go/internal/media"
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
	store         *store.Store
	media         mediaURLs
	adminUser     string
	adminPassword string
}
type adminContextKey struct{}
type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type leadRequest struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Message string `json:"message"`
}
type leadStatusRequest struct {
	IsProcessed bool `json:"isProcessed"`
}

type publicationRequest struct {
	IsPublished bool `json:"isPublished"`
}
type contentUpdateRequest struct {
	Title       string `json:"title"`
	IsPublished bool   `json:"isPublished"`
	Position    int    `json:"position"`
}
type adminClassDTO struct {
	ID              int64      `json:"id,omitempty"`
	Title           string     `json:"title"`
	AgeRange        string     `json:"ageRange"`
	Description     string     `json:"description"`
	Schedule        string     `json:"schedule"`
	Price           string     `json:"price"`
	Emoji           string     `json:"emoji"`
	Color           string     `json:"color"`
	CoverMediaID    *int64     `json:"coverMediaId"`
	PhotoMediaIDs   []int64    `json:"photoMediaIds"`
	FullDescription []string   `json:"fullDescription"`
	Highlights      []string   `json:"highlights"`
	Position        int        `json:"position"`
	IsPublished     bool       `json:"isPublished"`
	CoverMedia      *mediaDTO  `json:"coverMedia,omitempty"`
	PhotoMedia      []mediaDTO `json:"photoMedia"`
}
type adminNewsDTO struct {
	ID            int64      `json:"id,omitempty"`
	Title         string     `json:"title"`
	Date          string     `json:"date"`
	Category      string     `json:"category"`
	Excerpt       string     `json:"excerpt"`
	Emoji         string     `json:"emoji"`
	CoverMediaID  *int64     `json:"coverMediaId"`
	PhotoMediaIDs []int64    `json:"photoMediaIds"`
	Content       []string   `json:"content"`
	IsPublished   bool       `json:"isPublished"`
	CoverMedia    *mediaDTO  `json:"coverMedia,omitempty"`
	PhotoMedia    []mediaDTO `json:"photoMedia"`
}
type adminServiceDTO struct {
	ID            int64      `json:"id,omitempty"`
	Title         string     `json:"title"`
	Tagline       string     `json:"tagline"`
	Description   string     `json:"description"`
	Features      []string   `json:"features"`
	Price         string     `json:"price"`
	Emoji         string     `json:"emoji"`
	Path          string     `json:"to"`
	Color         string     `json:"color"`
	CoverMediaID  *int64     `json:"coverMediaId"`
	PhotoMediaIDs []int64    `json:"photoMediaIds"`
	Position      int        `json:"position"`
	IsPublished   bool       `json:"isPublished"`
	CoverMedia    *mediaDTO  `json:"coverMedia,omitempty"`
	PhotoMedia    []mediaDTO `json:"photoMedia"`
}
type adminStaffDTO struct {
	ID           int64     `json:"id,omitempty"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	Bio          string    `json:"bio"`
	Initials     string    `json:"initials"`
	Color        string    `json:"color"`
	PhotoMediaID *int64    `json:"photoMediaId"`
	Position     int       `json:"position"`
	IsPublished  bool      `json:"isPublished"`
	PhotoMedia   *mediaDTO `json:"photoMedia,omitempty"`
}
type adminDocumentDTO struct {
	ID              int64     `json:"id,omitempty"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	FileMediaID     int64     `json:"fileMediaId"`
	DisplayFileType string    `json:"displayFileType"`
	DisplayFileSize string    `json:"displayFileSize"`
	Position        int       `json:"position"`
	IsPublished     bool      `json:"isPublished"`
	FileMedia       *mediaDTO `json:"fileMedia,omitempty"`
}
type adminReviewDTO struct {
	ID          int64  `json:"id,omitempty"`
	Author      string `json:"author"`
	Role        string `json:"role"`
	Text        string `json:"text"`
	Rating      int16  `json:"rating"`
	Position    int    `json:"position"`
	IsPublished bool   `json:"isPublished"`
}
type adminGalleryDTO struct {
	ID       int64    `json:"id"`
	MediaID  int64    `json:"mediaId"`
	Title    string   `json:"title"`
	Position int      `json:"position"`
	Media    mediaDTO `json:"media"`
}
type addGalleryRequest struct {
	MediaID int64  `json:"mediaId"`
	Title   string `json:"title"`
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
func (h handler) summerShifts(w http.ResponseWriter, r *http.Request) {
	values, err := h.store.SummerShifts(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	out := make([]classDTO, 0, len(values))
	for _, v := range values {
		dto, dtoErr := h.classDTO(r.Context(), v)
		if dtoErr != nil {
			h.failure(w, dtoErr)
			return
		}
		out = append(out, dto)
	}
	writeJSON(w, http.StatusOK, out)
}
func (h handler) summerShift(w http.ResponseWriter, r *http.Request) {
	v, err := h.store.SummerShift(r.Context(), pathID(r))
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

func (h handler) createLead(w http.ResponseWriter, r *http.Request) {
	var input leadRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.validationFailure(w, "form", "Не удалось прочитать заявку")
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Phone = strings.TrimSpace(input.Phone)
	input.Message = strings.TrimSpace(input.Message)
	if input.Name == "" || len(input.Phone) < 5 {
		h.validationFailure(w, "phone", "Укажите имя и корректный телефон")
		return
	}
	if len(input.Name) > 150 || len(input.Phone) > 50 || len(input.Message) > 3000 {
		h.validationFailure(w, "form", "Заявка содержит слишком длинные поля")
		return
	}
	id, err := h.store.CreateLead(r.Context(), input.Name, input.Phone, input.Message)
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (h handler) adminLeads(w http.ResponseWriter, r *http.Request) {
	items, err := h.store.Leads(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (h handler) updateAdminLead(w http.ResponseWriter, r *http.Request) {
	var input leadStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.validationFailure(w, "isProcessed", "Не удалось прочитать статус заявки")
		return
	}
	if err := h.store.SetLeadProcessed(r.Context(), pathID(r), input.IsProcessed); err != nil {
		h.failure(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h handler) deleteAdminLead(w http.ResponseWriter, r *http.Request) {
	if err := h.store.DeleteLead(r.Context(), pathID(r)); err != nil {
		h.failure(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h handler) adminOverview(w http.ResponseWriter, r *http.Request) {
	counts, err := h.store.AdminCounts(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{
		"classes": counts.Classes, "summer-shifts": counts.SummerShifts, "news": counts.News, "services": counts.Services,
		"staff": counts.Staff, "documents": counts.Documents, "gallery": counts.Gallery,
		"reviews": counts.Reviews,
		"leads":   counts.Leads,
	})
}
func (h handler) login(w http.ResponseWriter, r *http.Request) {
	var input loginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || strings.TrimSpace(input.Login) == "" || input.Password == "" {
		h.validationFailure(w, "login", "Укажите логин и пароль")
		return
	}
	if !validAdminCredentials(strings.TrimSpace(input.Login), input.Password, h.adminUser, h.adminPassword) {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: APIError{Code: "invalid_credentials", Message: "Неверный логин или пароль"}})
		return
	}
	session, err := h.store.CreateSession(r.Context(), h.adminUser)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{Error: APIError{Code: "invalid_credentials", Message: "Неверный логин или пароль"}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": session.Token, "expiresAt": session.ExpiresAt, "user": map[string]any{"id": session.User.ID, "login": session.User.Login, "role": session.User.Role}})
}

func validAdminCredentials(login, password, expectedLogin, expectedPassword string) bool {
	loginHash := sha256.Sum256([]byte(login))
	expectedLoginHash := sha256.Sum256([]byte(expectedLogin))
	passwordHash := sha256.Sum256([]byte(password))
	expectedPasswordHash := sha256.Sum256([]byte(expectedPassword))
	return subtle.ConstantTimeCompare(loginHash[:], expectedLoginHash[:]) == 1 &&
		subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1
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
	uploadBody := io.Reader(file)
	uploadSize := header.Size
	if strings.HasPrefix(mime, "image/") {
		optimized, err := mediafiles.OptimizeImage(file)
		if err != nil {
			h.validationFailure(w, "file", "Не удалось обработать изображение. Загрузите JPEG, PNG, GIF или WebP")
			return
		}
		mime = "image/jpeg"
		ext = ".jpg"
		uploadSize = int64(len(optimized))
		uploadBody = bytes.NewReader(optimized)
	}
	key := "uploads/" + time.Now().UTC().Format("2006/01/02") + "/" + uuid.NewString() + ext
	if err := storage.Upload(r.Context(), key, mime, uploadSize, uploadBody); err != nil {
		h.failure(w, err)
		return
	}
	item, err := h.store.CreateMedia(r.Context(), key, mime, uploadSize, strings.TrimSpace(r.FormValue("alt")))
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

func (h handler) adminClassForm(w http.ResponseWriter, r *http.Request) {
	item, err := h.store.AdminClass(r.Context(), pathID(r))
	if err != nil {
		h.failure(w, err)
		return
	}
	ids := append([]int64{}, item.PhotoMediaIDs...)
	if item.CoverMediaID != nil {
		ids = append([]int64{*item.CoverMediaID}, ids...)
	}
	mediaItems, err := h.store.MediaByIDs(r.Context(), ids)
	if err != nil {
		h.failure(w, err)
		return
	}
	mediaByID := make(map[int64]mediaDTO, len(mediaItems))
	for _, mediaItem := range mediaItems {
		url, urlErr := h.media.URL(r.Context(), mediaItem.ObjectKey)
		if urlErr != nil {
			h.failure(w, urlErr)
			return
		}
		mediaByID[mediaItem.ID] = mediaDTO{ID: mediaItem.ID, URL: url, MIME: mediaItem.MIME, SizeBytes: mediaItem.SizeBytes, Alt: mediaItem.Alt}
	}
	output := adminClassDTO{ID: item.ID, Title: item.Title, AgeRange: item.AgeRange, Description: item.Description, Schedule: item.Schedule,
		Price: item.Price, Emoji: item.Emoji, Color: item.Color, CoverMediaID: item.CoverMediaID, PhotoMediaIDs: item.PhotoMediaIDs,
		FullDescription: item.FullDescription, Highlights: item.Highlights, Position: item.Position, IsPublished: item.IsPublished, PhotoMedia: []mediaDTO{}}
	if item.CoverMediaID != nil {
		if cover, ok := mediaByID[*item.CoverMediaID]; ok {
			output.CoverMedia = &cover
		}
	}
	for _, id := range item.PhotoMediaIDs {
		if photo, ok := mediaByID[id]; ok {
			output.PhotoMedia = append(output.PhotoMedia, photo)
		}
	}
	writeJSON(w, http.StatusOK, output)
}

func (h handler) adminSummerShiftForm(w http.ResponseWriter, r *http.Request) {
	item, err := h.store.AdminSummerShift(r.Context(), pathID(r))
	if err != nil {
		h.failure(w, err)
		return
	}
	ids := append([]int64{}, item.PhotoMediaIDs...)
	if item.CoverMediaID != nil {
		ids = append([]int64{*item.CoverMediaID}, ids...)
	}
	mediaItems, err := h.store.MediaByIDs(r.Context(), ids)
	if err != nil {
		h.failure(w, err)
		return
	}
	mediaByID := make(map[int64]mediaDTO, len(mediaItems))
	for _, mediaItem := range mediaItems {
		url, urlErr := h.media.URL(r.Context(), mediaItem.ObjectKey)
		if urlErr != nil {
			h.failure(w, urlErr)
			return
		}
		mediaByID[mediaItem.ID] = mediaDTO{ID: mediaItem.ID, URL: url, MIME: mediaItem.MIME, SizeBytes: mediaItem.SizeBytes, Alt: mediaItem.Alt}
	}
	output := adminClassDTO{ID: item.ID, Title: item.Title, AgeRange: item.AgeRange, Description: item.Description, Schedule: item.Schedule, Price: item.Price, Emoji: item.Emoji, Color: item.Color, CoverMediaID: item.CoverMediaID, PhotoMediaIDs: item.PhotoMediaIDs, FullDescription: item.FullDescription, Highlights: item.Highlights, Position: item.Position, IsPublished: item.IsPublished, PhotoMedia: []mediaDTO{}}
	if item.CoverMediaID != nil {
		if cover, ok := mediaByID[*item.CoverMediaID]; ok {
			output.CoverMedia = &cover
		}
	}
	for _, id := range item.PhotoMediaIDs {
		if photo, ok := mediaByID[id]; ok {
			output.PhotoMedia = append(output.PhotoMedia, photo)
		}
	}
	writeJSON(w, http.StatusOK, output)
}

func (h handler) deleteAdminClassMedia(w http.ResponseWriter, r *http.Request) {
	storage, ok := h.media.(mediaStorage)
	if !ok {
		h.failure(w, errors.New("media storage does not support deletion"))
		return
	}
	classID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || classID < 1 {
		h.failure(w, store.ErrNotFound)
		return
	}
	mediaID, err := strconv.ParseInt(chi.URLParam(r, "mediaId"), 10, 64)
	if err != nil || mediaID < 1 {
		h.failure(w, store.ErrNotFound)
		return
	}
	key, err := h.store.DeleteClassMedia(r.Context(), classID, mediaID)
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
func (h handler) deleteAdminSummerShiftMedia(w http.ResponseWriter, r *http.Request) {
	storage, ok := h.media.(mediaStorage)
	if !ok {
		h.failure(w, errors.New("media storage does not support deletion"))
		return
	}
	shiftID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || shiftID < 1 {
		h.failure(w, store.ErrNotFound)
		return
	}
	mediaID, err := strconv.ParseInt(chi.URLParam(r, "mediaId"), 10, 64)
	if err != nil || mediaID < 1 {
		h.failure(w, store.ErrNotFound)
		return
	}
	key, err := h.store.DeleteSummerShiftMedia(r.Context(), shiftID, mediaID)
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

func (h handler) adminNewsForm(w http.ResponseWriter, r *http.Request) {
	item, err := h.store.AdminNews(r.Context(), pathID(r))
	if err != nil {
		h.failure(w, err)
		return
	}
	ids := append([]int64{}, item.PhotoMediaIDs...)
	if item.CoverMediaID != nil {
		ids = append([]int64{*item.CoverMediaID}, ids...)
	}
	mediaItems, err := h.store.MediaByIDs(r.Context(), ids)
	if err != nil {
		h.failure(w, err)
		return
	}
	mediaByID := make(map[int64]mediaDTO, len(mediaItems))
	for _, mediaItem := range mediaItems {
		url, urlErr := h.media.URL(r.Context(), mediaItem.ObjectKey)
		if urlErr != nil {
			h.failure(w, urlErr)
			return
		}
		mediaByID[mediaItem.ID] = mediaDTO{ID: mediaItem.ID, URL: url, MIME: mediaItem.MIME, SizeBytes: mediaItem.SizeBytes, Alt: mediaItem.Alt}
	}
	output := adminNewsDTO{ID: item.ID, Title: item.Title, Date: item.Date, Category: item.Category, Excerpt: item.Excerpt, Emoji: item.Emoji, CoverMediaID: item.CoverMediaID, PhotoMediaIDs: item.PhotoMediaIDs, Content: item.Content, IsPublished: item.IsPublished, PhotoMedia: []mediaDTO{}}
	if item.CoverMediaID != nil {
		if cover, ok := mediaByID[*item.CoverMediaID]; ok {
			output.CoverMedia = &cover
		}
	}
	for _, id := range item.PhotoMediaIDs {
		if photo, ok := mediaByID[id]; ok {
			output.PhotoMedia = append(output.PhotoMedia, photo)
		}
	}
	writeJSON(w, http.StatusOK, output)
}

func (h handler) adminNewsCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.store.NewsCategories(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, categories)
}

func (h handler) saveAdminNews(w http.ResponseWriter, r *http.Request) {
	var input adminNewsDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.validationFailure(w, "form", "Не удалось прочитать данные формы")
		return
	}
	if chi.URLParam(r, "id") != "" {
		input.ID = pathID(r)
	}
	input.Title = strings.TrimSpace(input.Title)
	input.Category = strings.TrimSpace(input.Category)
	input.Excerpt = strings.TrimSpace(input.Excerpt)
	if input.Title == "" {
		h.validationFailure(w, "title", "Укажите заголовок новости")
		return
	}
	if _, err := time.Parse("2006-01-02", input.Date); err != nil {
		h.validationFailure(w, "date", "Укажите дату публикации")
		return
	}
	if input.Category == "" {
		h.validationFailure(w, "category", "Укажите рубрику")
		return
	}
	if input.Excerpt == "" {
		h.validationFailure(w, "excerpt", "Добавьте краткий анонс")
		return
	}
	id, err := h.store.SaveAdminNews(r.Context(), store.AdminNews{ID: input.ID, Title: input.Title, Date: input.Date, Category: input.Category, Excerpt: input.Excerpt, Emoji: strings.TrimSpace(input.Emoji), CoverMediaID: input.CoverMediaID, PhotoMediaIDs: input.PhotoMediaIDs, Content: input.Content, IsPublished: input.IsPublished})
	if err != nil {
		h.failure(w, err)
		return
	}
	status := http.StatusOK
	if input.ID == 0 {
		status = http.StatusCreated
	}
	writeJSON(w, status, map[string]int64{"id": id})
}

func (h handler) deleteAdminNewsMedia(w http.ResponseWriter, r *http.Request) {
	storage, ok := h.media.(mediaStorage)
	if !ok {
		h.failure(w, errors.New("media storage does not support deletion"))
		return
	}
	newsID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || newsID < 1 {
		h.failure(w, store.ErrNotFound)
		return
	}
	mediaID, err := strconv.ParseInt(chi.URLParam(r, "mediaId"), 10, 64)
	if err != nil || mediaID < 1 {
		h.failure(w, store.ErrNotFound)
		return
	}
	key, err := h.store.DeleteNewsMedia(r.Context(), newsID, mediaID)
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

func (h handler) adminServiceForm(w http.ResponseWriter, r *http.Request) {
	item, err := h.store.AdminService(r.Context(), pathID(r))
	if err != nil {
		h.failure(w, err)
		return
	}
	ids := append([]int64{}, item.PhotoMediaIDs...)
	if item.CoverMediaID != nil {
		ids = append([]int64{*item.CoverMediaID}, ids...)
	}
	mediaItems, err := h.store.MediaByIDs(r.Context(), ids)
	if err != nil {
		h.failure(w, err)
		return
	}
	byID := map[int64]mediaDTO{}
	for _, m := range mediaItems {
		url, urlErr := h.media.URL(r.Context(), m.ObjectKey)
		if urlErr != nil {
			h.failure(w, urlErr)
			return
		}
		byID[m.ID] = mediaDTO{ID: m.ID, URL: url, MIME: m.MIME, SizeBytes: m.SizeBytes, Alt: m.Alt}
	}
	out := adminServiceDTO{ID: item.ID, Title: item.Title, Tagline: item.Tagline, Description: item.Description, Features: item.Features, Price: item.Price, Emoji: item.Emoji, Path: item.Path, Color: item.Color, CoverMediaID: item.CoverMediaID, PhotoMediaIDs: item.PhotoMediaIDs, Position: item.Position, IsPublished: item.IsPublished, PhotoMedia: []mediaDTO{}}
	if item.CoverMediaID != nil {
		if cover, ok := byID[*item.CoverMediaID]; ok {
			out.CoverMedia = &cover
		}
	}
	for _, id := range item.PhotoMediaIDs {
		if photo, ok := byID[id]; ok {
			out.PhotoMedia = append(out.PhotoMedia, photo)
		}
	}
	writeJSON(w, http.StatusOK, out)
}
func (h handler) saveAdminService(w http.ResponseWriter, r *http.Request) {
	var input adminServiceDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.validationFailure(w, "form", "Не удалось прочитать данные формы")
		return
	}
	if chi.URLParam(r, "id") != "" {
		input.ID = pathID(r)
	}
	input.Title = strings.TrimSpace(input.Title)
	input.Tagline = strings.TrimSpace(input.Tagline)
	input.Description = strings.TrimSpace(input.Description)
	input.Path = strings.TrimSpace(input.Path)
	if input.Title == "" || input.Tagline == "" || input.Description == "" {
		h.validationFailure(w, "form", "Заполните название, подзаголовок и описание")
		return
	}
	if !strings.HasPrefix(input.Path, "/special/") {
		h.validationFailure(w, "to", "Адрес должен начинаться с /special/")
		return
	}
	if len(input.Color) != 7 || input.Color[0] != '#' {
		h.validationFailure(w, "color", "Укажите цвет в формате #1a6757")
		return
	}
	id, err := h.store.SaveAdminService(r.Context(), store.AdminService{ID: input.ID, Title: input.Title, Tagline: input.Tagline, Description: input.Description, Features: input.Features, Price: strings.TrimSpace(input.Price), Emoji: strings.TrimSpace(input.Emoji), Path: input.Path, Color: input.Color, CoverMediaID: input.CoverMediaID, PhotoMediaIDs: input.PhotoMediaIDs, Position: input.Position, IsPublished: input.IsPublished})
	if err != nil {
		h.failure(w, err)
		return
	}
	status := http.StatusOK
	if input.ID == 0 {
		status = http.StatusCreated
	}
	writeJSON(w, status, map[string]int64{"id": id})
}
func (h handler) deleteAdminServiceMedia(w http.ResponseWriter, r *http.Request) {
	storage, ok := h.media.(mediaStorage)
	if !ok {
		h.failure(w, errors.New("media storage does not support deletion"))
		return
	}
	serviceID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || serviceID < 1 {
		h.failure(w, store.ErrNotFound)
		return
	}
	mediaID, err := strconv.ParseInt(chi.URLParam(r, "mediaId"), 10, 64)
	if err != nil || mediaID < 1 {
		h.failure(w, store.ErrNotFound)
		return
	}
	key, err := h.store.DeleteServiceMedia(r.Context(), serviceID, mediaID)
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

func (h handler) saveAdminClass(w http.ResponseWriter, r *http.Request) {
	var input adminClassDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.validationFailure(w, "form", "Не удалось прочитать данные формы")
		return
	}
	if id := chi.URLParam(r, "id"); id != "" {
		input.ID = pathID(r)
	}
	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		h.validationFailure(w, "title", "Укажите название занятия")
		return
	}
	if input.Position < 0 {
		h.validationFailure(w, "position", "Порядок не может быть отрицательным")
		return
	}
	if strings.TrimSpace(input.Color) == "" {
		input.Color = "#e46129"
	} else if len(input.Color) != 7 || input.Color[0] != '#' {
		h.validationFailure(w, "color", "Укажите цвет в формате #1a6757")
		return
	}
	if strings.TrimSpace(input.Emoji) == "" {
		input.Emoji = "🎨"
	}
	if input.FullDescription == nil {
		input.FullDescription = []string{}
	}
	if input.Highlights == nil {
		input.Highlights = []string{}
	}
	item := store.AdminClass{ID: input.ID, Title: input.Title, AgeRange: strings.TrimSpace(input.AgeRange), Description: strings.TrimSpace(input.Description),
		Schedule: strings.TrimSpace(input.Schedule), Price: strings.TrimSpace(input.Price), Emoji: strings.TrimSpace(input.Emoji), Color: input.Color,
		CoverMediaID: input.CoverMediaID, PhotoMediaIDs: input.PhotoMediaIDs, FullDescription: input.FullDescription, Highlights: input.Highlights,
		Position: input.Position, IsPublished: input.IsPublished}
	id, err := h.store.SaveAdminClass(r.Context(), item)
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, map[bool]int{true: http.StatusCreated, false: http.StatusOK}[input.ID == 0], map[string]int64{"id": id})
}

func (h handler) saveAdminSummerShift(w http.ResponseWriter, r *http.Request) {
	var input adminClassDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.validationFailure(w, "form", "Не удалось прочитать данные формы")
		return
	}
	if id := chi.URLParam(r, "id"); id != "" {
		input.ID = pathID(r)
	}
	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		h.validationFailure(w, "title", "Укажите название смены")
		return
	}
	if input.Position < 0 {
		h.validationFailure(w, "position", "Порядок не может быть отрицательным")
		return
	}
	if strings.TrimSpace(input.Color) == "" {
		input.Color = "#e46129"
	} else if len(input.Color) != 7 || input.Color[0] != '#' {
		h.validationFailure(w, "color", "Укажите цвет в формате #1a6757")
		return
	}
	if strings.TrimSpace(input.Emoji) == "" {
		input.Emoji = "☀️"
	}
	if input.FullDescription == nil {
		input.FullDescription = []string{}
	}
	if input.Highlights == nil {
		input.Highlights = []string{}
	}
	item := store.AdminSummerShift{ID: input.ID, Title: input.Title, AgeRange: strings.TrimSpace(input.AgeRange), Description: strings.TrimSpace(input.Description), Schedule: strings.TrimSpace(input.Schedule), Price: strings.TrimSpace(input.Price), Emoji: strings.TrimSpace(input.Emoji), Color: input.Color, CoverMediaID: input.CoverMediaID, PhotoMediaIDs: input.PhotoMediaIDs, FullDescription: input.FullDescription, Highlights: input.Highlights, Position: input.Position, IsPublished: input.IsPublished}
	id, err := h.store.SaveAdminSummerShift(r.Context(), item)
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, map[bool]int{true: http.StatusCreated, false: http.StatusOK}[input.ID == 0], map[string]int64{"id": id})
}

func (h handler) mediaDTOByID(ctx context.Context, id int64) (*mediaDTO, error) {
	items, err := h.store.MediaByIDs(ctx, []int64{id})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, store.ErrNotFound
	}
	url, err := h.media.URL(ctx, items[0].ObjectKey)
	if err != nil {
		return nil, err
	}
	return &mediaDTO{ID: items[0].ID, URL: url, MIME: items[0].MIME, SizeBytes: items[0].SizeBytes, Alt: items[0].Alt}, nil
}
func (h handler) adminStaffForm(w http.ResponseWriter, r *http.Request) {
	item, err := h.store.AdminStaff(r.Context(), pathID(r))
	if err != nil {
		h.failure(w, err)
		return
	}
	out := adminStaffDTO{ID: item.ID, Name: item.Name, Role: item.Role, Bio: item.Bio, Initials: item.Initials, Color: item.Color, PhotoMediaID: item.PhotoMediaID, Position: item.Position, IsPublished: item.IsPublished}
	if item.PhotoMediaID != nil {
		out.PhotoMedia, err = h.mediaDTOByID(r.Context(), *item.PhotoMediaID)
		if err != nil {
			h.failure(w, err)
			return
		}
	}
	writeJSON(w, http.StatusOK, out)
}
func (h handler) saveAdminStaff(w http.ResponseWriter, r *http.Request) {
	var input adminStaffDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.validationFailure(w, "form", "Не удалось прочитать данные формы")
		return
	}
	if chi.URLParam(r, "id") != "" {
		input.ID = pathID(r)
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Role = strings.TrimSpace(input.Role)
	input.Bio = strings.TrimSpace(input.Bio)
	input.Initials = strings.TrimSpace(input.Initials)
	if input.Name == "" || input.Role == "" {
		h.validationFailure(w, "name", "Укажите имя и должность")
		return
	}
	if len(input.Color) != 7 || input.Color[0] != '#' {
		h.validationFailure(w, "color", "Укажите цвет в формате #1a6757")
		return
	}
	id, err := h.store.SaveAdminStaff(r.Context(), store.AdminStaff{ID: input.ID, Name: input.Name, Role: input.Role, Bio: input.Bio, Initials: input.Initials, Color: input.Color, PhotoMediaID: input.PhotoMediaID, Position: input.Position, IsPublished: input.IsPublished})
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, map[bool]int{true: http.StatusCreated, false: http.StatusOK}[input.ID == 0], map[string]int64{"id": id})
}
func (h handler) deleteAdminStaffMedia(w http.ResponseWriter, r *http.Request) {
	storage, ok := h.media.(mediaStorage)
	if !ok {
		h.failure(w, errors.New("media storage does not support deletion"))
		return
	}
	staffID := pathID(r)
	mediaID, err := strconv.ParseInt(chi.URLParam(r, "mediaId"), 10, 64)
	if err != nil {
		h.failure(w, store.ErrNotFound)
		return
	}
	key, err := h.store.DeleteStaffMedia(r.Context(), staffID, mediaID)
	if err != nil {
		h.failure(w, err)
		return
	}
	if err = storage.Delete(r.Context(), key); err != nil {
		h.failure(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h handler) adminDocumentForm(w http.ResponseWriter, r *http.Request) {
	item, err := h.store.AdminDocument(r.Context(), pathID(r))
	if err != nil {
		h.failure(w, err)
		return
	}
	file, err := h.mediaDTOByID(r.Context(), item.FileMediaID)
	if err != nil {
		h.failure(w, err)
		return
	}
	out := adminDocumentDTO{ID: item.ID, Title: item.Title, Description: item.Description, FileMediaID: item.FileMediaID, Position: item.Position, IsPublished: item.IsPublished, FileMedia: file}
	if item.DisplayFileType != nil {
		out.DisplayFileType = *item.DisplayFileType
	}
	if item.DisplayFileSize != nil {
		out.DisplayFileSize = *item.DisplayFileSize
	}
	writeJSON(w, http.StatusOK, out)
}
func (h handler) saveAdminDocument(w http.ResponseWriter, r *http.Request) {
	var input adminDocumentDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.validationFailure(w, "form", "Не удалось прочитать данные формы")
		return
	}
	if chi.URLParam(r, "id") != "" {
		input.ID = pathID(r)
	}
	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" || input.FileMediaID < 1 {
		h.validationFailure(w, "file", "Укажите название и загрузите файл")
		return
	}
	var fileType, fileSize *string
	if value := strings.TrimSpace(input.DisplayFileType); value != "" {
		fileType = &value
	}
	if value := strings.TrimSpace(input.DisplayFileSize); value != "" {
		fileSize = &value
	}
	id, err := h.store.SaveAdminDocument(r.Context(), store.AdminDocument{ID: input.ID, Title: input.Title, Description: strings.TrimSpace(input.Description), FileMediaID: input.FileMediaID, DisplayFileType: fileType, DisplayFileSize: fileSize, Position: input.Position, IsPublished: input.IsPublished})
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, map[bool]int{true: http.StatusCreated, false: http.StatusOK}[input.ID == 0], map[string]int64{"id": id})
}
func (h handler) adminReviewForm(w http.ResponseWriter, r *http.Request) {
	item, err := h.store.AdminReview(r.Context(), pathID(r))
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, http.StatusOK, adminReviewDTO{ID: item.ID, Author: item.Author, Role: item.Role, Text: item.Text, Rating: item.Rating, Position: item.Position, IsPublished: item.IsPublished})
}
func (h handler) saveAdminReview(w http.ResponseWriter, r *http.Request) {
	var input adminReviewDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.validationFailure(w, "form", "Не удалось прочитать данные формы")
		return
	}
	if chi.URLParam(r, "id") != "" {
		input.ID = pathID(r)
	}
	input.Author = strings.TrimSpace(input.Author)
	input.Text = strings.TrimSpace(input.Text)
	if input.Author == "" || input.Text == "" {
		h.validationFailure(w, "author", "Укажите автора и текст отзыва")
		return
	}
	if input.Rating < 1 || input.Rating > 5 {
		h.validationFailure(w, "rating", "Оценка должна быть от 1 до 5")
		return
	}
	id, err := h.store.SaveAdminReview(r.Context(), store.AdminReview{ID: input.ID, Author: input.Author, Role: strings.TrimSpace(input.Role), Text: input.Text, Rating: input.Rating, Position: input.Position, IsPublished: input.IsPublished})
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, map[bool]int{true: http.StatusCreated, false: http.StatusOK}[input.ID == 0], map[string]int64{"id": id})
}
func (h handler) adminGallery(w http.ResponseWriter, r *http.Request) {
	items, err := h.store.AdminGallery(r.Context())
	if err != nil {
		h.failure(w, err)
		return
	}
	out := make([]adminGalleryDTO, 0, len(items))
	for _, item := range items {
		url, urlErr := h.media.URL(r.Context(), item.Media.ObjectKey)
		if urlErr != nil {
			h.failure(w, urlErr)
			return
		}
		out = append(out, adminGalleryDTO{ID: item.ID, MediaID: item.MediaID, Title: item.Title, Position: item.Position, Media: mediaDTO{ID: item.MediaID, URL: url, MIME: item.Media.MIME, SizeBytes: item.Media.SizeBytes, Alt: item.Media.Alt}})
	}
	writeJSON(w, http.StatusOK, out)
}
func (h handler) addAdminGallery(w http.ResponseWriter, r *http.Request) {
	var input addGalleryRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.MediaID < 1 {
		h.validationFailure(w, "mediaId", "Сначала загрузите фотографию")
		return
	}
	id, err := h.store.AddGalleryItem(r.Context(), input.MediaID, strings.TrimSpace(input.Title))
	if err != nil {
		h.failure(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
}
func (h handler) deleteAdminGallery(w http.ResponseWriter, r *http.Request) {
	storage, ok := h.media.(mediaStorage)
	if !ok {
		h.failure(w, errors.New("media storage does not support deletion"))
		return
	}
	key, err := h.store.DeleteGalleryItem(r.Context(), pathID(r))
	if err != nil {
		h.failure(w, err)
		return
	}
	if err = storage.Delete(r.Context(), key); err != nil {
		h.failure(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
	resource := chi.URLParam(r, "resource")
	if resource == "staff" || resource == "documents" {
		storage, ok := h.media.(mediaStorage)
		if !ok {
			h.failure(w, errors.New("media storage does not support deletion"))
			return
		}
		key, err := h.store.DeleteAdminAsset(r.Context(), resource, pathID(r))
		if err != nil {
			h.failure(w, err)
			return
		}
		if key != "" {
			if err = storage.Delete(r.Context(), key); err != nil {
				h.failure(w, err)
				return
			}
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err := h.store.DeleteAdminItem(r.Context(), resource, pathID(r)); err != nil {
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
