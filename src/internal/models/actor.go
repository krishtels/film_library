package models

import (
	"context"
	"net/http"
	"time"
)

type Actor struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Sex      string    `json:"sex"`
	Birthday time.Time `json:"birthday"`
	Films    []string  `json:"films"`
}

type ActorRepository interface {
	Get(ctx context.Context, id int) (*Actor, error)
	Add(ctx context.Context, a *Actor) (*Actor, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, a *Actor) error
	GetAll(ctx context.Context) ([]*Actor, error)
}

type ActorService interface {
	GetAll(ctx context.Context) ([]*ActorResponse, error)
	Add(ctx context.Context, req *ActorInfo) (*ActorResponse, error)
	Get(ctx context.Context, req *ActorIdRequest) (*ActorResponse, error)
	Update(ctx context.Context, req *ActorIdInfoRequest) (*ActorResponse, error)
	Delete(ctx context.Context, req *ActorIdRequest) (*ActorResponse, error)
}

type ActorHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type ActorInfo struct {
	Name     string `json:"name"`
	Sex      string `json:"sex"`
	Birthday string `json:"birthday"`
}

type ActorResponse struct {
	ID    int       `json:"id"`
	Info  ActorInfo `json:"info"`
	Films []string  `json:"films,omitempty"`
}

type ActorIdRequest struct {
	ID string
}

type ActorIdInfoRequest struct {
	ID   string
	Info ActorInfo
}
