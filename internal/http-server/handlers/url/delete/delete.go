package delete

import (
	"log/slog"
	"net/http"
	res "url-shortener/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type DeleteURL interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, deleteUrl DeleteURL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.getall.New"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, res.Error("invalid request"))
			return
		}
		err := deleteUrl.DeleteURL(alias)
		if err != nil {
			log.Error("failed to delete url", slog.Any("err", err))
			render.JSON(w, r, res.Error("internal error"))
			return
		}
		log.Info("url was deleted")
		render.JSON(w, r, res.OK())
	}
}
