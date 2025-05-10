package httpgin

import (
	"github.com/gin-gonic/gin"
	"homework9/internal/ads"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type adResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	AuthorID    int64  `json:"author_id"`
	Published   bool   `json:"published"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type DeleteAdRequest struct {
	AuthorId int64 `json:"author_id"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

type userResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func AdSuccessResponse(ad *ads.Ad) gin.H {
	return gin.H{
		"data": adResponse{
			ID:          ad.ID,
			Title:       ad.Title,
			Text:        ad.Text,
			AuthorID:    ad.AuthorID,
			Published:   ad.Published,
			DateCreated: ad.DateCreated.Format("2006-01-02 15:04:05"),
			DateUpdated: ad.DateUpdated.Format("2006-01-02 15:04:05"),
		},
		"error": nil,
	}
}

func UserSuccessResponse(user *ads.User) gin.H {
	return gin.H{
		"data": userResponse{
			ID:   user.ID,
			Name: user.Name,
		},
		"error": nil,
	}
}

func AdListSuccessResponse(ad []*ads.Ad) gin.H {
	resp := make([]adResponse, len(ad))
	for i := range ad {
		resp[i] = adResponse{
			ID:          ad[i].ID,
			Title:       ad[i].Title,
			Text:        ad[i].Text,
			AuthorID:    ad[i].AuthorID,
			Published:   ad[i].Published,
			DateCreated: ad[i].DateCreated.Format("2006-01-02 15:04:05"),
			DateUpdated: ad[i].DateUpdated.Format("2006-01-02 15:04:05"),
		}
	}
	return gin.H{
		"data":  resp,
		"error": nil,
	}
}

func ErrorResponse(err error) gin.H {
	return gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
