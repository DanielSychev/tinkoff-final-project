package adrepo

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"homework9/internal/ads"
	"homework9/internal/app"
	"sync"
	"time"
)

var ErrNotAuthor = errors.New("not author")
var ErrValidate = errors.New("validation error")
var ErrNotCreated = errors.New("not created")
var ErrWasDeleted = errors.New("has been already deleted")

type Repo struct {
	index      int64
	a          []ads.Ad
	users      []ads.User
	usersIndex int64
	mu         *sync.Mutex
	conn       *pgx.Conn
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
	r.a = append(r.a, ads.Ad{
		ID:          r.index,
		Title:       Title,
		Text:        Text,
		AuthorID:    UserID,
		Published:   false,
		DateCreated: time.Now().Format("2006-01-02 15:04:05"),
		DateUpdated: time.Now().Format("2006-01-02 15:04:05"),
	})
	r.index++
	return &r.a[r.index-1], nil
}

func (r *Repo) UpdatePublished(ID int64, UserID int64, Published bool) (*ads.Ad, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.a[ID].AuthorID != UserID {
		return nil, ErrNotAuthor
	}
	r.a[ID].Published = Published
	r.a[ID].DateUpdated = time.Now().Format("2006-01-02 15:04:05")
	return &r.a[ID], nil
}

func (r *Repo) UpdateTextAndTitle(ID int64, UserID int64, Title string, Text string) (*ads.Ad, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.a[ID].AuthorID != UserID {
		return nil, ErrNotAuthor
	}
	if !validate(Title, Text) {
		return nil, ErrValidate
	}
	r.a[ID].Text = Text
	r.a[ID].Title = Title
	r.a[ID].DateUpdated = time.Now().Format("2006-01-02 15:04:05")
	return &r.a[ID], nil
}

func (r *Repo) GetList(filter ads.AdFilter) ([]*ads.Ad, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var res = make([]*ads.Ad, 0)
	for _, elem := range r.a {
		if elem.Deleted {
			continue
		}
		if filter.Pub && !elem.Published {
			continue
		}
		if filter.Auth != -1 && elem.AuthorID != filter.Auth {
			continue
		}
		if filter.Title != "" && elem.Title != filter.Title {
			continue
		}
		res = append(res, &elem)
	}
	return res, nil
}

func (r *Repo) GetByID(ID int64) (*ads.Ad, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ID >= r.index {
		return nil, ErrNotCreated
	}
	if r.a[ID].Deleted {
		return nil, ErrWasDeleted
	}
	return &r.a[ID], nil
}

func (r *Repo) DeleteAd(ID int64, UserID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ID >= r.index {
		return ErrNotCreated
	}
	if r.a[ID].AuthorID != UserID {
		return ErrNotAuthor
	}
	if r.a[ID].Deleted {
		return ErrWasDeleted
	}
	r.a[ID].Deleted = true
	return nil
}

func (r *Repo) CreateUser(Name string) (*ads.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users = append(r.users, ads.User{Name: Name, ID: r.usersIndex})
	r.usersIndex++
	return &r.users[r.usersIndex-1], nil
}

func (r *Repo) GetUser(ID int64) (*ads.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ID >= r.usersIndex {
		return nil, ErrNotCreated
	}
	if r.users[ID].Deleted {
		return nil, ErrWasDeleted
	}
	return &r.users[ID], nil
}

func (r *Repo) DeleteUser(ID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ID >= r.usersIndex {
		return ErrNotCreated
	}
	if r.users[ID].Deleted {
		return ErrWasDeleted
	}
	r.users[ID].Deleted = true
	return nil
}

func New() app.Repository {
	return &Repo{index: 0, a: []ads.Ad{}, mu: new(sync.Mutex), users: []ads.User{}, usersIndex: 0}
}
