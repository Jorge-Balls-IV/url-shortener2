package redirect

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

//go:generate mockery --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(logger *slog.Logger, db URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		const origin = "handlers.url.redirect.New"

		logger := logger.With(
			slog.String("origin", origin),
			slog.String("request_id", middleware.GetReqID(req.Context())),
		)

		alias := chi.URLParam(req, "alias")

		if alias == "" {
			logger.Info("cannot redirect on empty alias")
			render.JSON(w, req, response.Error("cannot redirect on empty alias"))
			return
		}

		urlToSend, err := db.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlNotFound) {
				logger.Error("url not found", slog.String("alias", alias))
				render.JSON(w, req, response.Error("url not found"))
				return
			}

			logger.Error("error getting url", logging.Err(err))
			render.JSON(w, req, response.Error("internal error"))
			return
		}

		logger.Info("got url", slog.String("url", urlToSend))

		// перенаправляем по найденному url
		http.Redirect(w, req, urlToSend, http.StatusFound)
		logger.Info("user redirected")
	}
}
