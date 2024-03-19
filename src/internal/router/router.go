package router

import (
	"net/http"

	"film-library/src/internal/config"
	"film-library/src/internal/film"
	"film-library/src/internal/models"
	"film-library/src/internal/user"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(cfg *config.Config, uh user.UserHandler, ah models.ActorHandler, fh film.FilmHandler) *Router {
	mux := http.NewServeMux()

	authMW := NewAuthMiddleware(cfg.SigningKey, false)
	adminOnlyMW := NewAuthMiddleware(cfg.SigningKey, true)
	logMW := NewLogMiddleware()

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("GET /docs/html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeFile(w, r, cfg.DocsHTML)
	})
	mux.HandleFunc("GET /docs/yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.DocsYAML)
	})

	mux.Handle("POST /signup", logMW(http.HandlerFunc(uh.CreateUser)))
	mux.Handle("POST /signin", logMW(http.HandlerFunc(uh.Login)))
	mux.Handle("DELETE /signout", logMW(http.HandlerFunc(uh.Logout)))

	mux.Handle("GET /actors", logMW(authMW(http.HandlerFunc(ah.GetAll))))
	mux.Handle("POST /actors", logMW(adminOnlyMW(http.HandlerFunc(ah.Add))))
	mux.Handle("GET /actors/{id}", logMW(authMW(http.HandlerFunc(ah.Get))))
	mux.Handle("PUT /actors/{id}", logMW(adminOnlyMW(http.HandlerFunc(ah.Update))))
	mux.Handle("DELETE /actors/{id}", logMW(adminOnlyMW(http.HandlerFunc(ah.Delete))))

	mux.Handle("GET /films", logMW(authMW(http.HandlerFunc(fh.GetFilms))))
	mux.Handle("POST /films", logMW(adminOnlyMW(http.HandlerFunc(fh.AddFilm))))
	mux.Handle("GET /films/{id}", logMW(authMW(http.HandlerFunc(fh.GetFilm))))
	mux.Handle("PUT /films/{id}", logMW(adminOnlyMW(http.HandlerFunc(fh.UpdateFilm))))
	mux.Handle("DELETE /films/{id}", logMW(adminOnlyMW(http.HandlerFunc(fh.DeleteFilm))))
	mux.Handle("GET /films/{id}/actors", logMW(authMW(http.HandlerFunc(fh.GetFilmActors))))
	mux.Handle("PUT /films/{id}/actors", logMW(adminOnlyMW(http.HandlerFunc(fh.AddFilmActors))))
	mux.Handle("DELETE /films/{id}/actors", logMW(adminOnlyMW(http.HandlerFunc(fh.DeleteFilmActors))))

	return &Router{
		mux: mux,
	}
}

func (r *Router) Run(addr string) error {
	return http.ListenAndServe(addr, r.mux)
}
