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
type SummerShift = Class
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
	Classes, SummerShifts, News, Services, Staff, Documents, Gallery, Reviews int64
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
type AdminClass struct {
	ID                                            int64
	Title, AgeRange, Description, Schedule, Price string
	Emoji, Color                                  string
	CoverMediaID                                  *int64
	FullDescription, Highlights                   []string
	PhotoMediaIDs                                 []int64
	Position                                      int
	IsPublished                                   bool
}
type AdminSummerShift = AdminClass
type AdminStaff struct {
	ID                               int64
	Name, Role, Bio, Initials, Color string
	PhotoMediaID                     *int64
	Position                         int
	IsPublished                      bool
}
type AdminDocument struct {
	ID                               int64
	Title, Description               string
	FileMediaID                      int64
	DisplayFileType, DisplayFileSize *string
	Position                         int
	IsPublished                      bool
}
type AdminReview struct {
	ID                 int64
	Author, Role, Text string
	Rating             int16
	Position           int
	IsPublished        bool
}
type AdminGalleryItem struct {
	ID, MediaID int64
	Title       string
	Position    int
	Media       Media
}
type AdminNews struct {
	ID                             int64
	Title, Date, Category, Excerpt string
	Emoji                          string
	CoverMediaID                   *int64
	Content                        []string
	PhotoMediaIDs                  []int64
	IsPublished                    bool
}
type AdminService struct {
	ID                                                     int64
	Title, Tagline, Description, Price, Emoji, Path, Color string
	Features                                               []string
	CoverMediaID                                           *int64
	PhotoMediaIDs                                          []int64
	Position                                               int
	IsPublished                                            bool
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
const summerShiftSelect = `select c.id,c.title,c.age_range,c.description,c.schedule,c.price,c.emoji,c.color,cover.object_key,
 coalesce(array_agg(photo.object_key order by cp.position) filter (where photo.id is not null),'{}'),c.full_description,c.highlights
 from summer_shifts c left join media cover on cover.id=c.cover_media_id left join summer_shift_photos cp on cp.summer_shift_id=c.id left join media photo on photo.id=cp.media_id
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
func (s *Store) SummerShifts(ctx context.Context) ([]SummerShift, error) {
	return scanClasses(ctx, s.pool, fmt.Sprintf(summerShiftSelect, ""))
}
func (s *Store) SummerShift(ctx context.Context, id int64) (SummerShift, error) {
	rows, err := scanClasses(ctx, s.pool, fmt.Sprintf(summerShiftSelect, "and c.id=$1"), id)
	if err != nil {
		return SummerShift{}, err
	}
	if len(rows) == 0 {
		return SummerShift{}, ErrNotFound
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

func (s *Store) MediaByIDs(ctx context.Context, ids []int64) ([]Media, error) {
	if len(ids) == 0 {
		return []Media{}, nil
	}
	rows, err := s.pool.Query(ctx, `select id,object_key,mime,size_bytes,alt from media where id=any($1)`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	byID := make(map[int64]Media, len(ids))
	for rows.Next() {
		var item Media
		if err := rows.Scan(&item.ID, &item.ObjectKey, &item.MIME, &item.SizeBytes, &item.Alt); err != nil {
			return nil, err
		}
		byID[item.ID] = item
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	result := make([]Media, 0, len(ids))
	for _, id := range ids {
		if item, ok := byID[id]; ok {
			result = append(result, item)
		}
	}
	return result, nil
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
		union all select 1 from summer_shifts where cover_media_id=$1 union all select 1 from summer_shift_photos where media_id=$1
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

// DeleteClassMedia atomically detaches a photo from one class and removes its
// database record. The object key is returned so the HTTP layer can remove the
// corresponding object from storage.
func (s *Store) DeleteClassMedia(ctx context.Context, classID, mediaID int64) (string, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	var key string
	if err := tx.QueryRow(ctx, `select object_key from media where id=$1 for update`, mediaID).Scan(&key); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	command, err := tx.Exec(ctx, `update classes set cover_media_id=null,updated_at=now() where id=$1 and cover_media_id=$2`, classID, mediaID)
	if err != nil {
		return "", err
	}
	removed := command.RowsAffected() > 0
	command, err = tx.Exec(ctx, `delete from class_photos where class_id=$1 and media_id=$2`, classID, mediaID)
	if err != nil {
		return "", err
	}
	removed = removed || command.RowsAffected() > 0
	if !removed {
		return "", ErrNotFound
	}
	var used bool
	err = tx.QueryRow(ctx, `select exists(
		select 1 from classes where cover_media_id=$1 union all select 1 from class_photos where media_id=$1
		union all select 1 from news where cover_media_id=$1 union all select 1 from news_photos where media_id=$1
		union all select 1 from services where cover_media_id=$1 union all select 1 from service_photos where media_id=$1
		union all select 1 from staff where photo_media_id=$1 union all select 1 from documents where file_media_id=$1
		union all select 1 from gallery where media_id=$1 union all select 1 from home_photos where media_id=$1)`, mediaID).Scan(&used)
	if err != nil {
		return "", err
	}
	if used {
		return "", ErrMediaInUse
	}
	if _, err = tx.Exec(ctx, `delete from media where id=$1`, mediaID); err != nil {
		return "", err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", err
	}
	return key, nil
}

func (s *Store) DeleteSummerShiftMedia(ctx context.Context, shiftID, mediaID int64) (string, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	var key string
	if err := tx.QueryRow(ctx, `select object_key from media where id=$1 for update`, mediaID).Scan(&key); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	command, err := tx.Exec(ctx, `update summer_shifts set cover_media_id=null,updated_at=now() where id=$1 and cover_media_id=$2`, shiftID, mediaID)
	if err != nil {
		return "", err
	}
	removed := command.RowsAffected() > 0
	command, err = tx.Exec(ctx, `delete from summer_shift_photos where summer_shift_id=$1 and media_id=$2`, shiftID, mediaID)
	if err != nil {
		return "", err
	}
	removed = removed || command.RowsAffected() > 0
	if !removed {
		return "", ErrNotFound
	}
	var used bool
	err = tx.QueryRow(ctx, `select exists(
		select 1 from classes where cover_media_id=$1 union all select 1 from class_photos where media_id=$1
		union all select 1 from summer_shifts where cover_media_id=$1 union all select 1 from summer_shift_photos where media_id=$1
		union all select 1 from news where cover_media_id=$1 union all select 1 from news_photos where media_id=$1
		union all select 1 from services where cover_media_id=$1 union all select 1 from service_photos where media_id=$1
		union all select 1 from staff where photo_media_id=$1 union all select 1 from documents where file_media_id=$1
		union all select 1 from gallery where media_id=$1 union all select 1 from home_photos where media_id=$1)`, mediaID).Scan(&used)
	if err != nil {
		return "", err
	}
	if used {
		return "", ErrMediaInUse
	}
	if _, err = tx.Exec(ctx, `delete from media where id=$1`, mediaID); err != nil {
		return "", err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", err
	}
	return key, nil
}

func (s *Store) DeleteNewsMedia(ctx context.Context, newsID, mediaID int64) (string, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	var key string
	if err := tx.QueryRow(ctx, `select object_key from media where id=$1 for update`, mediaID).Scan(&key); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	command, err := tx.Exec(ctx, `update news set cover_media_id=null,updated_at=now() where id=$1 and cover_media_id=$2`, newsID, mediaID)
	if err != nil {
		return "", err
	}
	removed := command.RowsAffected() > 0
	command, err = tx.Exec(ctx, `delete from news_photos where news_id=$1 and media_id=$2`, newsID, mediaID)
	if err != nil {
		return "", err
	}
	removed = removed || command.RowsAffected() > 0
	if !removed {
		return "", ErrNotFound
	}
	var used bool
	err = tx.QueryRow(ctx, `select exists(
		select 1 from classes where cover_media_id=$1 union all select 1 from class_photos where media_id=$1
		union all select 1 from news where cover_media_id=$1 union all select 1 from news_photos where media_id=$1
		union all select 1 from services where cover_media_id=$1 union all select 1 from service_photos where media_id=$1
		union all select 1 from staff where photo_media_id=$1 union all select 1 from documents where file_media_id=$1
		union all select 1 from gallery where media_id=$1 union all select 1 from home_photos where media_id=$1)`, mediaID).Scan(&used)
	if err != nil {
		return "", err
	}
	if used {
		return "", ErrMediaInUse
	}
	if _, err = tx.Exec(ctx, `delete from media where id=$1`, mediaID); err != nil {
		return "", err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", err
	}
	return key, nil
}

func (s *Store) DeleteServiceMedia(ctx context.Context, serviceID, mediaID int64) (string, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	var key string
	if err := tx.QueryRow(ctx, `select object_key from media where id=$1 for update`, mediaID).Scan(&key); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	command, err := tx.Exec(ctx, `update services set cover_media_id=null where id=$1 and cover_media_id=$2`, serviceID, mediaID)
	if err != nil {
		return "", err
	}
	removed := command.RowsAffected() > 0
	command, err = tx.Exec(ctx, `delete from service_photos where service_id=$1 and media_id=$2`, serviceID, mediaID)
	if err != nil {
		return "", err
	}
	removed = removed || command.RowsAffected() > 0
	if !removed {
		return "", ErrNotFound
	}
	var used bool
	err = tx.QueryRow(ctx, `select exists(select 1 from classes where cover_media_id=$1 union all select 1 from class_photos where media_id=$1 union all select 1 from news where cover_media_id=$1 union all select 1 from news_photos where media_id=$1 union all select 1 from services where cover_media_id=$1 union all select 1 from service_photos where media_id=$1 union all select 1 from staff where photo_media_id=$1 union all select 1 from documents where file_media_id=$1 union all select 1 from gallery where media_id=$1 union all select 1 from home_photos where media_id=$1)`, mediaID).Scan(&used)
	if err != nil {
		return "", err
	}
	if used {
		return "", ErrMediaInUse
	}
	if _, err = tx.Exec(ctx, `delete from media where id=$1`, mediaID); err != nil {
		return "", err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", err
	}
	return key, nil
}

func (s *Store) AdminCounts(ctx context.Context) (AdminCounts, error) {
	var result AdminCounts
	row := s.pool.QueryRow(ctx, `select
	  (select count(*) from classes), (select count(*) from summer_shifts), (select count(*) from news),
  (select count(*) from services), (select count(*) from staff),
  (select count(*) from documents), (select count(*) from gallery), (select count(*) from reviews)`)
	if err := row.Scan(&result.Classes, &result.SummerShifts, &result.News, &result.Services, &result.Staff, &result.Documents, &result.Gallery, &result.Reviews); err != nil {
		return AdminCounts{}, err
	}
	return result, nil
}

func (s *Store) SetPublished(ctx context.Context, resource string, id int64, published bool) error {
	tables := map[string]string{
		"classes": "classes", "summer-shifts": "summer_shifts", "news": "news", "services": "services",
		"staff": "staff", "documents": "documents", "reviews": "reviews",
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
	if resource == "reviews" {
		command, err := s.pool.Exec(ctx, `update reviews set author=$1,position=$2,is_published=$3 where id=$4`, input.Title, input.Position, input.IsPublished, id)
		if err != nil {
			return err
		}
		if command.RowsAffected() == 0 {
			return ErrNotFound
		}
		return nil
	}
	tables := map[string]string{"classes": "classes", "summer-shifts": "summer_shifts", "news": "news", "services": "services", "staff": "staff", "documents": "documents", "reviews": "reviews"}
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
	tables := map[string]string{"classes": "classes", "summer-shifts": "summer_shifts", "news": "news", "services": "services", "staff": "staff", "documents": "documents", "reviews": "reviews"}
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
func (s *Store) DeleteAdminAsset(ctx context.Context, resource string, id int64) (string, error) {
	if resource != "staff" && resource != "documents" {
		return "", ErrNotFound
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	var key *string
	var mediaID *int64
	query := `select m.id,m.object_key from staff s left join media m on m.id=s.photo_media_id where s.id=$1`
	table := "staff"
	if resource == "documents" {
		query = `select m.id,m.object_key from documents d join media m on m.id=d.file_media_id where d.id=$1 for update`
		table = "documents"
	}
	if err = tx.QueryRow(ctx, query, id).Scan(&mediaID, &key); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	if _, err = tx.Exec(ctx, `delete from `+table+` where id=$1`, id); err != nil {
		return "", err
	}
	if mediaID != nil {
		if _, err = tx.Exec(ctx, `delete from media where id=$1`, *mediaID); err != nil {
			return "", err
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return "", err
	}
	if key == nil {
		return "", nil
	}
	return *key, nil
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
		"classes":       `select id,title,is_published,position from classes order by position,id`,
		"summer-shifts": `select id,title,is_published,position from summer_shifts order by position,id`,
		"news":          `select id,title,is_published,0 from news order by published_on desc,id`,
		"services":      `select id,title,is_published,position from services order by position,id`,
		"staff":         `select id,name,is_published,position from staff order by position,id`,
		"documents":     `select id,title,is_published,position from documents order by position,id`,
		"reviews":       `select id,author,is_published,position from reviews order by position,id`,
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

func (s *Store) AdminClass(ctx context.Context, id int64) (AdminClass, error) {
	var item AdminClass
	err := s.pool.QueryRow(ctx, `select c.id,c.title,c.age_range,c.description,c.schedule,c.price,c.emoji,c.color,
		c.cover_media_id,c.full_description,c.highlights,c.position,c.is_published,
		coalesce(array_agg(cp.media_id order by cp.position) filter (where cp.media_id is not null),'{}')
		from classes c left join class_photos cp on cp.class_id=c.id where c.id=$1 group by c.id`, id).
		Scan(&item.ID, &item.Title, &item.AgeRange, &item.Description, &item.Schedule, &item.Price,
			&item.Emoji, &item.Color, &item.CoverMediaID, &item.FullDescription, &item.Highlights,
			&item.Position, &item.IsPublished, &item.PhotoMediaIDs)
	if errors.Is(err, pgx.ErrNoRows) {
		return AdminClass{}, ErrNotFound
	}
	return item, err
}

func (s *Store) SaveAdminClass(ctx context.Context, item AdminClass) (int64, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	if item.ID == 0 {
		err = tx.QueryRow(ctx, `insert into classes
			(title,age_range,description,schedule,price,emoji,color,cover_media_id,full_description,highlights,position,is_published,publication_status)
			values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,case when $12 then 'published' else 'draft' end) returning id`,
			item.Title, item.AgeRange, item.Description, item.Schedule, item.Price, item.Emoji, item.Color, item.CoverMediaID,
			item.FullDescription, item.Highlights, item.Position, item.IsPublished).Scan(&item.ID)
	} else {
		command, updateErr := tx.Exec(ctx, `update classes set title=$1,age_range=$2,description=$3,schedule=$4,price=$5,
			emoji=$6,color=$7,cover_media_id=$8,full_description=$9,highlights=$10,position=$11,is_published=$12,
			publication_status=case when $12 then 'published' else 'draft' end,publish_at=null,updated_at=now() where id=$13`,
			item.Title, item.AgeRange, item.Description, item.Schedule, item.Price, item.Emoji, item.Color, item.CoverMediaID,
			item.FullDescription, item.Highlights, item.Position, item.IsPublished, item.ID)
		err = updateErr
		if err == nil && command.RowsAffected() == 0 {
			err = ErrNotFound
		}
	}
	if err != nil {
		return 0, err
	}
	if _, err = tx.Exec(ctx, `delete from class_photos where class_id=$1`, item.ID); err != nil {
		return 0, err
	}
	for position, mediaID := range item.PhotoMediaIDs {
		if _, err = tx.Exec(ctx, `insert into class_photos(class_id,media_id,position) values ($1,$2,$3)`, item.ID, mediaID, position); err != nil {
			return 0, err
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return item.ID, nil
}

func (s *Store) AdminSummerShift(ctx context.Context, id int64) (AdminSummerShift, error) {
	var item AdminSummerShift
	err := s.pool.QueryRow(ctx, `select c.id,c.title,c.age_range,c.description,c.schedule,c.price,c.emoji,c.color,
		c.cover_media_id,c.full_description,c.highlights,c.position,c.is_published,
		coalesce(array_agg(cp.media_id order by cp.position) filter (where cp.media_id is not null),'{}')
		from summer_shifts c left join summer_shift_photos cp on cp.summer_shift_id=c.id where c.id=$1 group by c.id`, id).
		Scan(&item.ID, &item.Title, &item.AgeRange, &item.Description, &item.Schedule, &item.Price, &item.Emoji, &item.Color,
			&item.CoverMediaID, &item.FullDescription, &item.Highlights, &item.Position, &item.IsPublished, &item.PhotoMediaIDs)
	if errors.Is(err, pgx.ErrNoRows) {
		return AdminSummerShift{}, ErrNotFound
	}
	return item, err
}

func (s *Store) SaveAdminSummerShift(ctx context.Context, item AdminSummerShift) (int64, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	if item.ID == 0 {
		err = tx.QueryRow(ctx, `insert into summer_shifts (title,age_range,description,schedule,price,emoji,color,cover_media_id,full_description,highlights,position,is_published) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) returning id`, item.Title, item.AgeRange, item.Description, item.Schedule, item.Price, item.Emoji, item.Color, item.CoverMediaID, item.FullDescription, item.Highlights, item.Position, item.IsPublished).Scan(&item.ID)
	} else {
		command, updateErr := tx.Exec(ctx, `update summer_shifts set title=$1,age_range=$2,description=$3,schedule=$4,price=$5,emoji=$6,color=$7,cover_media_id=$8,full_description=$9,highlights=$10,position=$11,is_published=$12,updated_at=now() where id=$13`, item.Title, item.AgeRange, item.Description, item.Schedule, item.Price, item.Emoji, item.Color, item.CoverMediaID, item.FullDescription, item.Highlights, item.Position, item.IsPublished, item.ID)
		err = updateErr
		if err == nil && command.RowsAffected() == 0 {
			err = ErrNotFound
		}
	}
	if err != nil {
		return 0, err
	}
	if _, err = tx.Exec(ctx, `delete from summer_shift_photos where summer_shift_id=$1`, item.ID); err != nil {
		return 0, err
	}
	for position, mediaID := range item.PhotoMediaIDs {
		if _, err = tx.Exec(ctx, `insert into summer_shift_photos(summer_shift_id,media_id,position) values ($1,$2,$3)`, item.ID, mediaID, position); err != nil {
			return 0, err
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return item.ID, nil
}

func (s *Store) AdminStaff(ctx context.Context, id int64) (AdminStaff, error) {
	var item AdminStaff
	err := s.pool.QueryRow(ctx, `select id,name,role,bio,initials,color,photo_media_id,position,is_published from staff where id=$1`, id).Scan(&item.ID, &item.Name, &item.Role, &item.Bio, &item.Initials, &item.Color, &item.PhotoMediaID, &item.Position, &item.IsPublished)
	if errors.Is(err, pgx.ErrNoRows) {
		return AdminStaff{}, ErrNotFound
	}
	return item, err
}
func (s *Store) SaveAdminStaff(ctx context.Context, item AdminStaff) (int64, error) {
	if item.ID == 0 {
		err := s.pool.QueryRow(ctx, `insert into staff(name,role,bio,initials,color,photo_media_id,position,is_published,publication_status) values($1,$2,$3,$4,$5,$6,$7,$8,case when $8 then 'published' else 'draft' end) returning id`, item.Name, item.Role, item.Bio, item.Initials, item.Color, item.PhotoMediaID, item.Position, item.IsPublished).Scan(&item.ID)
		return item.ID, err
	}
	command, err := s.pool.Exec(ctx, `update staff set name=$1,role=$2,bio=$3,initials=$4,color=$5,photo_media_id=$6,position=$7,is_published=$8,publication_status=case when $8 then 'published' else 'draft' end,publish_at=null where id=$9`, item.Name, item.Role, item.Bio, item.Initials, item.Color, item.PhotoMediaID, item.Position, item.IsPublished, item.ID)
	if err == nil && command.RowsAffected() == 0 {
		err = ErrNotFound
	}
	return item.ID, err
}
func (s *Store) AdminDocument(ctx context.Context, id int64) (AdminDocument, error) {
	var item AdminDocument
	err := s.pool.QueryRow(ctx, `select id,title,description,file_media_id,display_file_type,display_file_size,position,is_published from documents where id=$1`, id).Scan(&item.ID, &item.Title, &item.Description, &item.FileMediaID, &item.DisplayFileType, &item.DisplayFileSize, &item.Position, &item.IsPublished)
	if errors.Is(err, pgx.ErrNoRows) {
		return AdminDocument{}, ErrNotFound
	}
	return item, err
}
func (s *Store) SaveAdminDocument(ctx context.Context, item AdminDocument) (int64, error) {
	if item.ID == 0 {
		err := s.pool.QueryRow(ctx, `insert into documents(title,description,file_media_id,display_file_type,display_file_size,position,is_published,publication_status) values($1,$2,$3,$4,$5,$6,$7,case when $7 then 'published' else 'draft' end) returning id`, item.Title, item.Description, item.FileMediaID, item.DisplayFileType, item.DisplayFileSize, item.Position, item.IsPublished).Scan(&item.ID)
		return item.ID, err
	}
	command, err := s.pool.Exec(ctx, `update documents set title=$1,description=$2,file_media_id=$3,display_file_type=$4,display_file_size=$5,position=$6,is_published=$7,publication_status=case when $7 then 'published' else 'draft' end,publish_at=null where id=$8`, item.Title, item.Description, item.FileMediaID, item.DisplayFileType, item.DisplayFileSize, item.Position, item.IsPublished, item.ID)
	if err == nil && command.RowsAffected() == 0 {
		err = ErrNotFound
	}
	return item.ID, err
}

func (s *Store) DeleteStaffMedia(ctx context.Context, staffID, mediaID int64) (string, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	var key string
	if err = tx.QueryRow(ctx, `select object_key from media where id=$1 for update`, mediaID).Scan(&key); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	command, err := tx.Exec(ctx, `update staff set photo_media_id=null where id=$1 and photo_media_id=$2`, staffID, mediaID)
	if err != nil {
		return "", err
	}
	if command.RowsAffected() == 0 {
		return "", ErrNotFound
	}
	if _, err = tx.Exec(ctx, `delete from media where id=$1`, mediaID); err != nil {
		return "", err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", err
	}
	return key, nil
}
func (s *Store) AdminReview(ctx context.Context, id int64) (AdminReview, error) {
	var item AdminReview
	err := s.pool.QueryRow(ctx, `select id,author,role,text,rating,position,is_published from reviews where id=$1`, id).Scan(&item.ID, &item.Author, &item.Role, &item.Text, &item.Rating, &item.Position, &item.IsPublished)
	if errors.Is(err, pgx.ErrNoRows) {
		return AdminReview{}, ErrNotFound
	}
	return item, err
}
func (s *Store) SaveAdminReview(ctx context.Context, item AdminReview) (int64, error) {
	if item.ID == 0 {
		err := s.pool.QueryRow(ctx, `insert into reviews(author,role,text,rating,position,is_published) values($1,$2,$3,$4,$5,$6) returning id`, item.Author, item.Role, item.Text, item.Rating, item.Position, item.IsPublished).Scan(&item.ID)
		return item.ID, err
	}
	command, err := s.pool.Exec(ctx, `update reviews set author=$1,role=$2,text=$3,rating=$4,position=$5,is_published=$6 where id=$7`, item.Author, item.Role, item.Text, item.Rating, item.Position, item.IsPublished, item.ID)
	if err == nil && command.RowsAffected() == 0 {
		err = ErrNotFound
	}
	return item.ID, err
}
func (s *Store) AdminGallery(ctx context.Context) ([]AdminGalleryItem, error) {
	rows, err := s.pool.Query(ctx, `select g.id,g.media_id,g.title,g.position,m.object_key,m.mime,m.size_bytes,m.alt from gallery g join media m on m.id=g.media_id order by g.position,g.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []AdminGalleryItem{}
	for rows.Next() {
		var v AdminGalleryItem
		if err = rows.Scan(&v.ID, &v.MediaID, &v.Title, &v.Position, &v.Media.ObjectKey, &v.Media.MIME, &v.Media.SizeBytes, &v.Media.Alt); err != nil {
			return nil, err
		}
		v.Media.ID = v.MediaID
		out = append(out, v)
	}
	return out, rows.Err()
}
func (s *Store) AddGalleryItem(ctx context.Context, mediaID int64, title string) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `insert into gallery(media_id,title,position) values($1,$2,coalesce((select max(position)+1 from gallery),0)) returning id`, mediaID, title).Scan(&id)
	return id, err
}
func (s *Store) DeleteGalleryItem(ctx context.Context, id int64) (string, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	var mediaID int64
	var key string
	if err = tx.QueryRow(ctx, `select m.id,m.object_key from gallery g join media m on m.id=g.media_id where g.id=$1`, id).Scan(&mediaID, &key); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	if _, err = tx.Exec(ctx, `delete from gallery where id=$1`, id); err != nil {
		return "", err
	}
	if _, err = tx.Exec(ctx, `delete from media where id=$1`, mediaID); err != nil {
		return "", err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", err
	}
	return key, nil
}

func (s *Store) AdminNews(ctx context.Context, id int64) (AdminNews, error) {
	var item AdminNews
	err := s.pool.QueryRow(ctx, `select n.id,n.title,n.published_on::text,n.category,n.excerpt,n.emoji,n.cover_media_id,n.content,n.is_published,
		coalesce(array_agg(np.media_id order by np.position) filter (where np.media_id is not null),'{}')
		from news n left join news_photos np on np.news_id=n.id where n.id=$1 group by n.id`, id).
		Scan(&item.ID, &item.Title, &item.Date, &item.Category, &item.Excerpt, &item.Emoji, &item.CoverMediaID, &item.Content, &item.IsPublished, &item.PhotoMediaIDs)
	if errors.Is(err, pgx.ErrNoRows) {
		return AdminNews{}, ErrNotFound
	}
	return item, err
}

func (s *Store) NewsCategories(ctx context.Context) ([]string, error) {
	rows, err := s.pool.Query(ctx, `select distinct category from news where trim(category)<>'' order by category`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []string{}
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		result = append(result, category)
	}
	return result, rows.Err()
}

func (s *Store) SaveAdminNews(ctx context.Context, item AdminNews) (int64, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	if item.ID == 0 {
		err = tx.QueryRow(ctx, `insert into news(title,published_on,category,excerpt,emoji,cover_media_id,content,is_published,publication_status)
			values($1,$2,$3,$4,$5,$6,$7,$8,case when $8 then 'published' else 'draft' end) returning id`,
			item.Title, item.Date, item.Category, item.Excerpt, item.Emoji, item.CoverMediaID, item.Content, item.IsPublished).Scan(&item.ID)
	} else {
		command, updateErr := tx.Exec(ctx, `update news set title=$1,published_on=$2,category=$3,excerpt=$4,emoji=$5,cover_media_id=$6,content=$7,
			is_published=$8,publication_status=case when $8 then 'published' else 'draft' end,publish_at=null,updated_at=now() where id=$9`,
			item.Title, item.Date, item.Category, item.Excerpt, item.Emoji, item.CoverMediaID, item.Content, item.IsPublished, item.ID)
		err = updateErr
		if err == nil && command.RowsAffected() == 0 {
			err = ErrNotFound
		}
	}
	if err != nil {
		return 0, err
	}
	if _, err = tx.Exec(ctx, `delete from news_photos where news_id=$1`, item.ID); err != nil {
		return 0, err
	}
	for position, mediaID := range item.PhotoMediaIDs {
		if _, err = tx.Exec(ctx, `insert into news_photos(news_id,media_id,position) values($1,$2,$3)`, item.ID, mediaID, position); err != nil {
			return 0, err
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return item.ID, nil
}

func (s *Store) AdminService(ctx context.Context, id int64) (AdminService, error) {
	var item AdminService
	err := s.pool.QueryRow(ctx, `select s.id,s.title,s.tagline,s.description,s.features,s.price,s.emoji,s.path,s.color,s.cover_media_id,s.position,s.is_published,
		coalesce(array_agg(sp.media_id order by sp.position) filter(where sp.media_id is not null),'{}') from services s left join service_photos sp on sp.service_id=s.id where s.id=$1 group by s.id`, id).
		Scan(&item.ID, &item.Title, &item.Tagline, &item.Description, &item.Features, &item.Price, &item.Emoji, &item.Path, &item.Color, &item.CoverMediaID, &item.Position, &item.IsPublished, &item.PhotoMediaIDs)
	if errors.Is(err, pgx.ErrNoRows) {
		return AdminService{}, ErrNotFound
	}
	return item, err
}

func (s *Store) SaveAdminService(ctx context.Context, item AdminService) (int64, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)
	if item.ID == 0 {
		err = tx.QueryRow(ctx, `insert into services(title,tagline,description,features,price,emoji,path,color,cover_media_id,position,is_published,publication_status)
		values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,case when $11 then 'published' else 'draft' end) returning id`, item.Title, item.Tagline, item.Description, item.Features, item.Price, item.Emoji, item.Path, item.Color, item.CoverMediaID, item.Position, item.IsPublished).Scan(&item.ID)
	} else {
		command, updateErr := tx.Exec(ctx, `update services set title=$1,tagline=$2,description=$3,features=$4,price=$5,emoji=$6,path=$7,color=$8,cover_media_id=$9,position=$10,is_published=$11,publication_status=case when $11 then 'published' else 'draft' end,publish_at=null where id=$12`, item.Title, item.Tagline, item.Description, item.Features, item.Price, item.Emoji, item.Path, item.Color, item.CoverMediaID, item.Position, item.IsPublished, item.ID)
		err = updateErr
		if err == nil && command.RowsAffected() == 0 {
			err = ErrNotFound
		}
	}
	if err != nil {
		return 0, err
	}
	if _, err = tx.Exec(ctx, `delete from service_photos where service_id=$1`, item.ID); err != nil {
		return 0, err
	}
	for position, mediaID := range item.PhotoMediaIDs {
		if _, err = tx.Exec(ctx, `insert into service_photos(service_id,media_id,position)values($1,$2,$3)`, item.ID, mediaID, position); err != nil {
			return 0, err
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return item.ID, nil
}

type pgxQuerier interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}
