package main

import (
	"log/slog"
	"net/http"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/getall"
	"url-shortener/internal/http-server/handlers/url/redirect"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()
	log := config.SetUpLogger(cfg.Env)
	log.Info("starting url-shortener", slog.Any("env", cfg.Env))
	log.Debug("debug messages are enabled")
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", slog.Any("err", err))
		return
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.User: cfg.Password,
		}))
		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete.New(log, storage))
	})

	router.Get("/url/all", getall.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped")
	// TODO : init config : cleanenv
	// TODO : init logger : slog
	// TODO : init storage : sqllite
	// TODO : init router : chi, chi:render
	// TODO : run server :
}

// id, err := storage.SaveURL("https://amazon.com", "amazon")
// if err != nil {
// 	log.Error("failed to save url", sl.Any("err", err))
// 	return
// }
// log.Info("url was saved", sl.Any("ID", id))
// ursl, err := storage.GetAllURLs()
// if err != nil {
// 	log.Error("failed to fetch all urls", sl.Any("err", err))
// 	return
// }
// log.Info("urls was fetched", sl.Any("urls", ursl))
// google, err := storage.GetURL("google")
// if err != nil {
// 	log.Error("failed to fetch google url", sl.Any("err", err))
// 	return
// }
// log.Info("google", sl.Any("url", google))
// amazon, err := storage.GetURL("amazon")
// if err != nil {
// 	log.Error("failed to fetch amazon url", sl.Any("err", err))
// 	return
// }
// log.Info("amazon", sl.Any("url", amazon))
// err = storage.DeleteURL("notexists")
// if err != nil {
// 	log.Error("failed to delete url", sl.Any("err", err))
// 	return
// }
// log.Info("url deleted successfully")
