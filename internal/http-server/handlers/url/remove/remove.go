package remove

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener2/internal/logging"
	"url-shortener2/internal/response"
	"url-shortener2/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(logger *slog.Logger, db URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		const origin = "http-server.handlers.url.delete.New"

		logger = logger.With(
			slog.String("origin", origin),
			slog.String("request_id", middleware.GetReqID(req.Context())),
		)

		alias := chi.URLParam(req, "alias")
		if alias == "" {
			logger.Info("empty alias")
			render.JSON(w, req, response.Error("alias must not be empty"))
			return
		}

		err := db.DeleteURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlNotDeleted) {
				logger.Info("alias doesnt exist", slog.String("alias", alias))
				render.JSON(w, req, response.Error("url not found"))
				return
			}

			logger.Error("unknown error", logging.Err(err))
			render.JSON(w, req, response.Error("internal error"))
			return
		}

		logger.Info("url deleted", slog.String("alias", alias))
		render.JSON(w, req, response.OK("alias"))
	}
}
