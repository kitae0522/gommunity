package dto

import (
	"github.com/kitae0522/gommunity/internal/model"
)

type CreateThreadRequest struct {
	UserID       string  `json:"userID"`
	Title        *string `json:"title"`
	ImgUrl       *string `json:"imgUrl"`
	Content      string  `json:"content" validate:"required"`
	ParentThread *int    `json:"parentThread"`
	NextThread   *int    `json:"nextThread"`
	PrevThread   *int    `json:"prevThread"`
}

type CreateThreadReponse struct {
	IsError    bool              `json:"isError"`
	StatusCode int               `json:"statusCode"`
	Message    string            `json:"message"`
	Thread     model.ThreadModel `json:"thread"`
}

type ListThreadRequest struct {
	PageNumber int `query:"pageNumber"`
	PageSize   int `query:"pageSize"`
}

type ListThreadResponse struct {
	IsError    bool                `json:"isError"`
	StatusCode int                 `json:"statusCode"`
	Message    string              `json:"message"`
	Threads    []model.ThreadModel `json:"threads"`
}

type ListThreadByHandleRequest struct {
	Handle string `params:"handle" validate:"required"`
}

type ListThreadByHandleResponse struct {
	IsError    bool                `json:"isError"`
	StatusCode int                 `json:"statusCode"`
	Message    string              `json:"message"`
	Handle     string              `json:"handle"`
	Threads    []model.ThreadModel `json:"threads"`
}

type GetThreadByIDRequest struct {
	ThreadID int `params:"threadID" validate:"required"`
}

type GetThreadByIDResponse struct {
	IsError    bool                `json:"isError"`
	StatusCode int                 `json:"statusCode"`
	Message    string              `json:"message"`
	Thread     *model.ThreadModel  `json:"thread"`
	SubThread  []model.ThreadModel `json:"subThread"`
}

type RemoveThreadByIDRequest struct {
	ID       string `json:"id" validate:"required"`
	ThreadID int    `params:"threadID" validate:"required"`
}

type InteractionRequest struct {
	ThreadID int `json:"threadID" validate:"required"`
}
