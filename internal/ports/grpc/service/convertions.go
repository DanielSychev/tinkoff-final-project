package service

import (
	"homework9/internal/ads"
	"homework9/internal/ports/grpc"
)

func ToAdResponse(a *ads.Ad) *grpc.AdResponse {
	return &grpc.AdResponse{
		Id:          a.ID,
		Title:       a.Title,
		Text:        a.Text,
		AuthorId:    a.AuthorID,
		Published:   a.Published,
		DateCreated: a.DateCreated,
		DateUpdated: a.DateUpdated,
	}
}

func ToListAdResponse(a []*ads.Ad) *grpc.ListAdResponse {
	var list = make([]*grpc.AdResponse, len(a))
	for i := range a {
		list[i] = ToAdResponse(a[i])
	}
	return &grpc.ListAdResponse{List: list}
}

func ToUserResponse(u *ads.User) *grpc.UserResponse {
	return &grpc.UserResponse{
		Id:   u.ID,
		Name: u.Name,
	}
}
