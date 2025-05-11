package adrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"homework9/internal/ads"
	"homework9/internal/app"
	"sync"
)

var ErrNotAuthor = errors.New("not author")
var ErrValidate = errors.New("validation error")
var ErrNotCreated = errors.New("not created")
var ErrWasDeleted = errors.New("has been already deleted")

const insertAdd = "INSERT INTO adds(title, text, author_id) VALUES($1, $2, $3) RETURNING *"
const selectAuthorId = "SELECT author_id FROM adds WHERE id = $1"
const selectAdd = "SELECT * FROM adds WHERE id = $1"
const updateAddPublished = "UPDATE adds SET published = $2 WHERE id = $1 RETURNING *"
const updateTextAndTitle = "UPDATE adds SET title = $2, text = $3 WHERE id = $1 RETURNING *"
const deleteAdd = "DELETE FROM adds WHERE id = $1"

const insertUser = "INSERT INTO users(name) VALUES($1) RETURNING *"
const selectUser = "SELECT * FROM users WHERE id = $1"
const deleteUser = "DELETE FROM users WHERE id = $1"

type Repo struct {
	mu   *sync.Mutex
	conn *pgx.Conn
	ctx  context.Context
}

func validate(Title string, Text string) bool {
	return Title != "" && len(Title) < 100 && Text != "" && len(Text) < 500
}

func (r *Repo) Create(Title string, Text string, UserID int64) (*ads.Ad, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !validate(Title, Text) {
		return nil, ErrValidate
	}
	ad := &ads.Ad{}
	err := r.conn.QueryRow(r.ctx, insertAdd, Title, Text, UserID).Scan(
		&ad.ID, &ad.Title, &ad.Text, &ad.AuthorID,
		&ad.Published, &ad.DateCreated, &ad.DateUpdated,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create ad: %w", err)
	}
	return ad, nil
}

func (r *Repo) UpdatePublished(ID int64, UserID int64, Published bool) (*ads.Ad, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var auId int64
	err := r.conn.QueryRow(r.ctx, selectAuthorId, ID).Scan(&auId)
	if err != nil {
		return nil, fmt.Errorf("unable to select with such id: %w", err)
	}
	if auId != UserID {
		return nil, ErrNotAuthor
	}
	ad := &ads.Ad{}
	err = r.conn.QueryRow(r.ctx, updateAddPublished, ID, Published).Scan(
		&ad.ID, &ad.Title, &ad.Text, &ad.AuthorID,
		&ad.Published, &ad.DateCreated, &ad.DateUpdated,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to update published ad: %w", err)
	}
	return ad, nil
}

func (r *Repo) UpdateTextAndTitle(ID int64, UserID int64, Title string, Text string) (*ads.Ad, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !validate(Title, Text) {
		return nil, ErrValidate
	}
	var auId int64
	err := r.conn.QueryRow(r.ctx, selectAuthorId, ID).Scan(&auId)
	if err != nil {
		return nil, fmt.Errorf("unable to select with such id: %w", err)
	}
	if auId != UserID {
		return nil, ErrNotAuthor
	}
	ad := &ads.Ad{}
	err = r.conn.QueryRow(r.ctx, updateTextAndTitle, ID, Title, Text).Scan(
		&ad.ID, &ad.Title, &ad.Text, &ad.AuthorID,
		&ad.Published, &ad.DateCreated, &ad.DateUpdated,
	)
	return ad, nil
}

func (r *Repo) GetList(filter ads.AdFilter) ([]*ads.Ad, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var res = make([]*ads.Ad, 0)
	var i = 1
	for {
		ad := &ads.Ad{}
		err := r.conn.QueryRow(r.ctx, selectAdd, i).Scan(
			&ad.ID, &ad.Title, &ad.Text, &ad.AuthorID,
			&ad.Published, &ad.DateCreated, &ad.DateUpdated,
		)
		i += 1
		if err != nil {
			break
		}
		if filter.Pub && !ad.Published {
			continue
		}
		if filter.Auth != -1 && ad.AuthorID != filter.Auth {
			continue
		}
		if filter.Title != "" && ad.Title != filter.Title {
			continue
		}
		res = append(res, ad)
	}
	return res, nil
}

func (r *Repo) GetByID(ID int64) (*ads.Ad, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	ad := &ads.Ad{}
	err := r.conn.QueryRow(r.ctx, selectAdd, ID).Scan(
		&ad.ID, &ad.Title, &ad.Text, &ad.AuthorID,
		&ad.Published, &ad.DateCreated, &ad.DateUpdated,
	)
	if err != nil {
		return nil, ErrNotCreated
	}
	return ad, nil
}

func (r *Repo) DeleteAd(ID int64, UserID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var auId int64
	err := r.conn.QueryRow(r.ctx, selectAuthorId, ID).Scan(&auId)
	if err != nil {
		return ErrNotCreated
	}
	if auId != UserID {
		return ErrNotAuthor
	}
	_, err = r.conn.Exec(r.ctx, deleteAdd, ID)
	if err != nil {
		return fmt.Errorf("unable to delete ad: %w", err)
	}
	return nil
}

func (r *Repo) CreateUser(Name string) (*ads.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user := &ads.User{}
	r.conn.QueryRow(r.ctx, insertUser, Name).Scan(&user.ID, &user.Name)
	return user, nil
}

func (r *Repo) GetUser(ID int64) (*ads.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user := &ads.User{}
	if err := r.conn.QueryRow(r.ctx, selectUser, ID).Scan(&user.ID, &user.Name); err != nil {
		return nil, ErrNotCreated
	}
	return user, nil
}

func (r *Repo) DeleteUser(ID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, err := r.conn.Exec(r.ctx, deleteUser, ID); err != nil {
		return ErrNotCreated
	}
	return nil
}

func New(ctx context.Context, conn *pgx.Conn) app.Repository {
	return &Repo{ctx: ctx, conn: conn, mu: new(sync.Mutex)}
}
