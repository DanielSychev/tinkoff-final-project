package app

import (
	"context"
	"homework9/internal/ads"
)

type App interface {
	CreateAd(c context.Context, Title string, Text string, UserID int64) (*ads.Ad, error)
	ChangeAdStatus(c context.Context, ID int64, UserID int64, Published bool) (*ads.Ad, error)
	UpdateAd(c context.Context, ID int64, UserID int64, Title string, Text string) (*ads.Ad, error)
	GetList(c context.Context, filter ads.AdFilter) ([]*ads.Ad, error)
	GetByID(c context.Context, ID int64) (*ads.Ad, error)
	DeleteAd(c context.Context, ID int64, UserID int64) error
	CreateUser(c context.Context, Name string) (*ads.User, error)
	GetUser(c context.Context, ID int64) (*ads.User, error)
	DeleteUser(c context.Context, ID int64) error
}

type Repository interface {
	Create(Title string, Text string, UserID int64) (*ads.Ad, error)
	UpdatePublished(ID int64, UserID int64, Published bool) (*ads.Ad, error)
	UpdateTextAndTitle(ID int64, UserID int64, Title string, Text string) (*ads.Ad, error)
	GetList(filter ads.AdFilter) ([]*ads.Ad, error)
	GetByID(ID int64) (*ads.Ad, error)
	DeleteAd(ID int64, UserID int64) error
	CreateUser(Name string) (*ads.User, error)
	GetUser(ID int64) (*ads.User, error)
	DeleteUser(ID int64) error
}

type AppMethods struct {
	r Repository
}

func (apm *AppMethods) CreateAd(c context.Context, Title string, Text string, UserID int64) (*ads.Ad, error) {
	ad, err := apm.r.Create(Title, Text, UserID)
	if err != nil {
		return nil, err
	}
	return ad, nil
}

func (apm *AppMethods) ChangeAdStatus(c context.Context, ID int64, UserID int64, Published bool) (*ads.Ad, error) {
	ad, err := apm.r.UpdatePublished(ID, UserID, Published)
	if err != nil {
		return nil, err
	}
	return ad, nil
}

func (apm *AppMethods) UpdateAd(c context.Context, ID int64, UserID int64, Title string, Text string) (*ads.Ad, error) {
	ad, err := apm.r.UpdateTextAndTitle(ID, UserID, Title, Text)
	if err != nil {
		return nil, err
	}
	return ad, nil
}

func (apm *AppMethods) GetList(c context.Context, filter ads.AdFilter) ([]*ads.Ad, error) {
	return apm.r.GetList(filter)
}

func (apm *AppMethods) GetByID(c context.Context, ID int64) (*ads.Ad, error) {
	return apm.r.GetByID(ID)
}

func (apm *AppMethods) DeleteAd(c context.Context, ID int64, UserID int64) error {
	return apm.r.DeleteAd(ID, UserID)
}

func (apm *AppMethods) CreateUser(c context.Context, Name string) (*ads.User, error) {
	return apm.r.CreateUser(Name)
}

func (apm *AppMethods) GetUser(c context.Context, ID int64) (*ads.User, error) {
	return apm.r.GetUser(ID)
}

func (apm *AppMethods) DeleteUser(c context.Context, ID int64) error {
	return apm.r.DeleteUser(ID)
}

func NewApp(repo Repository) App {
	return &AppMethods{r: repo}
}
