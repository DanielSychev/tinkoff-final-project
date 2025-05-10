package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework9/internal/ads"
	"homework9/internal/app"
	"homework9/internal/ports/grpc"
)

type MyServer struct {
	a app.App
	grpc.UnimplementedAdServiceServer
}

func NewMyServer(ap app.App) *MyServer {
	return &MyServer{a: ap}
}

func (s *MyServer) CreateAd(c context.Context, adReq *grpc.CreateAdRequest) (*grpc.AdResponse, error) {
	//log.Println("you are here")
	adResp, err := s.a.CreateAd(c, adReq.Title, adReq.Text, adReq.UserId)
	if err != nil {
		return nil, err
	}
	return ToAdResponse(adResp), nil
}

func (s *MyServer) ChangeAdStatus(c context.Context, adReq *grpc.ChangeAdStatusRequest) (*grpc.AdResponse, error) {
	adResp, err := s.a.ChangeAdStatus(c, adReq.AdId, adReq.UserId, adReq.Published)
	if err != nil {
		return nil, err
	}
	return ToAdResponse(adResp), nil
}

func (s *MyServer) UpdateAd(c context.Context, adReq *grpc.UpdateAdRequest) (*grpc.AdResponse, error) {
	adResp, err := s.a.UpdateAd(c, adReq.AdId, adReq.UserId, adReq.Title, adReq.Text)
	if err != nil {
		return nil, err
	}
	return ToAdResponse(adResp), nil
}

func (s *MyServer) ListAds(c context.Context, _ *emptypb.Empty) (*grpc.ListAdResponse, error) {
	filter := ads.AdFilter{Pub: true, Auth: -1, Title: ""}
	adResp, err := s.a.GetList(c, filter)
	if err != nil {
		return nil, err
	}
	return ToListAdResponse(adResp), nil
}

func (s *MyServer) CreateUser(c context.Context, in *grpc.CreateUserRequest) (*grpc.UserResponse, error) {
	resp, err := s.a.CreateUser(c, in.Name)
	if err != nil {
		return nil, err
	}
	return ToUserResponse(resp), nil
}

func (s *MyServer) GetUser(c context.Context, in *grpc.GetUserRequest) (*grpc.UserResponse, error) {
	resp, err := s.a.GetUser(c, in.Id)
	if err != nil {
		return nil, err
	}
	return ToUserResponse(resp), nil
}

func (s *MyServer) DeleteUser(c context.Context, in *grpc.DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.a.DeleteUser(c, in.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *MyServer) DeleteAd(c context.Context, in *grpc.DeleteAdRequest) (*emptypb.Empty, error) {
	err := s.a.DeleteAd(c, in.AdId, in.AuthorId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
