package store

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var ErrNotFound = errors.New("record not found")
var ErrMediaInUse = errors.New("media is in use")

type Store struct{ pool *pgxpool.Pool }

func New(pool *pgxpool.Pool) *Store             { return &Store{pool: pool} }
func (s *Store) Ping(ctx context.Context) error { return s.pool.Ping(ctx) }
func (s *Store) Close()                         { s.pool.Close() }

type Class struct {
	ID                                                          int64
	Title, AgeRange, Description, Schedule, Price, Emoji, Color string
	CoverKey                                                    *string
	Photos, FullDescription, Highlights                         []string
}
type News struct {
	ID                                    int64
	Title, Date, Category, Excerpt, Emoji string
	CoverKey                              *string
	Photos, Content                       []string
}
type Service struct {
	ID                                                     int64
	Title, Tagline, Description, Price, Emoji, Path, Color string
	CoverKey                                               *string
	Features, Photos                                       []string
	Detail                                                 []byte
}
type Staff struct {
	ID                               int64
	Name, Role, Bio, Initials, Color string
	PhotoKey                         *string
}
type Document struct {
	ID                                  int64
	Title, Description, ObjectKey, MIME string
	DisplayFileType, DisplayFileSize    *string
	SizeBytes                           int64
}
type GalleryItem struct {
	ID        int64
	ObjectKey string
}
type Media struct {
	ID        int64
	ObjectKey string
	MIME      string
	SizeBytes int64
	Alt       string
}
type Review struct {
	ID                 int64
	Author, Role, Text string
	Rating             int16
}

type AdminCounts struct {
	Classes, News, Services, Staff, Documents, Gallery int64
}

type AdminItem struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	IsPublished bool   `json:"isPublished"`
	Position    int    `json:"position"`
}
type ContentUpdate struct {
	Title       string
	IsPublished bool
	Position    int
}
type AdminUser struct {
	ID    int64
	Login string
	Role  string
}
type Session struct {
	Token     string
	User      AdminUser
	ExpiresAt time.Time
}

func (s *Store) BootstrapAdmin(ctx context.Context, login, password string) error {
	var exists bool
	if err := s.pool.QueryRow(ctx, `select exists(select 1 from admin_users)`).Scan(&exists); err != nil || exists {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(ctx, `insert into admin_users(login,password_hash,role) values ($1,$2,'administrator')`, login, string(hash))
	return err
}
func (s *Store) Login(ctx context.Context, login, password string) (Session, error) {
	var user AdminUser
	var hash string
	err := s.pool.QueryRow(ctx, `select id,login,role,password_hash from admin_users where login=$1`, login).Scan(&user.ID, &user.Login, &user.Role, &hash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return Session{}, ErrNotFound
	}
	token := uuid.NewString() + uuid.NewString()
	sum := sha256.Sum256([]byte(token))
	expires := time.Now().UTC().Add(24 * time.Hour)
	_, err = s.pool.Exec(ctx, `insert into admin_sessions(id,user_id,token_hash,expires_at) values ($1,$2,$3,$4)`, uuid.New(), user.ID, hex.EncodeToString(sum[:]), expires)
	return Session{Token: token, User: user, ExpiresAt: expires}, err
}
func (s *Store) Session(ctx context.Context, token string) (AdminUser, error) {
	sum := sha256.Sum256([]byte(token))
	var user AdminUser
	err := s.pool.QueryRow(ctx, `select u.id,u.login,u.role from admin_sessions s join admin_users u on u.id=s.user_id where s.token_hash=$1 and s.revoked_at is null and s.expires_at>now()`, hex.EncodeToString(sum[:])).Scan(&user.ID, &user.Login, &user.Role)
	if errors.Is(err, pgx.ErrNoRows) {
		return AdminUser{}, ErrNotFound
	}
	return user, err
}
func (s *Store) RevokeSessions(ctx context.Context, userID int64) error {
	_, err := s.pool.Exec(ctx, `update admin_sessions set revoked_at=now() where user_id=$1 and revoked_at is null`, userID)
	return err
}

const classSelect = `select c.id,c.title,c.age_range,c.description,c.schedule,c.price,c.emoji,c.color,cover.object_key,
 coalesce(array_agg(photo.object_key order by cp.position) filter (where photo.id is not null),'{}'),c.full_description,c.highlights
 from classes c left join media cover on cover.id=c.cover_media_id left join class_photos cp on cp.class_id=c.id left join media photo on photo.id=cp.media_id
 where c.is_published and (c.publication_status='published' or (c.publication_status='scheduled' and c.publish_at<=now())) %s group by c.id,cover.object_key order by c.position,c.id`
const newsSelect = `select n.id,n.title,n.published_on::text,n.category,n.excerpt,n.emoji,cover.object_key,
 coalesce(array_agg(photo.object_key order by np.position) filter (where photo.id is not null),'{}'),n.content
 from news n left join media cover on cover.id=n.cover_media_id left join news_photos np on np.news_id=n.id left join media photo on photo.id=np.media_id
 where n.is_published and (n.publication_status='published' or (n.publication_status='scheduled' and n.publish_at<=now())) %s group by n.id,cover.object_key order by n.published_on desc,n.id`
const serviceSelect = `select s.id,s.title,s.tagline,s.description,s.features,s.price,s.emoji,s.path,s.color,cover.object_key,
 coalesce(array_agg(photo.object_key order by sp.position) filter (where photo.id is not null),'{}'),s.detail
 from services s left join media cover on cover.id=s.cover_media_id left join service_photos sp on sp.service_id=s.id left join media photo on photo.id=sp.media_id
 where s.is_published and (s.publication_status='published' or (s.publication_status='scheduled' and s.publish_at<=now())) %s group by s.id,cover.object_key order by s.position,s.id`

func (s *Store) Classes(ctx context.Context) ([]Class, error) {
	return scanClasses(ctx, s.pool, fmt.Sprintf(classSelect, ""))
}
func (s *Store) Class(ctx context.Context, id int64) (Class, error) {
	rows, err := scanClasses(ctx, s.pool, fmt.Sprintf(classSelect, "and c.id=$1"), id)
	if err != nil {
		return Class{}, err
	}
	if len(rows) == 0 {
		return Class{}, ErrNotFound
	}
	return rows[0], nil
}
func scanClasses(ctx context.Context, q pgxQuerier, query string, args ...any) ([]Class, error) {
	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Class{}
	for rows.Next() {
		var v Class
		if err := rows.Scan(&v.ID, &v.Title, &v.AgeRange, &v.Description, &v.Schedule, &v.Price, &v.Emoji, &v.Color, &v.CoverKey, &v.Photos, &v.FullDescription, &v.Highlights); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

func (s *Store) News(ctx context.Context) ([]News, error) {
	return scanNews(ctx, s.pool, fmt.Sprintf(newsSelect, ""))
}
func (s *Store) NewsItem(ctx context.Context, id int64) (News, error) {
	rows, err := scanNews(ctx, s.pool, fmt.Sprintf(newsSelect, "and n.id=$1"), id)
	if err != nil {
		return News{}, err
	}
	if len(rows) == 0 {
		return News{}, ErrNotFound
	}
	return rows[0], nil
}
func scanNews(ctx context.Context, q pgxQuerier, query string, args ...any) ([]News, error) {
	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []News{}
	for rows.Next() {
		var v News
		if err := rows.Scan(&v.ID, &v.Title, &v.Date, &v.Category, &v.Excerpt, &v.Emoji, &v.CoverKey, &v.Photos, &v.Content); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (s *Store) Services(ctx context.Context) ([]Service, error) {
	return scanServices(ctx, s.pool, fmt.Sprintf(serviceSelect, ""))
}
func (s *Store) Service(ctx context.Context, id int64) (Service, error) {
	rows, err := scanServices(ctx, s.pool, fmt.Sprintf(serviceSelect, "and s.id=$1"), id)
	if err != nil {
		return Service{}, err
	}
	if len(rows) == 0 {
		return Service{}, ErrNotFound
	}
	return rows[0], nil
}
func scanServices(ctx context.Context, q pgxQuerier, query string, args ...any) ([]Service, error) {
	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Service{}
	for rows.Next() {
		var v Service
		if err := rows.Scan(&v.ID, &v.Title, &v.Tagline, &v.Description, &v.Features, &v.Price, &v.Emoji, &v.Path, &v.Color, &v.CoverKey, &v.Photos, &v.Detail); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (s *Store) Staff(ctx context.Context) ([]Staff, error) {
	rows, err := s.pool.Query(ctx, `select s.id,s.name,s.role,s.bio,s.initials,s.color,m.object_key from staff s left join media m on m.id=s.photo_media_id where s.is_published order by s.position,s.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Staff{}
	for rows.Next() {
		var v Staff
		if err := rows.Scan(&v.ID, &v.Name, &v.Role, &v.Bio, &v.Initials, &v.Color, &v.PhotoKey); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (s *Store) Documents(ctx context.Context) ([]Document, error) {
	rows, err := s.pool.Query(ctx, `select d.id,d.title,d.description,m.object_key,m.mime,m.size_bytes,d.display_file_type,d.display_file_size from documents d join media m on m.id=d.file_media_id where d.is_published order by d.position,d.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Document{}
	for rows.Next() {
		var v Document
		if err := rows.Scan(&v.ID, &v.Title, &v.Description, &v.ObjectKey, &v.MIME, &v.SizeBytes, &v.DisplayFileType, &v.DisplayFileSize); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (s *Store) Gallery(ctx context.Context) ([]GalleryItem, error) {
	rows, err := s.pool.Query(ctx, `select g.id,m.object_key from gallery g join media m on m.id=g.media_id order by g.position,g.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []GalleryItem{}
	for rows.Next() {
		var v GalleryItem
		if err := rows.Scan(&v.ID, &v.ObjectKey); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (s *Store) HomePhotoKeys(ctx context.Context) ([]string, error) {
	rows, err := s.pool.Query(ctx, `select m.object_key from home_photos h join media m on m.id=h.media_id order by h.position,m.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		out = append(out, key)
	}
	return out, rows.Err()
}
func (s *Store) Reviews(ctx context.Context) ([]Review, error) {
	rows, err := s.pool.Query(ctx, `select id,author,role,text,rating from reviews where is_published order by position,id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Review{}
	for rows.Next() {
		var v Review
		if err := rows.Scan(&v.ID, &v.Author, &v.Role, &v.Text, &v.Rating); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (s *Store) Site(ctx context.Context) ([]byte, error) {
	var data []byte
	err := s.pool.QueryRow(ctx, `select data from site_settings where id=true`).Scan(&data)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return data, err
}

// Media returns a stable page of library records. Object URLs are deliberately
// resolved in the HTTP layer: storage endpoints must not leak into the DB.
func (s *Store) Media(ctx context.Context, query string, limit, offset int) ([]Media, int64, error) {
	pattern := "%" + query + "%"
	var total int64
	if err := s.pool.QueryRow(ctx, `select count(*) from media where object_key ilike $1 or alt ilike $1`, pattern).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := s.pool.Query(ctx, `select id,object_key,mime,size_bytes,alt from media
		where object_key ilike $1 or alt ilike $1 order by created_at desc,id desc limit $2 offset $3`, pattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := []Media{}
	for rows.Next() {
		var item Media
		if err := rows.Scan(&item.ID, &item.ObjectKey, &item.MIME, &item.SizeBytes, &item.Alt); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func (s *Store) CreateMedia(ctx context.Context, objectKey, mime string, sizeBytes int64, alt string) (Media, error) {
	var item Media
	err := s.pool.QueryRow(ctx, `insert into media (object_key,mime,size_bytes,alt) values ($1,$2,$3,$4)
		returning id,object_key,mime,size_bytes,alt`, objectKey, mime, sizeBytes, alt).Scan(&item.ID, &item.ObjectKey, &item.MIME, &item.SizeBytes, &item.Alt)
	return item, err
}

// DeleteMedia only deletes unattached files. This avoids broken images and
// documents when an editor removes a file that is still used by content.
func (s *Store) DeleteMedia(ctx context.Context, id int64) (string, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	var key string
	if err := tx.QueryRow(ctx, `select object_key from media where id=$1 for update`, id).Scan(&key); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	var used bool
	err = tx.QueryRow(ctx, `select exists(
		select 1 from classes where cover_media_id=$1 union all select 1 from class_photos where media_id=$1
		union all select 1 from news where cover_media_id=$1 union all select 1 from news_photos where media_id=$1
		union all select 1 from services where cover_media_id=$1 union all select 1 from service_photos where media_id=$1
		union all select 1 from staff where photo_media_id=$1 union all select 1 from documents where file_media_id=$1
		union all select 1 from gallery where media_id=$1 union all select 1 from home_photos where media_id=$1)`, id).Scan(&used)
	if err != nil {
		return "", err
	}
	if used {
		return "", ErrMediaInUse
	}
	if _, err := tx.Exec(ctx, `delete from media where id=$1`, id); err != nil {
		return "", err
	}
	if err := tx.Commit(ctx); err != nil {
		return "", err
	}
	return key, nil
}

func (s *Store) AdminCounts(ctx context.Context) (AdminCounts, error) {
	var result AdminCounts
	row := s.pool.QueryRow(ctx, `select
  (select count(*) from classes), (select count(*) from news),
  (select count(*) from services), (select count(*) from staff),
  (select count(*) from documents), (select count(*) from gallery)`)
	if err := row.Scan(&result.Classes, &result.News, &result.Services, &result.Staff, &result.Documents, &result.Gallery); err != nil {
		return AdminCounts{}, err
	}
	return result, nil
}

func (s *Store) SetPublished(ctx context.Context, resource string, id int64, published bool) error {
	tables := map[string]string{
		"classes": "classes", "news": "news", "services": "services",
		"staff": "staff", "documents": "documents",
	}
	table, ok := tables[resource]
	if !ok {
		return ErrNotFound
	}
	command, err := s.pool.Exec(ctx, `update `+table+` set is_published=$1 where id=$2`, published, id)
	if err != nil {
		return err
	}
	if command.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// UpdateAdminItem covers fields shared by all editorial lists. Entity-specific
// content remains in its own columns, while title, order and visibility use one
// safe, whitelisted path.
func (s *Store) UpdateAdminItem(ctx context.Context, resource string, id int64, input ContentUpdate) error {
	tables := map[string]string{"classes": "classes", "news": "news", "services": "services", "staff": "staff", "documents": "documents"}
	titleColumn := map[string]string{"staff": "name"}
	table, ok := tables[resource]
	if !ok {
		return ErrNotFound
	}
	column := titleColumn[resource]
	if column == "" {
		column = "title"
	}
	command, err := s.pool.Exec(ctx, `update `+table+` set `+column+`=$1,position=$2,is_published=$3,publication_status=case when $3 then 'published' else 'draft' end,publish_at=null where id=$4`, input.Title, input.Position, input.IsPublished, id)
	if err != nil {
		return err
	}
	if command.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
func (s *Store) DeleteAdminItem(ctx context.Context, resource string, id int64) error {
	tables := map[string]string{"classes": "classes", "news": "news", "services": "services", "staff": "staff", "documents": "documents"}
	table, ok := tables[resource]
	if !ok {
		return ErrNotFound
	}
	command, err := s.pool.Exec(ctx, `delete from `+table+` where id=$1`, id)
	if err != nil {
		return err
	}
	if command.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) AdminContent(ctx context.Context, resource string, id int64) ([]byte, error) {
	tables := map[string]string{"classes": "classes", "news": "news", "services": "services", "staff": "staff", "documents": "documents", "gallery": "gallery"}
	table, ok := tables[resource]
	if !ok {
		return nil, ErrNotFound
	}
	var data []byte
	err := s.pool.QueryRow(ctx, `select to_jsonb(t) from `+table+` t where id=$1`, id).Scan(&data)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return data, err
}

// UpdateAdminContent applies a partial JSON document to editable columns. The
// table and column lists are fixed in code, so clients cannot inject SQL or
// write identifiers/audit timestamps.
func (s *Store) UpdateAdminContent(ctx context.Context, resource string, id int64, data []byte) error {
	queries := map[string]string{
		"classes":   `update classes c set (title,age_range,description,schedule,price,emoji,color,cover_media_id,full_description,highlights,position,is_published,publication_status,publish_at)=(select title,age_range,description,schedule,price,emoji,color,cover_media_id,full_description,highlights,position,is_published,publication_status,publish_at from jsonb_populate_record(c,$1::jsonb)) where c.id=$2`,
		"news":      `update news n set (title,published_on,category,excerpt,emoji,cover_media_id,content,is_published,publication_status,publish_at)=(select title,published_on,category,excerpt,emoji,cover_media_id,content,is_published,publication_status,publish_at from jsonb_populate_record(n,$1::jsonb)) where n.id=$2`,
		"services":  `update services s set (title,tagline,description,features,price,emoji,path,color,cover_media_id,detail,position,is_published,publication_status,publish_at)=(select title,tagline,description,features,price,emoji,path,color,cover_media_id,detail,position,is_published,publication_status,publish_at from jsonb_populate_record(s,$1::jsonb)) where s.id=$2`,
		"staff":     `update staff s set (name,role,bio,initials,color,photo_media_id,position,is_published,publication_status,publish_at)=(select name,role,bio,initials,color,photo_media_id,position,is_published,publication_status,publish_at from jsonb_populate_record(s,$1::jsonb)) where s.id=$2`,
		"documents": `update documents d set (title,description,file_media_id,display_file_type,display_file_size,position,is_published,publication_status,publish_at)=(select title,description,file_media_id,display_file_type,display_file_size,position,is_published,publication_status,publish_at from jsonb_populate_record(d,$1::jsonb)) where d.id=$2`,
		"gallery":   `update gallery g set (media_id,title,position)=(select media_id,title,position from jsonb_populate_record(g,$1::jsonb)) where g.id=$2`,
	}
	query, ok := queries[resource]
	if !ok {
		return ErrNotFound
	}
	command, err := s.pool.Exec(ctx, query, data, id)
	if err != nil {
		return err
	}
	if command.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) AdminItems(ctx context.Context, resource string) ([]AdminItem, error) {
	queries := map[string]string{
		"classes":   `select id,title,is_published,position from classes order by position,id`,
		"news":      `select id,title,is_published,0 from news order by published_on desc,id`,
		"services":  `select id,title,is_published,position from services order by position,id`,
		"staff":     `select id,name,is_published,position from staff order by position,id`,
		"documents": `select id,title,is_published,position from documents order by position,id`,
	}
	query, ok := queries[resource]
	if !ok {
		return nil, ErrNotFound
	}
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []AdminItem{}
	for rows.Next() {
		var item AdminItem
		if err := rows.Scan(&item.ID, &item.Title, &item.IsPublished, &item.Position); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

type pgxQuerier interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}
