package dropdb

import (
	"fmt"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct{}

type Response struct {
	response.Response
}

type DBDropper interface {
	DropTable() error
}

func New(log *slog.Logger, dbDropper DBDropper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.db.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		err := dbDropper.DropTable()

		if err != nil {
			log.Error("failed to drop database", sl.Err(err))
			render.JSON(w, r, fmt.Errorf("failed to drop db"))
			return
		}
		log.Info("db was dropped")
		render.JSON(w, r, Response{
			Response: response.OK(),
		})
	}

}
