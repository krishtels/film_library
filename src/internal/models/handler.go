package models

import (
	"errors"
	"film-library/src/internal/tools"
	"log"
	"net/http"
)

type Handler struct {
	service ActorService
}

func NewHandler(as ActorService) *Handler {
	return &Handler{
		service: as,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.GetAll(r.Context())
	if err != nil {
		log.Printf("ERROR: can't get actors err=%s\n", err.Error())
		tools.InternalServerError(w, r)
		return
	}

	tools.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	var req ActorInfo
	ok := tools.BindJSON(w, r, &req)
	if !ok {
		return
	}

	res, err := h.service.Add(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: can't to add actor err=%s\n", err.Error())

		var ve *tools.ValidationError
		if errors.As(err, &ve) {
			tools.JSON(w, r, http.StatusBadRequest, &tools.ErrorMessage{
				ErrorType: tools.ErrorTypeValidation,
				Body:      ve.Error(),
			})
			return
		}
		tools.InternalServerError(w, r)
		return
	}

	tools.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	req := ActorIdRequest{
		ID: r.PathValue("id"),
	}

	res, err := h.service.Get(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: can't to get actor err=%s\n", err.Error())
		if errors.Is(err, ErrIdInvalid) || errors.Is(err, ErrActorNotExist) {
			tools.NotFound(w, r)
			return
		}

		tools.InternalServerError(w, r)
		return
	}

	tools.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	req := ActorIdInfoRequest{
		ID: r.PathValue("id"),
	}

	ok := tools.BindJSON(w, r, &req.Info)
	if !ok {
		return
	}

	res, err := h.service.Update(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: can't update actor err=%s\n", err.Error())

		if errors.Is(err, ErrIdInvalid) || errors.Is(err, ErrActorNotExist) {
			tools.NotFound(w, r)
			return
		}

		var ve *tools.ValidationError
		if errors.As(err, &ve) {
			tools.JSON(w, r, http.StatusBadRequest, &tools.ErrorMessage{
				ErrorType: tools.ErrorTypeValidation,
				Body:      ve.Error(),
			})
			return
		}

		if errors.Is(err, ErrEmptyUpdate) {
			tools.JSON(w, r, http.StatusBadRequest, &tools.ErrorMessage{
				ErrorType: tools.ErrorTypeValidation,
				Body:      "empty update",
			})
			return
		}

		tools.InternalServerError(w, r)
		return
	}

	tools.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	req := ActorIdRequest{
		ID: r.PathValue("id"),
	}

	res, err := h.service.Delete(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: can't get actor err=%s\n", err.Error())
		if errors.Is(err, ErrIdInvalid) || errors.Is(err, ErrActorNotExist) {
			tools.NotFound(w, r)
			return
		}

		tools.InternalServerError(w, r)
		return
	}

	tools.JSON(w, r, http.StatusOK, res)
}
