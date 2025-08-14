package getall

import (
	"log/slog"
	"net/http"
	res "url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type AllURLGetter interface {
	GetAllURLs() ([]sqlite.URL, error)
}

func New(log *slog.Logger, allUrlGetter AllURLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.getall.New"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		urls, err := allUrlGetter.GetAllURLs()
		if err != nil {
			log.Error("failed to get all urls", slog.Any("err", err))
			render.JSON(w, r, res.Error("internal error"))
			return
		}
		log.Info("fetched all urls")
		render.JSON(w, r, urls)
	}
}
