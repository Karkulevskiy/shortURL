package delete

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.New()"
		log = log.With(
			slog.String("request_id", middleware.GetReqID(r.Context())),
			slog.String("op", op),
		)
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias can't be empty")
			render.JSON(w, r, response.Error("invalid alias"))
			return
		}
		err := urlDeleter.DeleteURL(alias)
		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, response.Error("internal error"))
		}
		log.Info("deleted url")
		render.JSON(w, r, http.StatusNoContent)
	}
}
