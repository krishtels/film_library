package user

import (
	"errors"
	"log"
	"net/http"

	"film-library/src/internal/tools"
)

type Handler struct {
	service UserService
}

func NewHandler(s UserService) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if ok := tools.BindJSON(w, r, &req); !ok {
		return
	}

	res, err := h.service.CreateUser(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to create user err=%s\n", err.Error())

		var ve *tools.ValidationError
		if errors.As(err, &ve) {
			tools.JSON(w, r, http.StatusBadRequest, &tools.ErrorMessage{
				ErrorType: tools.ErrorTypeValidation,
				Body:      ve.Error(),
			})
			return
		}

		if errors.Is(err, ErrUserExist) {
			tools.JSON(w, r, http.StatusBadRequest, &tools.ErrorMessage{
				ErrorType: tools.ErrorTypeConflict,
				Body:      "User already exists",
			})
			return
		}

		tools.InternalServerError(w, r)
		return
	}

	tools.JSON(w, r, http.StatusOK, res)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if ok := tools.BindJSON(w, r, &req); !ok {
		return
	}

	res, err := h.service.Login(r.Context(), &req)
	if err != nil {
		log.Printf("ERROR: failed to login user err=%s\n", err.Error())

		if errors.Is(err, ErrUserNotExist) || errors.Is(err, ErrPasswordIncorrect) {
			tools.JSON(w, r, http.StatusUnauthorized, &tools.ErrorMessage{
				ErrorType: tools.ErrorTypeValidation,
				Body:      "incorrect password",
			})
			return
		}

		tools.InternalServerError(w, r)
		return
	}

	tools.SetJWTCookie(w, res.AccessToken)
	tools.OK(w, r)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("jwt")
	if err != nil {
		log.Printf("ERROR: logout failed err=%s", err.Error())
		if errors.Is(err, http.ErrNoCookie) {
			tools.Unauthorized(w, r)
			return
		}

		tools.InternalServerError(w, r)
		return
	}

	tools.UnsetJWTCookie(w)
	tools.OK(w, r)
}
