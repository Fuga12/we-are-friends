// Command seed imports the current Litestar stub into PostgreSQL and MinIO.
// Start the stub on port 3001 first, then run `go run ./seed` from backend-go.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Fuga12/we-are-friends/backend-go/internal/config"
	"github.com/Fuga12/we-are-friends/backend-go/internal/media"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type home struct {
	HeroPhotos []string `json:"heroPhotos"`
}
type class struct {
	ID              int64    `json:"id"`
	Title           string   `json:"title"`
	AgeRange        string   `json:"ageRange"`
	Description     string   `json:"description"`
	Schedule        string   `json:"schedule"`
	Price           string   `json:"price"`
	Emoji           string   `json:"emoji"`
	Color           string   `json:"color"`
	Photo           string   `json:"photo"`
	Photos          []string `json:"photos"`
	FullDescription []string `json:"fullDescription"`
	Highlights      []string `json:"highlights"`
}
type news struct {
	ID       int64    `json:"id"`
	Title    string   `json:"title"`
	Date     string   `json:"date"`
	Category string   `json:"category"`
	Excerpt  string   `json:"excerpt"`
	Emoji    string   `json:"emoji"`
	Photo    string   `json:"photo"`
	Photos   []string `json:"photos"`
	Content  []string `json:"content"`
}
type service struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Tagline     string          `json:"tagline"`
	Description string          `json:"description"`
	Features    []string        `json:"features"`
	Price       string          `json:"price"`
	Emoji       string          `json:"emoji"`
	To          string          `json:"to"`
	Color       string          `json:"color"`
	Photo       string          `json:"photo"`
	Photos      []string        `json:"photos"`
	Detail      json.RawMessage `json:"detail"`
}
type staff struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Bio      string `json:"bio"`
	Initials string `json:"initials"`
	Color    string `json:"color"`
	Photo    string `json:"photo"`
}
type document struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	FileType    string `json:"fileType"`
	FileSize    string `json:"fileSize"`
}
type gallery struct {
	ID    int64  `json:"id"`
	Image string `json:"image"`
}

func main() {
	cfg, err := config.Load()
	must(err)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	must(err)
	defer pool.Close()
	minio, err := media.New(cfg)
	must(err)
	must(minio.EnsureBucket(ctx))
	stubURL := strings.TrimRight(valueOr("STUB_URL", "http://localhost:3001"), "/")
	stubMedia := valueOr("STUB_MEDIA_DIR", "../backend-stub/media")
	client := &http.Client{Timeout: 15 * time.Second}
	var h home
	get(ctx, client, stubURL+"/api/home", &h)
	var classes []class
	get(ctx, client, stubURL+"/api/classes", &classes)
	var articles []news
	get(ctx, client, stubURL+"/api/news", &articles)
	var services []service
	get(ctx, client, stubURL+"/api/services", &services)
	var people []staff
	get(ctx, client, stubURL+"/api/staff", &people)
	var documents []document
	get(ctx, client, stubURL+"/api/documents", &documents)
	var pictures []gallery
	get(ctx, client, stubURL+"/api/gallery", &pictures)

	tx, err := pool.Begin(ctx)
	must(err)
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `truncate home_photos, site_settings, reviews, gallery, documents, staff, service_photos, services, news_photos, news, class_photos, classes, media restart identity cascade`)
	must(err)
	mediaIDs := map[string]int64{}
	mediaID := func(raw string) int64 {
		key := mediaKey(raw)
		if id, ok := mediaIDs[key]; ok {
			return id
		}
		path := filepath.Join(stubMedia, key)
		file, err := os.Open(path)
		must(err)
		defer file.Close()
		info, err := file.Stat()
		must(err)
		must(minio.Upload(ctx, key, "image/jpeg", info.Size(), file))
		var id int64
		err = tx.QueryRow(ctx, `insert into media (object_key,mime,size_bytes) values ($1,$2,$3) returning id`, key, "image/jpeg", info.Size()).Scan(&id)
		must(err)
		mediaIDs[key] = id
		return id
	}

	for position, v := range classes {
		cover := mediaID(v.Photo)
		_, err = tx.Exec(ctx, `insert into classes (id,title,age_range,description,schedule,price,emoji,color,cover_media_id,full_description,highlights,position) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`, v.ID, v.Title, v.AgeRange, v.Description, v.Schedule, v.Price, v.Emoji, v.Color, cover, v.FullDescription, v.Highlights, position)
		must(err)
		for p, photo := range v.Photos {
			_, err = tx.Exec(ctx, `insert into class_photos (class_id,media_id,position) values ($1,$2,$3)`, v.ID, mediaID(photo), p)
			must(err)
		}
	}
	for _, v := range articles {
		cover := mediaID(v.Photo)
		_, err = tx.Exec(ctx, `insert into news (id,title,published_on,category,excerpt,emoji,cover_media_id,content) values ($1,$2,$3,$4,$5,$6,$7,$8)`, v.ID, v.Title, v.Date, v.Category, v.Excerpt, v.Emoji, cover, v.Content)
		must(err)
		for p, photo := range v.Photos {
			_, err = tx.Exec(ctx, `insert into news_photos (news_id,media_id,position) values ($1,$2,$3)`, v.ID, mediaID(photo), p)
			must(err)
		}
	}
	for position, v := range services {
		cover := mediaID(v.Photo)
		_, err = tx.Exec(ctx, `insert into services (id,title,tagline,description,features,price,emoji,path,color,cover_media_id,detail,position) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`, v.ID, v.Title, v.Tagline, v.Description, v.Features, v.Price, v.Emoji, v.To, v.Color, cover, nullableJSON(v.Detail), position)
		must(err)
		for p, photo := range v.Photos {
			_, err = tx.Exec(ctx, `insert into service_photos (service_id,media_id,position) values ($1,$2,$3)`, v.ID, mediaID(photo), p)
			must(err)
		}
	}
	for position, v := range people {
		_, err = tx.Exec(ctx, `insert into staff (id,name,role,bio,initials,color,photo_media_id,position) values ($1,$2,$3,$4,$5,$6,$7,$8)`, v.ID, v.Name, v.Role, v.Bio, v.Initials, v.Color, mediaID(v.Photo), position)
		must(err)
	}
	for position, v := range pictures {
		_, err = tx.Exec(ctx, `insert into gallery (id,media_id,position) values ($1,$2,$3)`, v.ID, mediaID(v.Image), position)
		must(err)
	}
	for position, v := range h.HeroPhotos {
		_, err = tx.Exec(ctx, `insert into home_photos (media_id,position) values ($1,$2)`, mediaID(v), position)
		must(err)
	}
	for position, v := range documents {
		key := fmt.Sprintf("docs/placeholder-%d.pdf", v.ID)
		contents := strings.NewReader("%PDF-1.4\n%%EOF\n")
		must(minio.Upload(ctx, key, "application/pdf", int64(contents.Len()), contents))
		var id int64
		err = tx.QueryRow(ctx, `insert into media (object_key,mime,size_bytes) values ($1,'application/pdf',$2) returning id`, key, 15).Scan(&id)
		must(err)
		_, err = tx.Exec(ctx, `insert into documents (id,title,description,file_media_id,display_file_type,display_file_size,position) values ($1,$2,$3,$4,$5,$6,$7)`, v.ID, v.Title, v.Description, id, v.FileType, v.FileSize, position)
		must(err)
	}
	must(tx.Commit(ctx))
	fmt.Printf("Imported %d classes, %d news items, %d services, %d staff members and %d gallery photos.\n", len(classes), len(articles), len(services), len(people), len(pictures))
}

func get(ctx context.Context, client *http.Client, address string, target any) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	must(err)
	res, err := client.Do(req)
	must(err)
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		must(fmt.Errorf("GET %s: %s: %s", address, res.Status, body))
	}
	must(json.NewDecoder(res.Body).Decode(target))
}
func mediaKey(raw string) string {
	parsed, err := url.Parse(raw)
	must(err)
	path := parsed.Path
	if i := strings.Index(path, "/media/"); i >= 0 {
		path = path[i+len("/media/"):]
	}
	path = strings.TrimLeft(path, "/")
	if !strings.HasPrefix(path, "photos/") {
		must(fmt.Errorf("unsupported stub media path %q", raw))
	}
	return path
}
func nullableJSON(value json.RawMessage) any {
	if len(value) == 0 || string(value) == "null" {
		return nil
	}
	return value
}
func valueOr(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}
func must(err error) {
	if err != nil {
		panic(err)
	}
}

var _ pgx.Tx
